package cpu

func (c *CPU) initArithmetic() {
	instrs := map[byte]Instr{
		0x69: {
			cycles:    2,
			immediate: c.adc,
		},
		0x65: {
			cycles:   3,
			zeroPage: c.adc,
		},
		0x75: {
			cycles:    4,
			zeroPageX: c.adc,
		},
		0x6D: {
			cycles:   4,
			absolute: c.adc,
		},
		0x7D: {
			cycles:    4,
			absoluteX: c.adc,
		},
		0x79: {
			cycles:    4,
			absoluteY: c.adc,
		},
		0x61: {
			cycles:    6,
			indirectX: c.adc,
		},
		0x71: {
			cycles:    5,
			indirectY: c.adc,
		},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}

// TODO: carry flag, page boundary
func (c *CPU) adc(v byte) (byte, bool) {
	if c.flagSet(FlagC) {
		v++
	}

	if c.flagSet(FlagD) {
		c.addDecimal(v)
	} else {
		c.addBinary(v)
	}

	c.setNZFromA()
	return 0, false
}

func (c *CPU) addDecimal(v byte) {
	result := fromBCD(c.a) + fromBCD(v)

	if result > 99 {
		result %= 100
		c.setFlag(FlagC)
	} else {
		c.clearFlag(FlagC)
	}

	c.a = toBCD(result)
}

func (c *CPU) addBinary(v byte) {
	positive := int16(c.a) >= 0

	c.setFlagTo(FlagC, uint16(c.a)+uint16(v) > 255)

	c.a += v

	signChanged := positive != (int16(c.a) >= 0)
	c.setFlagTo(FlagV, signChanged)
}
