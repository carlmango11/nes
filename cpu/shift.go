package cpu

func (c *CPU) initShift() {
	instrs := map[byte]Instr{
		// ASL
		0x0A: {
			cycles:      2,
			accumulator: c.asl,
		},
		0x06: {
			cycles:   3,
			zeroPage: c.asl,
		},
		0x16: {
			cycles:    4,
			zeroPageX: c.asl,
		},
		0x0E: {
			cycles:   4,
			absolute: c.asl,
		},
		0x1E: {
			cycles:    4,
			absoluteX: c.asl,
		},

		// LSR
		0x4A: {
			cycles:      2,
			accumulator: c.lsr,
		},
		0x4E: {
			cycles:   6,
			absolute: c.lsr,
		},
		0x5E: {
			cycles:    7,
			absoluteX: c.lsr,
		},
		0x46: {
			cycles:   5,
			zeroPage: c.lsr,
		},
		0x56: {
			cycles:    6,
			zeroPageX: c.lsr,
		},

		// ROR
		0x6A: {
			cycles:      2,
			accumulator: c.ror,
		},
		0x66: {
			cycles:   5,
			zeroPage: c.ror,
		},
		0x76: {
			cycles:    6,
			zeroPageX: c.ror,
		},
		0x6E: {
			cycles:   6,
			absolute: c.ror,
		},
		0x7E: {
			cycles:    7,
			absoluteX: c.ror,
		},

		// ROL
		0x2A: {
			cycles:      2,
			accumulator: c.rol,
		},
		0x26: {
			cycles:   5,
			zeroPage: c.rol,
		},
		0x36: {
			cycles:    6,
			zeroPageX: c.rol,
		},
		0x2E: {
			cycles:   6,
			absolute: c.rol,
		},
		0x3E: {
			cycles:    7,
			absoluteX: c.rol,
		},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}

func (c *CPU) lsr(v byte) (byte, bool) {
	v >>= 1

	c.setNZ(v)
	return v, true
}

func (c *CPU) asl(v byte) (byte, bool) {
	msb := (v & 0x80) >> 7
	c.setFlagTo(FlagC, msb == 1)

	v <<= 1

	c.setNZ(v)
	return v, true
}

func (c *CPU) ror(v byte) (byte, bool) {
	var left byte
	if c.flagSet(FlagC) {
		left = 1
	}

	left <<= 7

	c.setFlagTo(FlagC, (v&0x01) == 1)

	v >>= 1
	v |= left

	return v, true
}

func (c *CPU) rol(v byte) (byte, bool) {
	var right byte
	if c.flagSet(FlagC) {
		right = 1
	}

	c.setFlagTo(FlagC, (v&0x80) == 1)

	v <<= 1
	v |= right

	return v, true
}
