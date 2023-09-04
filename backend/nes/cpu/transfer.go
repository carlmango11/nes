package cpu

func (c *CPU) initTransfer() {
	instrs := map[byte]Instr{
		0xAA: {
			name:           "TAX",
			cycles:         2,
			impliedHandler: c.tax,
			addrMode:       Implied,
		},
		0xA8: {
			name:           "TAY",
			cycles:         2,
			impliedHandler: c.tay,
			addrMode:       Implied,
		},
		0xBA: {
			name:           "TSX",
			cycles:         2,
			impliedHandler: c.tsx,
			addrMode:       Implied,
		},
		0x8A: {
			name:           "TXA",
			cycles:         2,
			impliedHandler: c.txa,
			addrMode:       Implied,
		},
		0x9A: {
			name:           "TXS",
			cycles:         2,
			impliedHandler: c.txs,
			addrMode:       Implied,
		},
		0x98: {
			name:           "TYA",
			cycles:         2,
			impliedHandler: c.tya,
			addrMode:       Implied,
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
