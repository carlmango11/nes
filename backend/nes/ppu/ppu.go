package ppu

import (
	"fmt"
	"github.com/carlmango11/nes/backend/nes/log"
	"github.com/carlmango11/nes/backend/nes/rom"
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
	rom rom.ROM

	stateMu sync.Mutex
	state   [240][256]byte

	frame int
	line  int
	pixel int

	data []byte
	oam  [256]byte

	vramAddr  uint16
	vramRead  byte
	hiWritten bool

	registers        [8]byte
	verticalBlankNMI bool
}

func New(rom rom.ROM) *PPU {
	return &PPU{
		rom:  rom,
		line: -1,
		data: make([]byte, 16*1024),
	}
}

func (p *PPU) State() [240][256]byte {
	p.stateMu.Lock()
	defer p.stateMu.Unlock()

	return p.state
}

func (p *PPU) Tick() bool {
	var interrupt bool

	switch {
	case p.line == -1:
		p.preRenderScanline()
	case p.line == 0:
		//case r.line >= 1 && r.line <= 240:
		p.renderScanLine()
	case p.line == 1:
		p.renderFrame()
	case p.line >= 241 && p.line <= 260:
		p.registers[OAMADDR] = 0

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

func (p *PPU) WriteOAM(addr byte, v byte) {
	p.oam[addr] = v
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
	case OAMDATA:
		return p.readOAMDATA()
	case PPUDATA:
		return p.readPPUDATA()
	}

	return p.registers[addr]
}

func (p *PPU) writeRegister(addr uint16, val byte) {
	p.registers[addr] = val

	switch addr {
	case PPUCTRL:
		p.writePPUCTRL(val)
	case PPUSCROLL:
	//panic("no scroll")
	case OAMADDR:
		p.writeOAMADDR(val)
	case OAMDATA:
		p.writeOAMDATA(val)
	case PPUADDR:
		p.writePPUADDR(val)
	case PPUDATA:
		p.writePPUDATA(val)
	}
}

func (p *PPU) writeOAMADDR(val byte) {
	p.registers[OAMADDR] = val
}

func (p *PPU) writePPUADDR(val byte) {
	if !p.hiWritten {
		p.vramAddr = uint16(val) << 8
		p.hiWritten = true
	} else {
		p.vramAddr |= uint16(val)
	}
}

func (p *PPU) readOAMDATA() byte {
	addr := p.registers[OAMADDR]
	return p.oam[addr]
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

func (p *PPU) writeOAMDATA(val byte) {
	log.Debugf("writing %x to oam at %x", val, p.vramAddr)

	p.oam[p.registers[OAMADDR]] = val
	p.registers[OAMADDR]++
}

func (p *PPU) writePPUDATA(val byte) {
	log.Debugf("writing %x to vram at %x", val, p.vramAddr)

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
	//if p.enableBackground() {
	//	p.renderBackground()
	//}

	if p.enableSprites() {
		p.renderSprites()
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

func (p *PPU) renderSprites() {
	for i := 0; i < 64; i++ {
		start := i * 4
		p.renderSprite(p.oam[start : start+4])
	}
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

func (p *PPU) enableSprites() bool {
	return p.getBit(PPUMASK, 4) == 1
}

func (p *PPU) getPatternBase() uint16 {
	if p.getBit(PPUCTRL, 4) == 0 {
		return 0
	}

	return 0x100
}

func (p *PPU) renderBackgroundTile(x int, y int, addr uint16) {
	//log.Printf("background tile: %v  %v", y, x)

	patternAddr := p.getPatternBase() + (uint16(p.data[addr]) * 16)

	var p1, p2 [8]byte

	for i := 0; i < 8; i++ {
		p1[i] = p.data[patternAddr]
		patternAddr++
	}

	for i := 0; i < 8; i++ {
		p2[i] = p.data[patternAddr]
		patternAddr++
	}

	for i := 0; i < 8; i++ {
		p.processBackgroundLine(x, y, i, p1[i], p2[i])
	}
}

func (p *PPU) processBackgroundLine(x int, y int, line int, byte1 byte, byte2 byte) {
	//log.Printf("background %v %v", y, x)

	for i := 0; i < 8; i++ {
		bit1 := (byte1 >> i) & 1
		bit2 := (byte2 >> i) & 1

		val := calcVal(bit1, bit2)

		pixelX := (x * 8) + (7 - i)
		pixelY := (y * 8) + line

		if pixelX >= 256 || pixelY >= 240 {
			log.Printf("omg")
		}

		//log.Printf("background %v x %v = %v", pixelX, pixelY, val)
		p.state[pixelY][pixelX] = val
	}
}

func (p *PPU) renderSprite(sprite []byte) {
	if p.getBit(PPUCTRL, 5) == 1 {
		panic("not impl")
	}

	patternTable := uint16(0x0000)
	if p.getBit(PPUCTRL, 3) == 1 {
		patternTable = 0x1000
	}

	y := sprite[0] //+ 1 // it's offset by 1
	x := sprite[3]

	if y >= 240 {
		// hidden
		return
	}

	//attr := sprite[2]

	tileIndex := uint16(sprite[1])
	tileOffset := tileIndex * 16
	tileAddr := patternTable + tileOffset

	//log.Printf("sprite: %v (addr: %x)", sprite, tileAddr)

	for i := uint16(0); i < 8; i++ {
		p0 := p.readPatternTable(tileAddr + i)
		p1 := p.readPatternTable(tileAddr + i + 8)

		p.drawSpriteLine(x, y+byte(i), p0, p1)
	}
}

func (p *PPU) readPatternTable(addr uint16) byte {
	return p.rom.Read(addr)
}

func (p *PPU) drawSpriteLine(x, y, p0, p1 byte) {
	for i := byte(0); i < 8; i++ {
		b0 := p0 & 0x1
		b1 := p1 & 0x1

		p0 >>= 1
		p1 >>= 1

		xOffset := 7 - i
		p.state[y][x+xOffset] = calcColourIndex(b0, b1)

		log.Debugf("drew pixel %v/%v", y, x+xOffset)
	}
}

func calcColourIndex(b0, b1 byte) byte {
	if b0 == b1 {
		if b1 == 0 {
			return 0
		}

		return 3
	}

	if b0 == 1 {
		return 1
	}

	return 2
}

func calcVal(b1, b2 byte) byte {
	if b1 == b2 {
		if b1 == 0 {
			return 0
		}

		return 3
	}

	// TODO: check
	if b1 == 1 {
		return 1
	}

	return 2
}
