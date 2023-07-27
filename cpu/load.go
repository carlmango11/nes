package cpu

func (c *CPU) initLoad() {
	instrs := map[byte]Instr{
		// LDX
		0xA2: {
			name:     "LDX",
			cycles:   2,
			handler:  c.ldx,
			addrMode: Immediate,
		},
		0xA6: {
			name:     "LDX",
			cycles:   3,
			handler:  c.ldx,
			addrMode: ZeroPage,
		},
		0xB6: {
			name:     "LDX",
			cycles:   4,
			handler:  c.ldx,
			addrMode: ZeroPageY,
		},
		0xAE: {
			name:     "LDX",
			cycles:   4,
			handler:  c.ldx,
			addrMode: Absolute,
		},
		0xBE: {
			name:     "LDX",
			cycles:   4,
			handler:  c.ldx,
			addrMode: AbsoluteY,
		},

		// LDY
		0xA0: {
			name:     "LDY",
			cycles:   2,
			handler:  c.ldy,
			addrMode: Immediate,
		},
		0xA4: {
			name:     "LDY",
			cycles:   3,
			handler:  c.ldy,
			addrMode: ZeroPage,
		},
		0xB4: {
			name:     "LDY",
			cycles:   4,
			handler:  c.ldy,
			addrMode: ZeroPageY,
		},
		0xAC: {
			name:     "LDY",
			cycles:   4,
			handler:  c.ldy,
			addrMode: Absolute,
		},
		0xBC: {
			name:     "LDY",
			cycles:   4,
			handler:  c.ldy,
			addrMode: AbsoluteY,
		},

		0xA9: {
			name:     "LDA",
			cycles:   2,
			handler:  c.lda,
			addrMode: Immediate,
		},
		0xAD: {
			name:     "LDA",
			cycles:   4,
			handler:  c.lda,
			addrMode: Absolute,
		},
		0xBD: {
			name:     "LDA",
			cycles:   4,
			handler:  c.lda,
			addrMode: AbsoluteX,
		},
		0xB9: {
			name:     "LDA",
			cycles:   4,
			handler:  c.lda,
			addrMode: AbsoluteY,
		},
		0xA5: {
			name:     "LDA",
			cycles:   3,
			handler:  c.lda,
			addrMode: ZeroPage,
		},
		0xB5: {
			name:     "LDA",
			cycles:   4,
			handler:  c.lda,
			addrMode: ZeroPageX,
		},
		0xA1: {
			name:     "LDA",
			cycles:   6,
			handler:  c.lda,
			addrMode: IndirectX,
		},
		0xB1: {
			name:     "LDA",
			cycles:   5,
			handler:  c.lda,
			addrMode: IndirectY,
		},

		// STX
		0x8E: {
			cycles:   4,
			handler:  c.stx,
			addrMode: Absolute,
		},
		0x86: {
			cycles:   3,
			handler:  c.stx,
			addrMode: ZeroPage,
		},
		0x96: {
			cycles:   4,
			handler:  c.stx,
			addrMode: ZeroPageY,
		},

		// STY
		0x8C: {
			cycles:   4,
			handler:  c.sty,
			addrMode: Absolute,
		},
		0x84: {
			cycles:   3,
			handler:  c.sty,
			addrMode: ZeroPage,
		},
		0x94: {
			cycles:   4,
			handler:  c.sty,
			addrMode: ZeroPageY,
		},

		// STA
		0x8D: {
			name:     "STA",
			cycles:   4,
			handler:  c.sta,
			addrMode: Absolute,
		},
		0x9D: {
			cycles:   5,
			handler:  c.sta,
			addrMode: AbsoluteX,
		},
		0x99: {
			cycles:   5,
			handler:  c.sta,
			addrMode: AbsoluteY,
		},
		0x85: {
			cycles:   3,
			handler:  c.sta,
			addrMode: ZeroPage,
		},
		0x95: {
			cycles:   4,
			handler:  c.sta,
			addrMode: ZeroPageX,
		},
		0x81: {
			cycles:   6,
			handler:  c.sta,
			addrMode: IndirectX,
		},
		0x91: {
			cycles:   6,
			handler:  c.sta,
			addrMode: IndirectY,
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
