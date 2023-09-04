package ppu

import (
	"fmt"
	"github.com/carlmango11/nes/backend/nes/log"
	"log/slog"
	"sync"
)

//$0000-$0FFF	$1000	Pattern table 0
//$1000-$1FFF	$1000	Pattern table 1
//$2000-$23FF	$0400	Nametable 0
//$2400-$27FF	$0400	Nametable 1
//$2800-$2BFF	$0400	Nametable 2
//$2C00-$2FFF	$0400	Nametable 3
//$3000-$3EFF	$0F00	Mirrors of $2000-$2EFF
//$3F00-$3F1F	$0020	Palette RAM indexes
//$3F20-$3FFF	$00E0	Mirrors of $3F00-$3F1F

const (
	PPUCTRL   = 0
	PPUMASK   = 1
	PPUSTATUS = 2
	OAMADDR   = 3
	OAMDATA   = 4
	PPUSCROLL = 5
	PPUADDR   = 6
	PPUDATA   = 7
)

type PPU struct {
	log slog.Logger

	stateMu sync.Mutex
	state   [240][256]byte

	frame int
	line  int
	pixel int

	data []byte

	vramAddr  uint16
	vramRead  byte
	hiWritten bool

	registers        [8]byte
	verticalBlankNMI bool
}

func New() *PPU {
	return &PPU{
		line: -1,
	}
}

func (r *PPU) State() [240][256]byte {
	r.stateMu.Lock()
	defer r.stateMu.Unlock()

	return r.state
}

func (r *PPU) Tick() {
	//log.Debugf("PPU (%v, %v)", r.pixel, r.line)

	switch {
	case r.line == -1:
		r.preRenderScanline()
	case r.line == 0:
	case r.line >= 241 && r.line <= 260:
		if r.pixel == 1 {
			r.setVBlank()
		}
	}

	r.pixel++
	if r.pixel == 341 {
		r.pixel = 0

		r.line++
		if r.line == 260 {
			r.line = 0
		}
	}
}

func (r *PPU) Read(addr uint16) byte {
	if addr >= 0x2000 && addr <= 0x3FFF {
		return r.readRegister(addr % 8)
	}

	panic(fmt.Sprintf("ppu: invalid address %x", addr))
}

func (r *PPU) Write(addr uint16, val byte) {
	if addr >= 0x2000 && addr <= 0x3FFF {
		r.writeRegister(addr%8, val)
	}

	panic(fmt.Sprintf("ppu: invalid address %x", addr))
}

func (r *PPU) readRegister(addr uint16) byte {
	if addr == PPUSTATUS {
		// reset some statuses
		r.verticalBlankNMI = false
		r.hiWritten = false
	}

	return r.registers[addr]
}

func (r *PPU) writeRegister(addr uint16, val byte) {
	r.registers[addr] = val

	switch addr {
	case PPUCTRL:
		r.writePPUCTRL(val)
	case PPUSCROLL:
		panic("no scroll")
	case PPUADDR:
		r.writePPUADDR(val)
	case PPUDATA:
		r.writePPUDATA(val)
	}
}

func (r *PPU) writePPUADDR(val byte) {
	if r.hiWritten {
		r.vramAddr = uint16(val) << 8
		r.hiWritten = true
	} else {
		r.vramAddr |= uint16(val)
	}
}

func (r *PPU) readPPUDATA() byte {
	val := r.vramRead
	r.vramRead = r.data[r.vramAddr]

	inc := uint16(1)
	if r.getBit(PPUCTRL, 2) == 1 {
		inc = 32
	}

	r.vramAddr += inc

	return val
}

func (r *PPU) writePPUDATA(val byte) {
	r.data[r.vramAddr] = val
}

func (r *PPU) writePPUCTRL(val byte) {
	r.verticalBlankNMI = (val >> 7) == 1
}

func (r *PPU) setVBlank() {
	log.Debugf("PPU: setting vblank")
	r.setBit(PPUSTATUS, 7, 1)
}

func (r *PPU) preRenderScanline() {
}

func (r *PPU) setBit(register, bit, val byte) {
	r.registers[register] = r.registers[register] | (val << bit)
}

func (r *PPU) getBit(register, bit byte) byte {
	r.registers[register] = r.registers[register] | (val << bit)
}
