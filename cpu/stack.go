package cpu

func (c *CPU) pha() {
	c.ram.Write(c.stackAddr(), c.a)
	c.s--
}

func (c *CPU) php() {
	c.ram.Write(c.stackAddr(), c.p)
	c.s--
}

// TODO: whaaaat
func (c *CPU) pla() {
	c.ram.Read(c.stackAddr())
	c.s--
}

func (c *CPU) plp() {
	c.p = c.ram.Read(c.stackAddr())
	c.s++
}
