package cpu

import (
	"fmt"
	"github.com/carlmango11/nes/backend/nes/bus"
	"github.com/carlmango11/nes/backend/nes/log"
)

const ClockSpeedHz = 1660000

const (
	VectorNMI   = 0xFFFA
	VectorReset = 0xFFFC
	VectorIRQ   = 0xFFFE
)

type handler func(v byte) (byte, bool)
type impliedHandler func()
type addrHandler func(addr uint16)
type condition func() bool

type State struct {
	PC  uint16
	P   byte
	S   byte
	A   byte
	X   byte
	Y   byte
	RAM [][]uint16
}

type flagChange struct {
	flag Flag
	set  bool
}

type AddrMode string

const (
	Implied      AddrMode = "implied"
	Accumulator           = "accumulator"
	Immediate             = "immediate"
	ZeroPage              = "zeroPage"
	ZeroPageX             = "zeroPageX"
	ZeroPageY             = "zeroPageY"
	Absolute              = "absolute"
	AbsoluteAddr          = "absoluteAddr"
	AbsoluteX             = "absoluteX"
	AbsoluteY             = "absoluteY"
	Indirect              = "indirect"
	XIndirect             = "indirectX"
	IndirectY             = "indirectY"
	Relative              = "relative"
)

type Instr struct {
	name     string
	cycles   int
	addrMode AddrMode

	handler        handler
	addrHandler    addrHandler
	impliedHandler impliedHandler
	condition      condition
	flagChange     *flagChange
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
	bus     *bus.Bus
	opCodes map[byte]Instr

	pc            uint16
	s, p, a, x, y byte
	c             int
}

func New(bus *bus.Bus) *CPU {
	c := &CPU{
		bus:     bus,
		opCodes: map[byte]Instr{},
		s:       0xFF, // starts at top
	}

	c.initInstrs()
	c.vectorToPC(VectorReset)

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

func (c *CPU) Interrupt() {
	c.pushAddr(c.pc)
	c.pushStack(c.p)

	c.vectorToPC(VectorNMI)
}

func (c *CPU) HasOpCode(opCode byte) bool {
	_, ok := c.opCodes[opCode]
	return ok
}

func (c *CPU) State() *State {
	return &State{
		PC: c.pc,
		S:  c.s,
		A:  c.a,
		X:  c.x,
		Y:  c.y,
	}
}

func (c *CPU) LoadState(state State) {
	for _, e := range state.RAM {
		c.bus.Write(e[0], uint8(e[1]))
	}

	c.pc = state.PC
	c.p = state.P
	c.s = state.S
	c.a = state.A
	c.x = state.X
	c.y = state.Y
}

func (c *CPU) PrintState() {
	log.Debugf("C: %v\ta:%x x:%x y:%x s:%x pc:%x", c.c, c.a, c.x, c.y, c.s, c.pc)
	log.Debugf("N: %t V: %t B: %t D: %t I: %t Z: %t C: %t",
		c.flagSet(FlagN), c.flagSet(FlagV), c.flagSet(FlagB), c.flagSet(FlagD),
		c.flagSet(FlagI), c.flagSet(FlagZ), c.flagSet(FlagC))
}

func (c *CPU) Tick() {
	log.Debugf(">> executing at %x", c.pc)

	code := c.read()

	instr, ok := c.opCodes[code]
	if !ok {
		panic(fmt.Sprintf("unknown opcode %x", code))
	}

	log.Debugf("instr %v (%v) - %x", instr.name, instr.addrMode, code)

	if instr.flagChange != nil {
		c.execFlagChange(instr.flagChange)
		c.PrintState()
		return
	}

	switch instr.addrMode {
	case Accumulator:
		c.execAccumulator(instr.handler)

	case Implied:
		c.execImplied(instr.impliedHandler)

	case Immediate:
		c.execImmediate(instr.handler)

	case ZeroPage:
		c.execZeroPage(instr.handler)

	case ZeroPageX:
		c.execZeroPageX(instr.handler)

	case ZeroPageY:
		c.execZeroPageY(instr.handler)

	case Absolute:
		c.execAbsolute(instr.handler)

	case AbsoluteAddr:
		c.execAbsoluteAddr(instr.addrHandler)

	case AbsoluteX:
		c.execAbsoluteX(instr.handler)

	case AbsoluteY:
		c.execAbsoluteY(instr.handler)

	case Indirect:
		c.execIndirect(instr.addrHandler)

	case XIndirect:
		c.execIndirectX(instr.handler)

	case IndirectY:
		c.execIndirectY(instr.handler)

	case Relative:
		c.execRelative(instr.condition)
	}

	c.addCycles(instr.cycles)

	c.c++

	c.PrintState()
}

func (c *CPU) vectorToPC(vector uint16) {
	lo := c.bus.Read(vector)
	hi := c.bus.Read(vector + 1)
	c.pc = toAddr(hi, lo)
}

func (c *CPU) addCycles(count int) {

}

func (c *CPU) read() byte {
	val := c.bus.Read(c.pc)
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
	c.a, _ = f(c.a)
}

func (c *CPU) execImmediate(f handler) {
	f(c.read())
}

func (c *CPU) execZeroPage(f handler) {
	addr := uint16(c.read())

	newVal, write := f(c.bus.Read(addr))

	if write {
		c.bus.Write(addr, newVal)
	}
}

func (c *CPU) execZeroPageX(f handler) {
	c.execZeroPageGeneric(f, c.x)
}

func (c *CPU) execZeroPageY(f handler) {
	c.execZeroPageGeneric(f, c.y)
}

func (c *CPU) execZeroPageGeneric(f handler, register byte) {
	mem := c.read()
	addr := uint16(mem+register) % 0x100

	newVal, write := f(c.bus.Read(addr))

	if write {
		c.bus.Write(addr, newVal)
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
	lo := c.read()
	hi := c.read()

	loAddr := toAddr(hi, lo)
	hiAddr := toAddr(hi, lo+1) // add 1 to lo because we shouldn't cross page boundary

	targetLo := c.bus.Read(loAddr)
	targetHi := c.bus.Read(hiAddr)

	f(toAddr(targetHi, targetLo))
}

func (c *CPU) readAddr() uint16 {
	lo := uint16(c.read())
	hi := uint16(c.read())

	addr := hi << 8
	return addr | lo
}

func (c *CPU) execIndirectX(f handler) {
	zeroAddr := c.read()
	zeroAddr += c.x

	addr := uint16(zeroAddr)

	if addr >= 0x100 {
		panic(fmt.Sprintf("invalid zero page address %x", addr))
	}

	lo := c.bus.Read(addr)
	hi := c.bus.Read((addr + 1) % 0x100) // wrap around zero page

	finalAddr := toAddr(hi, lo)

	newVal, write := f(c.bus.Read(finalAddr))

	if write {
		c.bus.Write(finalAddr, newVal)
	}
}

func (c *CPU) execIndirectY(f handler) {
	zeroAddr := uint16(c.read())

	loAddr := c.bus.Read(zeroAddr)
	hiAddr := c.bus.Read((zeroAddr + 1) % 0x100) // wrap around zero page

	addr := toAddr(hiAddr, loAddr)
	addr += uint16(c.y)

	val := c.bus.Read(addr)

	newVal, write := f(val)

	if write {
		c.bus.Write(addr, newVal)
	}
}

func (c *CPU) execRelative(cond condition) {
	offset := int8(c.read()) // offset is signed

	var cycles int

	if cond() {
		newPC := uint16(int16(c.pc) + int16(offset))
		log.Debugf("branching to %x", newPC)

		c.pc = newPC
		cycles = 1
	} else {
		log.Debugf("not branching")
		cycles = 2
	}

	c.addCycles(cycles)
	// TODO: page boundary
}

func (c *CPU) execAbsoluteAddr(f addrHandler) {
	f(c.readAddr())
}

func (c *CPU) execAbsoluteGeneric(f handler, addition byte) {
	addr := c.readAddr()
	addr += uint16(addition)

	log.Debugf("read absolute: %x", addr)

	newVal, write := f(c.bus.Read(addr))

	if write {
		c.bus.Write(addr, newVal)
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

func (c *CPU) pushFlagsToStack() {
	c.pushStack(c.p | 0x10)
}

func (c *CPU) pushAddr(addr uint16) {
	hi, lo := fromAddr(addr)

	c.pushStack(hi)
	c.pushStack(lo)
}

func (c *CPU) pushStack(v byte) {
	c.bus.Write(c.stackAddr(), v)
	c.s--
}

func (c *CPU) popStack() byte {
	c.s++
	v := c.bus.Read(c.stackAddr())

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

func fromAddr(addr uint16) (byte, byte) {
	lo := byte(addr)
	hi := byte(addr >> 8)

	return hi, lo
}
