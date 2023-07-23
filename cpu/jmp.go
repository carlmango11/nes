package cpu

func (c *CPU) jmp(addr uint16) {
	c.pc = addr

	// TODO: do I need to implement the weird behaviour around end of page?
}
