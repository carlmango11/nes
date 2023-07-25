package cpu

func (c *CPU) initLoad() {
	instrs := map[byte]Instr{
		// LDX
		0xA2: {
			name:      "LDX",
			cycles:    2,
			immediate: c.ldx,
		},
		0xA6: {
			name:     "LDX",
			cycles:   3,
			zeroPage: c.ldx,
		},
		0xB6: {
			name:      "LDX",
			cycles:    4,
			zeroPageY: c.ldx,
		},
		0xAE: {
			name:     "LDX",
			cycles:   4,
			absolute: c.ldx,
		},
		0xBE: {
			name:      "LDX",
			cycles:    4,
			absoluteY: c.ldx,
		},

		// LDY
		0xA0: {
			name:      "LDY",
			cycles:    2,
			immediate: c.ldy,
		},
		0xA4: {
			name:     "LDY",
			cycles:   3,
			zeroPage: c.ldy,
		},
		0xB4: {
			name:      "LDY",
			cycles:    4,
			zeroPageY: c.ldy,
		},
		0xAC: {
			name:     "LDY",
			cycles:   4,
			absolute: c.ldy,
		},
		0xBC: {
			name:      "LDY",
			cycles:    4,
			absoluteY: c.ldy,
		},

		0xA9: {
			name:      "LDA",
			cycles:    2,
			immediate: c.lda,
		},
		0xAD: {
			name:     "LDA",
			cycles:   4,
			absolute: c.lda,
		},
		0xBD: {
			name:      "LDA",
			cycles:    4,
			absoluteX: c.lda,
		},
		0xB9: {
			name:      "LDA",
			cycles:    4,
			absoluteY: c.lda,
		},
		0xA5: {
			name:     "LDA",
			cycles:   3,
			zeroPage: c.lda,
		},
		0xB5: {
			name:      "LDA",
			cycles:    4,
			zeroPageX: c.lda,
		},
		0xA1: {
			name:      "LDA",
			cycles:    6,
			indirectX: c.lda,
		},
		0xB1: {
			name:      "LDA",
			cycles:    5,
			indirectY: c.lda,
		},

		// STX
		0x8E: {
			cycles:   4,
			absolute: c.stx,
		},
		0x86: {
			cycles:   3,
			zeroPage: c.stx,
		},
		0x96: {
			cycles:    4,
			zeroPageY: c.stx,
		},

		// STY
		0x8C: {
			cycles:   4,
			absolute: c.sty,
		},
		0x84: {
			cycles:   3,
			zeroPage: c.sty,
		},
		0x94: {
			cycles:    4,
			zeroPageY: c.sty,
		},

		// STA
		0x8D: {
			name:     "STA",
			cycles:   4,
			absolute: c.sta,
		},
		0x9D: {
			cycles:    5,
			absoluteX: c.sta,
		},
		0x99: {
			cycles:    5,
			absoluteY: c.sta,
		},
		0x85: {
			cycles:   3,
			zeroPage: c.sta,
		},
		0x95: {
			cycles:    4,
			zeroPageX: c.sta,
		},
		0x81: {
			cycles:    6,
			indirectX: c.sta,
		},
		0x91: {
			cycles:    6,
			indirectY: c.sta,
		},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}

func (c *CPU) ldx(v byte) (byte, bool) {
	c.x = v
	c.setNZ(c.x)

	return 0, false
}

func (c *CPU) ldy(v byte) (byte, bool) {
	c.y = v
	c.setNZ(c.y)

	return 0, false
}

func (c *CPU) lda(v byte) (byte, bool) {
	c.a = v
	c.setNZFromA()

	return 0, false
}

func (c *CPU) stx(v byte) (byte, bool) {
	return c.x, true
}

func (c *CPU) sty(v byte) (byte, bool) {
	return c.y, true
}

func (c *CPU) sta(v byte) (byte, bool) {
	return c.a, true
}
