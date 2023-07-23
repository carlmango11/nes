package cpu

func (c *CPU) bit(v byte) (byte, bool) {
	c.setFlagTo(FlagZ, c.a&v == 0)

	b7 := (v & 0x80) >> 7
	b6 := (v & 0x40) >> 6

	c.setFlagTo(FlagN, b7 == 1)
	c.setFlagTo(FlagV, b6 == 1)

	return 0, false
}
