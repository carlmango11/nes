package cpu

func (c *CPU) ldx(v byte) (byte, bool) {
	c.x = v

	return 0, false
}

func (c *CPU) ldy(v byte) (byte, bool) {
	c.y = v

	return 0, false
}
