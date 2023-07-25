package cpu

import "log"

const clockSpeedHz = 1660000

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

		// CMP
		0xC9: {
			cycles:    2,
			immediate: c.cmp,
		},
		0xC5: {
			cycles:   3,
			zeroPage: c.cmp,
		},
		0xD5: {
			cycles:    4,
			zeroPageX: c.cmp,
		},
		0xCD: {
			cycles:   4,
			absolute: c.cmp,
		},
		0xDD: {
			cycles:    4,
			absoluteX: c.cmp,
		},
		0xD9: {
			cycles:    4,
			absoluteY: c.cmp,
		},
		0xC1: {
			cycles:    6,
			indirectX: c.cmp,
		},
		0xD1: {
			cycles:    5,
			indirectY: c.cmp,
		},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}

func (c *CPU) cmp(v byte) (byte, bool) {
	return c.compareGeneric(c.a, v)
}

func (c *CPU) cpx(v byte) (byte, bool) {
	return c.compareGeneric(c.x, v)
}

func (c *CPU) cpy(v byte) (byte, bool) {
	return c.compareGeneric(c.y, v)
}

func (c *CPU) compareGeneric(register byte, memory byte) (byte, bool) {
	c.setFlagTo(FlagC, register >= memory)
	c.setFlagTo(FlagZ, register == memory)
	c.setFlagTo(FlagN, isNeg(register-memory))

	return 0, false
}

// TODO: carry flag, page boundary
func (c *CPU) sbc(v byte) (byte, bool) {
	hadBorrow := c.flagSet(FlagB)

	if c.flagSet(FlagD) {
		//c.subDecimal(v, hadBorrow)
		log.Panicf("tried to use decimal sbc: %v", v)
	} else {
		c.subBinary(v, hadBorrow)
	}

	c.setNZFromA()
	return 0, false
}

// TODO: carry flag, page boundary
func (c *CPU) adc(v byte) (byte, bool) {
	hadCarry := c.flagSet(FlagC)

	if c.flagSet(FlagD) {
		//c.addDecimal(v, hadCarry)
		log.Panicf("tried to use decimal adc: %v", v)
	} else {
		c.addBinary(v, hadCarry)
	}

	c.setNZFromA()
	return 0, false
}

//func (c *CPU) addDecimal(v byte, hadCarry bool) {
//	result := fromBCD(c.a) + fromBCD(v)
//
//	if hadCarry {
//		result++
//	}
//
//	if result > 99 {
//		result %= 100
//		c.setFlag(FlagC)
//	} else {
//		c.clearFlag(FlagC)
//	}
//
//	c.a = toBCD(result)
//}

func (c *CPU) addBinary(v byte, hadCarry bool) {
	wasNeg := isNeg(c.a)

	c.setFlagTo(FlagC, uint16(c.a)+uint16(v) > 255)

	c.a += v

	if hadCarry {
		c.a++
	}

	c.setFlagTo(FlagV, wasNeg != isNeg(c.a)) // TODO: wrong
}

//func (c *CPU) subDecimal(v byte, hadBorrow bool) {
//	result := fromBCD(c.a) + fromBCD(v)
//
//	if result > 99 {
//		result %= 100
//		c.setFlag(FlagC)
//	} else {
//		c.clearFlag(FlagC)
//	}
//
//	c.a = toBCD(result)
//}

func (c *CPU) subBinary(v byte, hadBorrow bool) {
	c.a -= v

	if hadBorrow {
		c.a--
	}
}

func isNeg(v byte) bool {
	v &= 0x80
	return (v >> 7) == 0x01
}
