package cpu

func (c *CPU) initNop() {
	c.opCodes[0xEA] = Instr{
		cycles: 2,
	}
}
