package apu

import "fmt"

type APU struct {
	frameCounter byte
}

func New() *APU {
	return &APU{}
}

func (a *APU) Read(addr uint16) byte {
	return 0
	panic(fmt.Sprintf("unhandled address %x", addr))
}

func (a *APU) Write(addr uint16, val byte) {
	return
	switch addr {
	case 0x4017:
		a.frameCounter = val
	default:
		panic(fmt.Sprintf("unhandled address %x (%x)", addr, val))
	}
}
