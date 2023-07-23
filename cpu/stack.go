package cpu

func (c *CPU) initStack() {
	instrs := map[byte]Instr{
		0x48: {
			cycles:  3,
			implied: c.pha,
		},
		0x08: {
			cycles:  3,
			implied: c.php,
		},
		0x68: {
			cycles:  3,
			implied: c.pla,
		},
		0x28: {
			cycles:  3,
			implied: c.plp,
		},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}

func (c *CPU) pha() {
	c.ram.Write(c.stackAddr(), c.a)
	c.s--
}

func (c *CPU) php() {
	c.ram.Write(c.stackAddr(), c.p)
	c.s--
}

// TODO: whaaaat
func (c *CPU) pla() {
	c.ram.Read(c.stackAddr())
	c.s--
}

func (c *CPU) plp() {
	c.p = c.ram.Read(c.stackAddr())
	c.s++
}
