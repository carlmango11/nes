package rom

import (
	"fmt"
	"github.com/carlmango11/nes/backend/nes/log"
)

type ROM3 struct {
	prg  []byte
	data *romData

	chrBank byte
}

func newROM3(data *romData) *ROM3 {
	log.Printf("banks prg: %v", len(data.prg))
	log.Printf("banks chr: %v", len(data.chr))

	var prg []byte
	for _, x := range data.prg {
		prg = append(prg, x...)
	}

	return &ROM3{
		prg:  prg,
		data: data,
	}
}

func (r *ROM3) Write(addr uint16, val byte) {
	if addr >= 0x8000 {
		// bank switch
		r.chrBank = val & 0x11
		log.Printf("Bank : %v", r.chrBank)
		return
	}

	panic(fmt.Sprintf("cannot write to rom: %x %x", addr, val))
}

func (r *ROM3) Read(addr uint16) byte {
	switch {
	case addr < 0x2000:
		return r.data.chr[r.chrBank][addr]
	case addr < 0x8000:
		panic(fmt.Sprintf("rom3: unhandler addr: %v", addr))
	case addr < 0xFFFF:
		log.Debugf("reading %x -> %x = %x", addr, addr-0x8000, r.prg[addr-0x8000])
		val := r.prg[addr-0x8000]

		return val
	}

	panic(fmt.Sprintf("rom3: unhandled addr: %x", addr))
}
