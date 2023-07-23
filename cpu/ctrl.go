package cpu

func (c *CPU) initCtrl() {
	instrs := map[byte]Instr{
		// JMP
		0x4C: {
			cycles:       3,
			absoluteAddr: c.jmp,
		},
		0x6C: {
			cycles:   5,
			indirect: c.jmp,
		},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}

func (c *CPU) jmp(addr uint16) {
	c.pc = addr
	// TODO: do I need to implement the weird behaviour around end of page?
}
