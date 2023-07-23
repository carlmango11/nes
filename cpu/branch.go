package cpu

func (c *CPU) bpl() bool {
	return !c.flagSet(FlagN)
}

func (c *CPU) bmi() bool {
	return c.flagSet(FlagN)
}

func (c *CPU) bvc() bool {
	return !c.flagSet(FlagV)
}

func (c *CPU) bvs() bool {
	return c.flagSet(FlagV)
}

func (c *CPU) bcc() bool {
	return !c.flagSet(FlagC)
}

func (c *CPU) bcs() bool {
	return c.flagSet(FlagC)
}

func (c *CPU) bne() bool {
	return !c.flagSet(FlagZ)
}

func (c *CPU) beq() bool {
	return c.flagSet(FlagZ)
}
