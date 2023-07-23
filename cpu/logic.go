package cpu

func (c *CPU) initLogic() {
	instrs := map[byte]Instr{
		// AND
		0x29: {
			cycles:    2,
			immediate: c.and,
		},
		0x25: {
			cycles:   3,
			zeroPage: c.and,
		},
		0x35: {
			cycles:    4,
			zeroPageX: c.and,
		},
		0x2D: {
			cycles:   4,
			absolute: c.and,
		},
		0x3D: {
			cycles:    4,
			absoluteX: c.and,
		},
		0x39: {
			cycles:    4,
			absoluteY: c.and,
		},
		0x21: {
			cycles:    6,
			indirectX: c.and,
		},
		0x31: {
			cycles:    5,
			indirectY: c.and,
		},

		// BIT
		0x24: {
			cycles:   2,
			zeroPage: c.bit,
		},
		0x2C: {
			cycles:   3,
			absolute: c.bit,
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

func (c *CPU) bit(v byte) (byte, bool) {
	c.setFlagTo(FlagZ, c.a&v == 0)

	b7 := (v & 0x80) >> 7
	b6 := (v & 0x40) >> 6

	c.setFlagTo(FlagN, b7 == 1)
	c.setFlagTo(FlagV, b6 == 1)

	return 0, false
}
