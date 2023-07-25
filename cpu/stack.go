package cpu

func (c *CPU) initStack() {
	instrs := map[byte]Instr{
		0x48: {
			name:    "PHA",
			cycles:  3,
			implied: c.pha,
		},
		0x08: {
			name:    "PHP",
			cycles:  3,
			implied: c.php,
		},
		0x68: {
			name:    "PLA",
			cycles:  4,
			implied: c.pla,
		},
		0x28: {
			name:    "PLP",
			cycles:  4,
			implied: c.plp,
		},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}

func (c *CPU) pha() {
	c.pushStack(c.a)
}

func (c *CPU) php() {
	c.pushStack(c.p)
}

func (c *CPU) pla() {
	c.a = c.popStack()
}

func (c *CPU) plp() {
	c.p = c.popStack()
}
