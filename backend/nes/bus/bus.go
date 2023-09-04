package bus

import (
	"fmt"
	"github.com/carlmango11/nes/backend/nes/apu"
	"github.com/carlmango11/nes/backend/nes/input"
	"github.com/carlmango11/nes/backend/nes/log"
	"github.com/carlmango11/nes/backend/nes/ppu"
	"github.com/carlmango11/nes/backend/nes/ram"
	"github.com/carlmango11/nes/backend/nes/rom"
)

type readComponent interface {
	Read(uint16) byte
}

type writeComponent interface {
	Write(uint16, byte)
}

//2000-2007 = PPU (graphics chip) registers
//4000-401F = Sound, Joypads, Sprite DMA
//6000-7FFF = Cartridge RAM if present
//8000-FFFF = Cartridge ROM

type Bus struct {
	apu   *apu.APU
	ram   *ram.RAM
	rom   rom.ROM
	ppu   *ppu.PPU
	input *input.Input
}

func New(rom rom.ROM, ppu *ppu.PPU) *Bus {
	return &Bus{
		apu:   apu.New(),
		ram:   ram.New(),
		rom:   rom,
		ppu:   ppu,
		input: input.New(),
	}
}

func (b *Bus) Read(addr uint16) byte {
	log.Debugf("bus: read %x", addr)
	return b.getReadComponent(addr).Read(addr)
}

func (b *Bus) Write(addr uint16, v byte) {
	log.Debugf("bus: write %x", addr)
	b.getWriteComponent(addr).Write(addr, v)
}

func (b *Bus) getReadComponent(addr uint16) readComponent {
	if addr < 0 {
		panic(fmt.Sprintf("bus: read from invalid address %x", addr))
	}

	switch {
	case addr <= 0x07FF:
		return b.ram
	case addr < 0x2000:
		panic(fmt.Sprintf("bus: read from invalid address %x", addr))
	case addr <= 0x3FFF:
		return b.ppu
	case addr >= 0x4000 && addr <= 0x4015:
		return b.apu
	case addr >= 0x4016 && addr <= 0x4017:
		return b.input
	case addr < 0x8000:
		panic(fmt.Sprintf("bus: read from unhandled address %x", addr))
	default:
		return b.rom
	}
}

func (b *Bus) getWriteComponent(addr uint16) writeComponent {
	if addr < 0 {
		panic(fmt.Sprintf("bus: write to invalid address %x", addr))
	}

	switch {
	case addr <= 0x07FF:
		return b.ram
	case addr < 0x2000:
		panic(fmt.Sprintf("bus: write to invalid address %x", addr))
	case addr <= 0x3FFF:
		return b.ppu
	case addr <= 0x4015:
		return b.apu
	case addr == 0x4016:
		return b.input
	case addr == 0x4017:
		return b.apu
	case addr < 0x8000:
		panic(fmt.Sprintf("bus: write to unhandled address %x", addr))
	default:
		return b.rom
	}
}
