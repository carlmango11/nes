package cpu

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
