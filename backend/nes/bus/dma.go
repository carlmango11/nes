package bus

import (
	"github.com/carlmango11/nes/backend/nes/ppu"
	"github.com/carlmango11/nes/backend/nes/ram"
)

type DMA struct {
	ram *ram.RAM
	ppu *ppu.PPU
}

func (d *DMA) Write(addr uint16, v byte) {
	for i := 0; i < 256; i++ {
		lo := uint16(i)
		hi := uint16(v) << 8
		dmaAddr := hi | lo

		dmaData := d.ram.Read(dmaAddr)

		d.ppu.WriteOAM(byte(i), dmaData)
	}
}
