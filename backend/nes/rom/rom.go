package rom

import (
	"bytes"
	"fmt"
	"github.com/carlmango11/nes/backend/nes/log"
)

const bankSize = 8 * 1024
const prgChunkSize = 16 * 1024
const crhChunkSize = 8 * 1024

type ROM interface {
	Read(addr uint16) byte
	Write(addr uint16, val byte)
}

type Metadata struct {
	prgSize byte
	chrSize byte
}

type romData struct {
	prg [][]byte
	chr [][]byte
}

func FromBytes(bs []byte) ROM {
	b := bytes.NewReader(bs)

	header := read(b, 16)

	fmt.Printf("%v\n", string(header[0:3]))

	metadata := Metadata{
		prgSize: header[4],
		chrSize: header[5],
	}

	f6 := header[6]
	if (f6>>2)&0x1 == 1 {
		panic("has trainer")
	}

	mapperLo := f6 >> 4

	f7 := header[7]
	if (f7>>1)&0x1 == 1 {
		panic("has playchoice")
	}

	if (f7>>2)&0x11 == 2 {
		panic("has nes 2")
	}

	mapperN := (f7 & 0xF0) | mapperLo
	log.Printf("mapper number: %v", mapperN)

	data := &romData{
		prg: readROMChunks(b, metadata.prgSize, prgChunkSize),
		chr: readROMChunks(b, metadata.chrSize, crhChunkSize),
	}

	switch mapperN {
	case 0:
		return newNROM(data)
	case 3:
		return newROM3(data)
	case 4:
		return newROM4(data)
	default:
		panic(fmt.Sprintf("unknown mapper: %v", mapperN))
	}
}

func read(b *bytes.Reader, size int) []byte {
	data := make([]byte, size)
	n, err := b.Read(data)
	if n != size {
		log.Panicf("unexpected size read. expected %v; read %v", size, n)
	}
	if err != nil {
		log.Panicf("error reading rom: %v", err)
	}

	return data
}

func readROMChunks(b *bytes.Reader, size byte, chunkSize int) [][]byte {
	var chunks [][]byte

	iters := int(size) * (chunkSize / bankSize)

	for i := 0; i < iters; i++ {
		chunk := make([]byte, bankSize)

		n, err := b.Read(chunk)
		if err != nil {
			log.Panicf("err reading prg: %v %v", i, err)
		}
		if n != bankSize {
			log.Panicf("unexpected prg chunk: %v %v", i, n)
		}

		if i == 31 {
			for j := range chunk {
				if j > 0x100 {
					//log.Debugf("ADDR %x = %x", j, chunk[j])
				}
			}
		}

		chunks = append(chunks, chunk)
	}

	return chunks
}
