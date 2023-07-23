package cpu

func (c *CPU) and(v byte) (byte, bool) {
	c.a &= v

	c.setNZFromA()
	return 0, false
}
