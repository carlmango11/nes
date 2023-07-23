package cpu

import (
	"Nes/ram"
)

type Flag byte

const (
	FlagN Flag = 7
	FlagV      = 6
	FlagB      = 4
	FlagD      = 3
	FlagI      = 2
	FlagZ      = 1
	FlagC      = 0
)

type CPU struct {
	pc uint16
	s  byte
	p  byte
	a  byte
	x  byte
	y  byte

	opCodes map[byte]Instr

	ram *ram.RAM
}

func New(ram *ram.RAM) *CPU {
	c := &CPU{
		ram: ram,
	}

	c.initInstrs()

	return c
}

func (c *CPU) exec() {
	instr := c.opCodes[c.read()]

	switch {
	case instr.accumulator != nil:
		c.execAccumulator(instr.accumulator)

	case instr.immediate != nil:
		c.execImmediate(instr.immediate)

	case instr.zeroPage != nil:
		c.execZeroPage(instr.zeroPage)

	case instr.zeroPageX != nil:
		c.execZeroPageX(instr.zeroPageX)

	case instr.absolute != nil:
		c.execAbsolute(instr.absolute)

	case instr.absoluteX != nil:
		c.execAbsoluteX(instr.absoluteX)

	case instr.absoluteY != nil:
		c.execAbsoluteY(instr.absoluteY)

	case instr.indirectX != nil:
		c.execIndirectX(instr.indirectX)

	case instr.indirectY != nil:
		c.execIndirectY(instr.indirectY)

	case instr.relative != nil:
		c.execRelative(instr.relative)

	case instr.flagChange != nil:
		c.execFlagChange(instr.flagChange)
	}

	c.addCycles(instr.cycles)
}

func (c *CPU) addCycles(count int) {

}

func (c *CPU) read() byte {
	val := c.ram.Read(c.pc)
	c.pc++

	return val
}

func (c *CPU) execFlagChange(fc *flagChange) {
	if fc.set {
		c.setFlag(fc.flag)
	} else {
		c.clearFlag(fc.flag)
	}
}

func (c *CPU) execAccumulator(f handler) {
	c.a, _ = f(c.read())
}

func (c *CPU) execImmediate(f handler) {
	f(c.read())
}

func (c *CPU) execZeroPage(f handler) {
	addr := uint16(c.read())

	newVal, write := f(c.ram.Read(addr))

	if write {
		c.ram.Write(addr, newVal)
	}
}

func (c *CPU) execZeroPageX(f handler) {
	addr := uint16(c.read()+c.x) % 255

	newVal, write := f(c.ram.Read(addr))

	if write {
		c.ram.Write(addr, newVal)
	}
}

func (c *CPU) execAbsolute(f handler) {
	c.execAbsoluteGeneric(f, 0)
}

func (c *CPU) execAbsoluteX(f handler) {
	c.execAbsoluteGeneric(f, c.x)
}

func (c *CPU) execAbsoluteY(f handler) {
	c.execAbsoluteGeneric(f, c.y)
}

func (c *CPU) execIndirectX(f handler) {
	c.execIndirectGeneric(f, c.x)
}

func (c *CPU) execIndirectY(f handler) {
	c.execIndirectGeneric(f, c.y)
}

func (c *CPU) execRelative(cond condition) {
	offset := c.read()

	var cycles int

	if cond() {
		c.pc += uint16(offset)
		cycles = 1
	} else {
		cycles = 2
	}

	c.addCycles(cycles)
	// TODO: page boundary
}

func (c *CPU) execIndirectGeneric(f handler, addition byte) {
	zeroAddr := c.read()
	zeroAddr += addition

	addr := uint16(zeroAddr)

	lo := c.ram.Read(addr)
	hi := c.ram.Read(addr + 1)

	finalAddr := (uint16(hi) << 8) | uint16(lo)

	f(c.ram.Read(finalAddr))
}

func (c *CPU) execAbsoluteGeneric(f handler, addition byte) {
	lo := uint16(c.read())
	hi := uint16(c.read())

	hi <<= 8
	hi |= lo

	hi += uint16(addition)

	newVal, write := f(c.ram.Read(hi))

	if write {
		c.ram.Write(hi, newVal)
	}
}

func (c *CPU) setFlagTo(index Flag, value bool) {
	if value {
		c.setFlag(index)
	} else {
		c.clearFlag(index)
	}
}

func (c *CPU) setFlag(index Flag) {
	v := byte(1) << index
	c.p = c.p | v
}

func (c *CPU) clearFlag(index Flag) {
	v := byte(1) << index
	c.p = c.p & ^v
}

func (c *CPU) flagSet(index Flag) bool {
	v := c.p >> index
	v &= 0x01

	return v == 1
}

func (c *CPU) setNZ(v byte) {
	c.setFlagTo(FlagZ, v == 0)
	c.setFlagTo(FlagN, int16(v) < 0)
}

func (c *CPU) setNZFromA() {
	c.setNZ(c.a)
}

func fromBCD(v byte) byte {
	dec := v & 0x0F
	return dec + (v>>4)*10
}

func toBCD(v byte) byte {
	res := v / 10
	res <<= 4

	return res | v%10
}
