package cpu

type handler func(v byte) (byte, bool)
type condition func() bool

type flagChange struct {
	flag Flag
	set  bool
}

type Instr struct {
	cycles int

	accumulator handler
	immediate   handler
	zeroPage    handler
	zeroPageX   handler
	absolute    handler
	absoluteX   handler
	absoluteY   handler
	indirectX   handler
	indirectY   handler
	relative    condition

	flagChange *flagChange
}

func (c *CPU) initInstrs() {
	// Move the function handler up a level?
	c.opCodes = map[byte]Instr{
		// ADC
		0x69: {
			cycles:    2,
			immediate: c.adc,
		},
		0x65: {
			cycles:   3,
			zeroPage: c.adc,
		},
		0x75: {
			cycles:    4,
			zeroPageX: c.adc,
		},
		0x6D: {
			cycles:   4,
			absolute: c.adc,
		},
		0x7D: {
			cycles:    4,
			absoluteX: c.adc,
		},
		0x79: {
			cycles:    4,
			absoluteY: c.adc,
		},
		0x61: {
			cycles:    6,
			indirectX: c.adc,
		},
		0x71: {
			cycles:    5,
			indirectY: c.adc,
		},

		// AND
		0x29: {
			cycles:    2,
			immediate: c.and,
		},
		0x25: {
			cycles:   3,
			zeroPage: c.and,
		},
		0x35: {
			cycles:    4,
			zeroPageX: c.and,
		},
		0x2D: {
			cycles:   4,
			absolute: c.and,
		},
		0x3D: {
			cycles:    4,
			absoluteX: c.and,
		},
		0x39: {
			cycles:    4,
			absoluteY: c.and,
		},
		0x21: {
			cycles:    6,
			indirectX: c.and,
		},
		0x31: {
			cycles:    5,
			indirectY: c.and,
		},

		// ASL
		0x0A: {
			cycles:      2,
			accumulator: c.asl,
		},
		0x06: {
			cycles:   3,
			zeroPage: c.asl,
		},
		0x16: {
			cycles:    4,
			zeroPageX: c.asl,
		},
		0x0E: {
			cycles:   4,
			absolute: c.asl,
		},
		0x1E: {
			cycles:    4,
			absoluteX: c.asl,
		},

		// BIT
		0x24: {
			cycles:   2,
			zeroPage: c.bit,
		},
		0x2C: {
			cycles:   3,
			absolute: c.bit,
		},

		// Branch
		0x10: {relative: c.bpl},
		0x30: {relative: c.bmi},
		0x50: {relative: c.bvc},
		0x70: {relative: c.bvs},
		0x90: {relative: c.bcc},
		0xB0: {relative: c.bcs},
		0xD0: {relative: c.bne},
		0xF0: {relative: c.beq},

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
}
