package rom

import (
	"Nes/log"
	"bytes"
	"fmt"
)

const bankSize = 8 * 1024
const prgChunkSize = 16 * 1024
const crhChunkSize = 8 * 1024

type ROM interface {
	Read(addr uint16) byte
}

type romData struct {
	prg [][]byte
	chr [][]byte
}

func FromBytes(bs []byte) ROM {
	metadata := Metadata{
		prgSize: bs[4],
		chrSize: bs[5],
	}

	f6 := bs[6]
	if (f6>>2)&0x1 == 1 {
		panic("has trainer")
	}

	mapperN := f6 >> 4
	log.Printf("mapper number: %v", mapperN)

	b := bytes.NewReader(bs)

	data := &romData{
		prg: readROMChunks(b, metadata.prgSize, prgChunkSize),
		chr: readROMChunks(b, metadata.chrSize, crhChunkSize),
	}

	switch mapperN {
	case 4:
		return newROM4(data)
	default:
		panic(fmt.Sprintf("unknown mapper: %v", mapperN))
	}
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

		chunks = append(chunks, chunk)
	}

	return chunks
}
