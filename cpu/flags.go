package cpu

func (c *CPU) initFlags() {
	instrs := map[byte]Instr{
		// Flags
		0x18: {
			cycles: 2,
			flagChange: &flagChange{
				flag: FlagC,
				set:  false,
			},
		},
		0x38: {
			cycles: 2,
			flagChange: &flagChange{
				flag: FlagC,
				set:  true,
			},
		},
		0x58: {
			cycles: 2,
			flagChange: &flagChange{
				flag: FlagI,
				set:  false,
			},
		},
		0x78: {
			cycles: 2,
			flagChange: &flagChange{
				flag: FlagI,
				set:  true,
			},
		},
		0xB8: {
			cycles: 2,
			flagChange: &flagChange{
				flag: FlagV,
				set:  false,
			},
		},
		0xD8: {
			name:   "CLD",
			cycles: 2,
			flagChange: &flagChange{
				flag: FlagD,
				set:  false,
			},
		},
		0xF8: {
			cycles: 2,
			flagChange: &flagChange{
				flag: FlagD,
				set:  true,
			},
		},
	}

	for code, instr := range instrs {
		c.opCodes[code] = instr
	}
}
