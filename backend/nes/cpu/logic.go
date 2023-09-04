package cpu

func (c *CPU) initLogic() {
	instrs := map[byte]Instr{
		// AND
		0x29: {
			cycles:   2,
			handler:  c.and,
			addrMode: Immediate,
		},
		0x25: {
			cycles:   3,
			handler:  c.and,
			addrMode: ZeroPage,
		},
		0x35: {
			cycles:   4,
			handler:  c.and,
			addrMode: ZeroPageX,
		},
		0x2D: {
			cycles:   4,
			handler:  c.and,
			addrMode: Absolute,
		},
		0x3D: {
			cycles:   4,
			handler:  c.and,
			addrMode: AbsoluteX,
		},
		0x39: {
			cycles:   4,
			handler:  c.and,
			addrMode: AbsoluteY,
		},
		0x21: {
			cycles:   6,
			handler:  c.and,
			addrMode: XIndirect,
		},
		0x31: {
			cycles:   5,
			handler:  c.and,
			addrMode: IndirectY,
		},

		// EOR
		0x49: {
			name:     "EOR",
			cycles:   2,
			handler:  c.eor,
			addrMode: Immediate,
		},
		0x4D: {
			name:     "EOR",
			cycles:   4,
			handler:  c.eor,
			addrMode: Absolute,
		},
		0x5D: {
			name:     "EOR",
			cycles:   4,
			handler:  c.eor,
			addrMode: AbsoluteX,
		},
		0x59: {
			name:     "EOR",
			cycles:   4,
			handler:  c.eor,
			addrMode: AbsoluteY,
		},
		0x45: {
			name:     "EOR",
			cycles:   3,
			handler:  c.eor,
			addrMode: ZeroPage,
		},
		0x55: {
			name:     "EOR",
			cycles:   4,
			handler:  c.eor,
			addrMode: ZeroPageX,
		},
		0x41: {
			name:     "EOR",
			cycles:   6,
			handler:  c.eor,
			addrMode: XIndirect,
		},
		0x51: {
			name:     "EOR",
			cycles:   5,
			handler:  c.eor,
			addrMode: IndirectY,
		},

		// ORA
		0x09: {
			name:     "ORA",
			cycles:   2,
			handler:  c.ora,
			addrMode: Immediate,
		},
		0x0D: {
			name:     "ORA",
			cycles:   4,
			handler:  c.ora,
			addrMode: Absolute,
		},
		0x1D: {
			name:     "ORA",
			cycles:   4,
			handler:  c.ora,
			addrMode: AbsoluteX,
		},
		0x19: {
			name:     "ORA",
			cycles:   4,
			handler:  c.ora,
			addrMode: AbsoluteY,
		},
		0x05: {
			name:     "ORA",
			cycles:   3,
			handler:  c.ora,
			addrMode: ZeroPage,
		},
		0x15: {
			name:     "ORA",
			cycles:   4,
			handler:  c.ora,
			addrMode: ZeroPageX,
		},
		0x01: {
			name:     "ORA",
			cycles:   6,
			handler:  c.ora,
			addrMode: XIndirect,
		},
		0x11: {
			name:     "ORA",
			cycles:   5,
			handler:  c.ora,
			addrMode: IndirectY,
		},

		// BIT
		0x24: {
			name:     "BIT",
			cycles:   2,
			handler:  c.bit,
			addrMode: ZeroPage,
		},
		0x2C: {
			name:     "BIT",
			cycles:   3,
			handler:  c.bit,
			addrMode: Absolute,
		},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}

func (c *CPU) and(v byte) (byte, bool) {
	c.a &= v

	c.setNZFromA()
	return 0, false
}

func (c *CPU) ora(v byte) (byte, bool) {
	c.a |= v

	c.setNZFromA()
	return 0, false
}

func (c *CPU) eor(v byte) (byte, bool) {
	c.a ^= v

	c.setNZFromA()
	return 0, false
}

func (c *CPU) bit(v byte) (byte, bool) {
	c.setFlagTo(FlagZ, c.a&v == 0)

	b7 := (v & 0x80) >> 7
	b6 := (v & 0x40) >> 6

	c.setFlagTo(FlagN, b7 == 1)
	c.setFlagTo(FlagV, b6 == 1)

	return 0, false
}
