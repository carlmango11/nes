package ram

import (
	"Nes/log"
)

const size = 100000

// Stack is $0100 - $01FF
type RAM struct {
	data [size]byte
}

func New() *RAM {
	return &RAM{
		data: [size]byte{},
	}
}

func (r *RAM) Dump() [][]uint16 {
	var dump [][]uint16

	for i, v := range r.data {
		dump = append(dump, []uint16{
			uint16(i),
			uint16(v),
		})
	}

	return dump
}

func (r *RAM) Read(addr uint16) byte {
	v := r.data[addr]
	log.Debugf("ram: read %x (%v) from %x (%v)", v, v, addr, addr)

	return v
}

func (r *RAM) Write(addr uint16, v byte) {
	log.Debugf("ram: write %x (%v) to %x (%v)", v, v, addr, addr)
	r.data[addr] = v
}

func (r *RAM) Load(bytes []byte) {
	for i, v := range bytes {
		r.data[i] = v
	}
}
