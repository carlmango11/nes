package cpu

import "log"

func (c *CPU) initArithmetic() {
	instrs := map[byte]Instr{
		0x69: {
			name:     "ADC",
			cycles:   2,
			handler:  c.adc,
			addrMode: Immediate,
		},
		0x65: {
			name:     "ADC",
			cycles:   3,
			handler:  c.adc,
			addrMode: ZeroPage,
		},
		0x75: {
			name:     "ADC",
			cycles:   4,
			handler:  c.adc,
			addrMode: ZeroPageX,
		},
		0x6D: {
			name:     "ADC",
			cycles:   4,
			handler:  c.adc,
			addrMode: Absolute,
		},
		0x7D: {
			name:     "ADC",
			cycles:   4,
			handler:  c.adc,
			addrMode: AbsoluteX,
		},
		0x79: {
			name:     "ADC",
			cycles:   4,
			handler:  c.adc,
			addrMode: AbsoluteY,
		},
		0x61: {
			name:     "ADC",
			cycles:   6,
			handler:  c.adc,
			addrMode: IndirectX,
		},
		0x71: {
			name:     "ADC",
			cycles:   5,
			handler:  c.adc,
			addrMode: IndirectY,
		},

		// CMP
		0xC9: {
			name:     "CMP",
			cycles:   2,
			handler:  c.cmp,
			addrMode: Immediate,
		},
		0xC5: {
			name:     "CMP",
			cycles:   3,
			handler:  c.cmp,
			addrMode: ZeroPage,
		},
		0xD5: {
			name:     "CMP",
			cycles:   4,
			handler:  c.cmp,
			addrMode: ZeroPageX,
		},
		0xCD: {
			name:     "CMP",
			cycles:   4,
			handler:  c.cmp,
			addrMode: Absolute,
		},
		0xDD: {
			name:     "CMP",
			cycles:   4,
			handler:  c.cmp,
			addrMode: AbsoluteX,
		},
		0xD9: {
			name:     "CMP",
			cycles:   4,
			handler:  c.cmp,
			addrMode: AbsoluteY,
		},
		0xC1: {
			name:     "CMP",
			cycles:   6,
			handler:  c.cmp,
			addrMode: IndirectX,
		},
		0xD1: {
			name:     "CMP",
			cycles:   5,
			handler:  c.cmp,
			addrMode: IndirectY,
		},

		// CPX
		0xE0: {
			name:     "CPX",
			cycles:   2,
			handler:  c.cpx,
			addrMode: Immediate,
		},
		0xEC: {
			name:     "CPX",
			cycles:   3,
			handler:  c.cpx,
			addrMode: Absolute,
		},
		0xE4: {
			name:     "CPX",
			cycles:   2,
			handler:  c.cpx,
			addrMode: ZeroPage,
		},

		// CPY
		0xC0: {
			name:     "CPY",
			cycles:   2,
			handler:  c.cpy,
			addrMode: Immediate,
		},
		0xCC: {
			name:     "CPY",
			cycles:   3,
			handler:  c.cpy,
			addrMode: Absolute,
		},
		0xC4: {
			name:     "CPY",
			cycles:   2,
			handler:  c.cpy,
			addrMode: ZeroPage,
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
