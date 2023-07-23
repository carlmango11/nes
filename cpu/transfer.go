package cpu

func (c *CPU) initTransfer() {
	instrs := map[byte]Instr{
		0xAA: {
			cycles:  3,
			implied: c.tax,
		},
		0xA8: {
			cycles:  3,
			implied: c.tay,
		},
		0xBA: {
			cycles:  3,
			implied: c.tsx,
		},
		0x8A: {
			cycles:  3,
			implied: c.txa,
		},
		0x9A: {
			cycles:  3,
			implied: c.txs,
		},
		0x98: {
			cycles:  3,
			implied: c.tya,
		},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}

func (c *CPU) tax() {
	c.x = c.a
	c.setNZ(c.x)
}

func (c *CPU) tay() {
	c.y = c.a
	c.setNZ(c.y)
}

func (c *CPU) tsx() {
	c.x = c.s
	c.setNZ(c.x)
}

func (c *CPU) txa() {
	c.a = c.x
	c.setNZFromA()
}

func (c *CPU) txs() {
	c.s = c.x
}

func (c *CPU) tya() {
	c.a = c.y
	c.setNZFromA()
}
