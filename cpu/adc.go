package cpu

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
