package ppu

import (
	"fmt"
	"github.com/carlmango11/nes/backend/nes/log"
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

var baseNametables = map[byte]uint16{
	0: 0x2000,
	1: 0x2400,
	2: 0x2800,
	3: 0x2C00,
}

type PPU struct {
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

	written bool
}

func New() *PPU {
	var state [240][256]byte
	for i := range state {
		for j := range state[i] {
			state[i][j] = 1
		}
	}

	return &PPU{
		line:  -1,
		state: state,
		data:  make([]byte, 16*1024),
	}
}

func (p *PPU) State() [240][256]byte {
	p.stateMu.Lock()
	defer p.stateMu.Unlock()

	return p.state
}

func (p *PPU) Tick() bool {
	var interrupt bool

	//log.Printf("PPU (%v, %v)", p.pixel, p.line)

	switch {
	case p.line == -1:
		p.preRenderScanline()
	case p.line == 0:
		//case r.line >= 1 && r.line <= 240:
		p.renderScanLine()
	case p.line == 1:
		p.renderFrame()
	case p.line >= 241 && p.line <= 260:
		if p.pixel == 1 {
			p.setVBlank()

			if p.getBit(PPUCTRL, 7) == 1 {
				interrupt = true
			}
		}
	}

	p.pixel++
	if p.pixel == 341 {
		p.pixel = 0

		p.line++
		if p.line == 260 {
			p.line = 0
		}
	}

	return interrupt
}

func (p *PPU) Read(addr uint16) byte {
	if addr >= 0x2000 && addr <= 0x3FFF {
		return p.readRegister(addr % 8)
	}

	panic(fmt.Sprintf("ppu: invalid address %x", addr))
}

func (p *PPU) Write(addr uint16, val byte) {
	if addr >= 0x2000 && addr <= 0x3FFF {
		p.writeRegister(addr%8, val)
		return
	}

	panic(fmt.Sprintf("ppu: invalid address %x", addr))
}

func (p *PPU) readRegister(addr uint16) byte {
	switch addr {
	case PPUSTATUS:
		// reset some statuses
		p.verticalBlankNMI = false
		p.hiWritten = false
	case PPUDATA:
		return p.readPPUDATA()
	}

	return p.registers[addr]
}

func (p *PPU) writeRegister(addr uint16, val byte) {
	p.registers[addr] = val

	//log.Printf("ppu: writing %x to %x", val, addr)

	switch addr {
	case PPUCTRL:
		p.writePPUCTRL(val)
	case PPUSCROLL:
		panic("no scroll")
	case PPUADDR:
		p.writePPUADDR(val)
	case PPUDATA:
		p.writePPUDATA(val)
	}
}

func (p *PPU) writePPUADDR(val byte) {
	if !p.hiWritten {
		p.vramAddr = uint16(val) << 8
		p.hiWritten = true
	} else {
		p.vramAddr |= uint16(val)
	}
}

func (p *PPU) readPPUDATA() byte {
	val := p.vramRead
	p.vramRead = p.data[p.vramAddr]

	p.incVramAddr()

	return val
}

func (p *PPU) incVramAddr() {
	if p.getBit(PPUCTRL, 2) == 1 {
		p.vramAddr += 32
	} else {
		p.vramAddr += 1
	}
}

func (p *PPU) writePPUDATA(val byte) {
	log.Printf("writing %x to vram at %x", val, p.vramAddr)
	p.data[p.vramAddr] = val
	p.incVramAddr()
}

func (p *PPU) writePPUCTRL(val byte) {
	p.verticalBlankNMI = (val >> 7) == 1
}

func (p *PPU) setVBlank() {
	log.Debugf("PPU: setting vblank")
	p.setBit(PPUSTATUS, 7, 1)
}

func (p *PPU) renderFrame() {
	if p.enableBackground() {
		p.renderBackground()

		if !p.written {
			p.written = true
			p.draw()
		}
	}
}

func (p *PPU) renderScanLine() {
}

func (p *PPU) preRenderScanline() {
}

func (p *PPU) setBit(register, bit, val byte) {
	p.registers[register] = p.registers[register] | (val << bit)
}

func (p *PPU) getBit(register, bit byte) byte {
	return p.registers[register] >> bit & 1
}

func (p *PPU) renderBackground() {
	i := baseNametables[p.registers[PPUCTRL]&0x3]

	for y := 0; y < 0x1D; y++ {
		for x := 0; x < 0x1F; x++ {
			p.renderBackgroundTile(x, y, i)
			i++
		}
	}
}

func (p *PPU) enableBackground() bool {
	return p.getBit(PPUMASK, 3) == 1
}

func (p *PPU) getPatternBase() uint16 {
	if p.getBit(PPUCTRL, 4) == 0 {
		return 0
	}

	return 0x100
}

func (p *PPU) renderBackgroundTile(x int, y int, addr uint16) {
	patternAddr := p.getPatternBase() + (uint16(p.data[addr]) * 16)

	var p1, p2 [8]byte

	for i := patternAddr; i < patternAddr+8; i++ {
		p1[i] = p.data[patternAddr]
	}

	for i := patternAddr + 8; i < patternAddr+16; i++ {
		p2[i] = p.data[patternAddr]
	}

	for i := 0; i < 8; i++ {
		p.processBackgroundLine(x, y, i, p1[i], p2[i])
	}
}

func (p *PPU) processBackgroundLine(x int, y int, line int, byte1 byte, byte2 byte) {
	for i := 0; i < 8; i++ {
		bit1 := (byte1 >> i) & 1
		bit2 := (byte2 >> i) & 1

		val := calcVal(bit1, bit2)

		pixelX := (x * 8) + (8 - i)
		pixelY := (y * 8) + line

		p.state[pixelX][pixelY] = val
	}
}

func (p *PPU) draw() {
	log.Printf("%v", p.state)
}

func calcVal(b1, b2 byte) byte {
	if b1 == b2 {
		if b1 == 0 {
			return 0
		}

		return 3
	}

	if b1 == 1 {
		return 1
	}

	return 2
}
