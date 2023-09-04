package cpu

func (c *CPU) initStack() {
	instrs := map[byte]Instr{
		0x48: {
			name:           "PHA",
			cycles:         3,
			impliedHandler: c.pha,
			addrMode:       Implied,
		},
		0x08: {
			name:           "PHP",
			cycles:         3,
			impliedHandler: c.php,
			addrMode:       Implied,
		},
		0x68: {
			name:           "PLA",
			cycles:         4,
			impliedHandler: c.pla,
			addrMode:       Implied,
		},
		0x28: {
			name:           "PLP",
			cycles:         4,
			impliedHandler: c.plp,
			addrMode:       Implied,
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
	c.pushFlagsToStack()
}

func (c *CPU) pla() {
	c.a = c.popStack()
}

func (c *CPU) plp() {
	c.p = c.popStack()
}
