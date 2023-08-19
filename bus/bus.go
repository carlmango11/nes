package bus

import (
	"Nes/log"
	"Nes/ram"
	"Nes/rom"
	"fmt"
)

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
	rom rom.ROM
}

func New(rom rom.ROM) *Bus {
	return &Bus{
		ram: ram.New(),
		rom: rom,
	}
}

func (b *Bus) Read(addr uint16) byte {
	log.Debugf("bus: read %x", addr)

	if addr < 0 {
		panic(fmt.Sprintf("bus: read from invalid address %x", addr))
	}

	switch {
	case addr <= 0x07FF:
		return b.ram.Read(addr)
	case addr < 0x8000:
		panic(fmt.Sprintf("bus: read from unhandled address %x", addr))
	default:
		return b.rom.Read(addr)
	}
}

func (b *Bus) Write(addr uint16, v byte) {
	if addr <= 0x07FF {
		b.ram.Write(addr, v)
		return
	}

	panic(fmt.Sprintf("bus: write to unhandled address %x", addr))
}
