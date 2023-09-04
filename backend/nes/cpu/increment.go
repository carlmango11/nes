package cpu

func (c *CPU) initIncrement() {
	instrs := map[byte]Instr{
		// INC
		0xEE: {
			cycles:   3,
			handler:  c.inc,
			addrMode: Absolute,
		},
		0xFE: {
			cycles:   3,
			handler:  c.inc,
			addrMode: AbsoluteX,
		},
		0xE6: {
			cycles:   2,
			handler:  c.inc,
			addrMode: ZeroPage,
		},
		0xF6: {
			cycles:   2,
			handler:  c.inc,
			addrMode: ZeroPageX,
		},
		0xE8: {
			cycles:         1,
			impliedHandler: c.inx,
			addrMode:       Implied,
		},
		0xC8: {
			cycles:         1,
			impliedHandler: c.iny,
			addrMode:       Implied,
		},

		0xCE: {
			name:     "DEC",
			cycles:   3,
			handler:  c.dec,
			addrMode: Absolute,
		},
		0xDE: {
			name:     "DEC",
			cycles:   3,
			handler:  c.dec,
			addrMode: AbsoluteX,
		},
		0xC6: {
			name:     "DEC",
			cycles:   2,
			handler:  c.dec,
			addrMode: ZeroPage,
		},
		0xD6: {
			name:     "DEC",
			cycles:   2,
			handler:  c.dec,
			addrMode: ZeroPageX,
		},
		0xCA: {
			name:           "DEX",
			cycles:         1,
			impliedHandler: c.dex,
			addrMode:       Implied,
		},
		0x88: {
			name:           "DEY",
			cycles:         1,
			impliedHandler: c.dey,
			addrMode:       Implied,
		},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}

func (c *CPU) dex() {
	c.x--
	c.setNZ(c.x)
}

func (c *CPU) dey() {
	c.y--
	c.setNZ(c.y)
}

func (c *CPU) dec(v byte) (byte, bool) {
	v--
	c.setNZ(v)

	return v, true
}

func (c *CPU) inx() {
	c.x++
	c.setNZ(c.x)
}

func (c *CPU) iny() {
	c.y++
	c.setNZ(c.y)
}

func (c *CPU) inc(v byte) (byte, bool) {
	v++
	c.setNZ(v)

	return v, true
}
