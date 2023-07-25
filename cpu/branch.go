package cpu

func (c *CPU) initBranch() {
	instrs := map[byte]Instr{
		// Branch
		0x10: {name: "BPL", relative: c.bpl},
		0x30: {name: "BMI", relative: c.bmi},
		0x50: {name: "BVC", relative: c.bvc},
		0x70: {name: "BVS", relative: c.bvs},
		0x90: {name: "BCC", relative: c.bcc},
		0xB0: {name: "BCS", relative: c.bcs},
		0xD0: {name: "BNE", relative: c.bne},
		0xF0: {name: "BEQ", relative: c.beq},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}

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
