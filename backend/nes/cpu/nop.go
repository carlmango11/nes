package cpu

func (c *CPU) initNop() {
	c.opCodes[0xEA] = Instr{
		name:   "NOP",
		cycles: 2,
	}
}
