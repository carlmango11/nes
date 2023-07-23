package cpu

func (c *CPU) initLoad() {
	instrs := map[byte]Instr{
		// LDX
		0xA2: {
			cycles:    2,
			immediate: c.ldx,
		},
		0xA6: {
			cycles:   3,
			zeroPage: c.ldx,
		},
		0xB6: {
			cycles:    4,
			zeroPageY: c.ldx,
		},
		0xAE: {
			cycles:   4,
			absolute: c.ldx,
		},
		0xBE: {
			cycles:    4,
			absoluteY: c.ldx,
		},

		// LDY
		0xA0: {
			cycles:    2,
			immediate: c.ldy,
		},
		0xA4: {
			cycles:   3,
			zeroPage: c.ldy,
		},
		0xB4: {
			cycles:    4,
			zeroPageY: c.ldy,
		},
		0xAC: {
			cycles:   4,
			absolute: c.ldy,
		},
		0xBC: {
			cycles:    4,
			absoluteY: c.ldy,
		},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}

func (c *CPU) ldx(v byte) (byte, bool) {
	c.x = v

	return 0, false
}

func (c *CPU) ldy(v byte) (byte, bool) {
	c.y = v

	return 0, false
}
