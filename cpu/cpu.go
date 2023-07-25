package cpu

import (
	"Nes/ram"
	"fmt"
)

type handler func(v byte) (byte, bool)
type impliedHandler func()
type addrHandler func(addr uint16)
type condition func() bool

type flagChange struct {
	flag Flag
	set  bool
}

type Instr struct {
	name   string
	cycles int

	implied      impliedHandler
	accumulator  handler
	immediate    handler
	zeroPage     handler
	zeroPageX    handler
	zeroPageY    handler
	absolute     handler
	absoluteAddr addrHandler // when an address is required
	absoluteX    handler
	absoluteY    handler
	indirect     addrHandler
	indirectX    handler
	indirectY    handler
	relative     condition

	flagChange *flagChange
}

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
	opCodes map[byte]Instr
	ram     *ram.RAM

	pc            uint16
	s, p, a, x, y byte
	c             int
}

func New(ram *ram.RAM, pc uint16) *CPU {
	c := &CPU{
		ram:     ram,
		opCodes: map[byte]Instr{},
		s:       0xFF, // starts at top
		pc:      pc,
	}

	c.initInstrs()

	return c
}

func (c *CPU) initInstrs() {
	c.initLoad()
	c.initTransfer()
	c.initStack()
	c.initShift()
	c.initLogic()
	c.initArithmetic()
	c.initIncrement()
	c.initCtrl()
	c.initBranch()
	c.initFlags()
	c.initNop()
}

func (c *CPU) PrintState() {
	fmt.Printf("\nC: %v\ta:%x x:%x y:%x s:%x pc:%x", c.c, c.a, c.x, c.y, c.s, c.pc)
	fmt.Printf("\nN: %t V: %t B: %t D: %t I: %t Z: %t C: %t",
		c.flagSet(FlagN), c.flagSet(FlagV), c.flagSet(FlagB), c.flagSet(FlagD),
		c.flagSet(FlagI), c.flagSet(FlagZ), c.flagSet(FlagC))
}

func (c *CPU) Exec() {
	code := c.read()

	instr, ok := c.opCodes[code]
	if !ok {
		panic(fmt.Sprintf("unknown opcode %x", code))
	}

	fmt.Printf("\ninstr %v - %x", instr.name, code)

	switch {
	case instr.accumulator != nil:
		c.execAccumulator(instr.accumulator)

	case instr.implied != nil:
		c.execImplied(instr.implied)

	case instr.immediate != nil:
		c.execImmediate(instr.immediate)

	case instr.zeroPage != nil:
		c.execZeroPage(instr.zeroPage)

	case instr.zeroPageX != nil:
		c.execZeroPageX(instr.zeroPageX)

	case instr.zeroPageY != nil:
		c.execZeroPageY(instr.zeroPageY)

	case instr.absolute != nil:
		c.execAbsolute(instr.absolute)

	case instr.absoluteAddr != nil:
		c.execAbsoluteAddr(instr.absoluteAddr)

	case instr.absoluteX != nil:
		c.execAbsoluteX(instr.absoluteX)

	case instr.absoluteY != nil:
		c.execAbsoluteY(instr.absoluteY)

	case instr.indirect != nil:
		c.execIndirect(instr.indirect)

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

	c.c++
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

func (c *CPU) execImplied(f impliedHandler) {
	f()
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
	c.execZeroPageGeneric(f, c.x)
}

func (c *CPU) execZeroPageY(f handler) {
	c.execZeroPageGeneric(f, c.y)
}

func (c *CPU) execZeroPageGeneric(f handler, register byte) {
	addr := uint16(c.read()+register) % 255

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

func (c *CPU) execIndirect(f addrHandler) {
	lo := uint16(c.read())
	hi := uint16(c.read())

	addr := hi << 8
	addr |= lo

	lo = uint16(c.ram.Read(addr))
	hi = uint16(c.ram.Read(addr + 1))

	addr = hi << 8
	addr |= lo

	f(addr)
}

func (c *CPU) readAddr() uint16 {
	lo := uint16(c.read())
	hi := uint16(c.read())

	addr := hi << 8
	return addr | lo
}

func (c *CPU) execIndirectX(f handler) {
	c.execIndirectGeneric(f, c.x)
}

func (c *CPU) execIndirectY(f handler) {
	c.execIndirectGeneric(f, c.y)
}

func (c *CPU) execRelative(cond condition) {
	offset := int8(c.read()) // offset is signed

	var cycles int

	if cond() {
		fmt.Printf("\nbranching: %v", offset)
		c.pc = uint16(int16(c.pc) + int16(offset))
		cycles = 1
	} else {
		fmt.Printf("\nnot branching")
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

	finalAddr := toAddr(hi, lo)

	f(c.ram.Read(finalAddr))
}

func (c *CPU) execAbsoluteAddr(f addrHandler) {
	f(c.readAddr())
}

func (c *CPU) execAbsoluteGeneric(f handler, addition byte) {
	addr := c.readAddr()
	addr += uint16(addition)

	newVal, write := f(c.ram.Read(addr))

	if write {
		c.ram.Write(addr, newVal)
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
	c.setFlagTo(FlagN, isNeg(v))
}

func (c *CPU) setNZFromA() {
	c.setNZ(c.a)
}

func (c *CPU) stackAddr() uint16 {
	return 0x0100 | uint16(c.s)
}

func (c *CPU) pushStack(v byte) {
	c.ram.Write(c.stackAddr(), v)
	c.s--
}

func (c *CPU) popStack() byte {
	v := c.ram.Read(c.stackAddr())
	c.s++

	return v
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

func toAddr(hi, lo byte) uint16 {
	addr := uint16(hi) << 8
	return addr | uint16(lo)
}
