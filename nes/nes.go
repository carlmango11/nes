package nes

import (
	"Nes/ram"
	"fmt"
)

type NES struct {
}

//0000-07FF= RAM
//2000-2007 = PPU (graphics chip) registers
//4000-401F = Sound, Joypads, Sprite DMA
//6000-7FFF = Cartridge RAM if present
//8000-FFFF = Cartridge ROM

// FFFA - NMI vector
// FFFC - Reset vector (Read the 16-bit value here, that's where code execution begins!)
// FFFE - IRQ vector
type Bus struct {
	ram *ram.RAM
}

func (b *Bus) Read(addr uint16) byte {
	if addr <= 0x07FF {
		return b.ram.Read(addr)
	}

	panic(fmt.Sprintf("bus: read from unhandled address %x", addr))
}
