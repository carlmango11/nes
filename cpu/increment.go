package cpu

func (c *CPU) initIncrement() {
	instrs := map[byte]Instr{
		// INC
		0xEE: {
			cycles:   3,
			absolute: c.inc,
		},
		0xFE: {
			cycles:    3,
			absoluteX: c.inc,
		},
		0xE6: {
			cycles:   2,
			zeroPage: c.inc,
		},
		0xF6: {
			cycles:    2,
			zeroPageX: c.inc,
		},
		0xE8: {
			cycles:  1,
			implied: c.inx,
		},
		0xC8: {
			cycles:  1,
			implied: c.iny,
		},

		0xCE: {
			name:     "DEC",
			cycles:   3,
			absolute: c.dec,
		},
		0xDE: {
			name:      "DEC",
			cycles:    3,
			absoluteX: c.dec,
		},
		0xC6: {
			name:     "DEC",
			cycles:   2,
			zeroPage: c.dec,
		},
		0xD6: {
			name:      "DEC",
			cycles:    2,
			zeroPageX: c.dec,
		},
		0xCA: {
			name:    "DEX",
			cycles:  1,
			implied: c.dex,
		},
		0x88: {
			name:    "DEY",
			cycles:  1,
			implied: c.dey,
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
