package cpu

import (
	"fmt"
)

func (c *CPU) initCtrl() {
	instrs := map[byte]Instr{
		0x00: {
			cycles:  7,
			implied: c.brk,
		},
		0x4C: {
			name:         "JMP",
			cycles:       3,
			absoluteAddr: c.jmp,
		},
		0x6C: {
			name:     "JMP",
			cycles:   5,
			indirect: c.jmp,
		},
		0x20: {
			cycles:       5,
			absoluteAddr: c.jsr,
		},
		0x60: {
			cycles:  6,
			implied: c.rts,
		},
		0x40: {
			cycles:  6,
			implied: c.rti,
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
	lo := byte(c.pc)
	hi := byte(c.pc >> 8)

	c.pushStack(hi)
	c.pushStack(lo)

	c.pc = addr
}

func (c *CPU) rts() {
	lo := c.popStack()
	hi := c.popStack()

	c.pc = toAddr(hi, lo)
}
