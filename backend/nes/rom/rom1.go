package rom

import (
	"fmt"
	"github.com/carlmango11/nes/backend/nes/log"
)

type NROM struct {
	prg  []byte
	data *romData
}

func newNROM(data *romData) *NROM {
	log.Debugf("banks prg: %v", len(data.prg))
	log.Debugf("banks chr: %v", len(data.chr))

	var prg []byte
	for _, x := range data.prg {
		prg = append(prg, x...)
	}

	return &NROM{
		prg:  prg,
		data: data,
	}
}

func (r *NROM) Write(addr uint16, val byte) {
	panic(fmt.Sprintf("cannot write to rom: %x %x", addr, val))
}

func (r *NROM) Read(addr uint16) byte {
	switch {
	case addr < 0x8000:
		panic(fmt.Sprintf("nrom: unhandler addr: %v", addr))
	case addr < 0xFFFF:
		log.Debugf("reading %x -> %x = %x", addr, addr-0x8000, r.prg[addr-0x8000])
		val := r.prg[addr-0x8000]

		return val
	}

	panic(fmt.Sprintf("nrom: unhandled addr: %x", addr))
}
