package cpu

func (c *CPU) initShift() {
	instrs := map[byte]Instr{
		// ASL
		0x0A: {
			name:     "ASL",
			cycles:   2,
			handler:  c.asl,
			addrMode: Accumulator,
		},
		0x06: {
			name:     "ASL",
			cycles:   3,
			handler:  c.asl,
			addrMode: ZeroPage,
		},
		0x16: {
			name:     "ASL",
			cycles:   4,
			handler:  c.asl,
			addrMode: ZeroPageX,
		},
		0x0E: {
			name:     "ASL",
			cycles:   4,
			handler:  c.asl,
			addrMode: Absolute,
		},
		0x1E: {
			name:     "ASL",
			cycles:   4,
			handler:  c.asl,
			addrMode: AbsoluteX,
		},

		// LSR
		0x4A: {
			name:     "LSR",
			cycles:   2,
			handler:  c.lsr,
			addrMode: Accumulator,
		},
		0x4E: {
			name:     "LSR",
			cycles:   6,
			handler:  c.lsr,
			addrMode: Absolute,
		},
		0x5E: {
			name:     "LSR",
			cycles:   7,
			handler:  c.lsr,
			addrMode: AbsoluteX,
		},
		0x46: {
			name:     "LSR",
			cycles:   5,
			handler:  c.lsr,
			addrMode: ZeroPage,
		},
		0x56: {
			name:     "LSR",
			cycles:   6,
			handler:  c.lsr,
			addrMode: ZeroPageX,
		},

		// ROR
		0x6A: {
			name:     "ROR",
			cycles:   2,
			handler:  c.ror,
			addrMode: Accumulator,
		},
		0x66: {
			name:     "ROR",
			cycles:   5,
			handler:  c.ror,
			addrMode: ZeroPage,
		},
		0x76: {
			name:     "ROR",
			cycles:   6,
			handler:  c.ror,
			addrMode: ZeroPageX,
		},
		0x6E: {
			name:     "ROR",
			cycles:   6,
			handler:  c.ror,
			addrMode: Absolute,
		},
		0x7E: {
			name:     "ROR",
			cycles:   7,
			handler:  c.ror,
			addrMode: AbsoluteX,
		},

		// ROL
		0x2A: {
			name:     "ROL",
			cycles:   2,
			handler:  c.rol,
			addrMode: Accumulator,
		},
		0x26: {
			name:     "ROL",
			cycles:   5,
			handler:  c.rol,
			addrMode: ZeroPage,
		},
		0x36: {
			name:     "ROL",
			cycles:   6,
			handler:  c.rol,
			addrMode: ZeroPageX,
		},
		0x2E: {
			name:     "ROL",
			cycles:   6,
			handler:  c.rol,
			addrMode: Absolute,
		},
		0x3E: {
			name:     "ROL",
			cycles:   7,
			handler:  c.rol,
			addrMode: AbsoluteX,
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

func (c *CPU) asl(x byte) (byte, bool) {
	msb := (x & 0x80) >> 7
	c.setFlagTo(FlagC, msb == 1)

	res := x << 1

	c.setNZ(res)
	return res, true
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
