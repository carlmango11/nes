package cpu

func (c *CPU) asl(v byte) (byte, bool) {
	msb := (v & 0x80) >> 7
	c.setFlagTo(FlagC, msb == 1)

	v <<= 1

	c.setNZ(v)
	return v, true
}
