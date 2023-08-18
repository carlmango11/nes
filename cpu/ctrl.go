package cpu

import (
	"fmt"
)

func (c *CPU) initCtrl() {
	instrs := map[byte]Instr{
		0x00: {
			cycles:         7,
			impliedHandler: c.brk,
			addrMode:       Implied,
		},
		0x4C: {
			name:        "JMP",
			cycles:      3,
			addrHandler: c.jmp,
			addrMode:    AbsoluteAddr,
		},
		0x6C: {
			name:        "JMP",
			cycles:      5,
			addrHandler: c.jmp,
			addrMode:    Indirect,
		},
		0x20: {
			cycles:      5,
			addrHandler: c.jsr,
			addrMode:    AbsoluteAddr,
		},
		0x60: {
			cycles:         6,
			impliedHandler: c.rts,
			addrMode:       Implied,
		},
		0x40: {
			cycles:         6,
			impliedHandler: c.rti,
			addrMode:       Implied,
		},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}

func (c *CPU) brk() {
	ret := c.pc + 1
	retLo := byte(ret)
	retHi := byte(ret >> 8)

	c.pushStack(retHi)
	c.pushStack(retLo)
	c.pushStack(c.p)

	lo := c.ram.Read(0xFFFE)
	hi := c.ram.Read(0xFFFF)

	targetAddr := toAddr(hi, lo)

	c.pc = targetAddr
}

func (c *CPU) rti() {
	c.p = c.popStack()
	lo := c.popStack()
	hi := c.popStack()

	c.pc = toAddr(hi, lo)
}

func (c *CPU) jmp(addr uint16) {
	fmt.Printf("\nJMP to %x", addr)
	c.pc = addr
	// TODO: do I need to implement the weird behaviour around end of page?
}

func (c *CPU) jsr(addr uint16) {
	returnAddr := c.pc - 1

	lo := byte(returnAddr)
	hi := byte(returnAddr >> 8)

	c.pushStack(hi)
	c.pushStack(lo)

	c.pc = addr
}

func (c *CPU) rts() {
	lo := c.popStack()
	hi := c.popStack()

	c.pc = toAddr(hi, lo) + 1
}
