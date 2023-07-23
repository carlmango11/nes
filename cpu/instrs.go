package cpu

type handler func(v byte) (byte, bool)
type impliedHandler func()
type addrHandler func(addr uint16)
type condition func() bool

type flagChange struct {
	flag Flag
	set  bool
}

type Instr struct {
	cycles int

	implied      impliedHandler
	accumulator  handler
	immediate    handler
	zeroPage     handler
	zeroPageX    handler
	zeroPageY    handler
	absolute     handler
	absoluteAddr addrHandler // when an address is required
	absoluteX    handler
	absoluteY    handler
	indirect     addrHandler
	indirectX    handler
	indirectY    handler
	relative     condition

	flagChange *flagChange
}

func (c *CPU) initInstrs2() {
	c.opCodes = map[byte]Instr{
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

		// JMP
		0x4C: {
			cycles:       3,
			absoluteAddr: c.jmp,
		},
		0x6C: {
			cycles:   5,
			indirect: c.jmp,
		},

		// ROR
		0x6A: {
			cycles:      2,
			accumulator: c.ror,
		},
		0x66: {
			cycles:   5,
			zeroPage: c.ror,
		},
		0x76: {
			cycles:    6,
			zeroPageX: c.ror,
		},
		0x6E: {
			cycles:   6,
			absolute: c.ror,
		},
		0x7E: {
			cycles:    7,
			absoluteX: c.ror,
		},

		// ROL
		0x2A: {
			cycles:      2,
			accumulator: c.rol,
		},
		0x26: {
			cycles:   5,
			zeroPage: c.rol,
		},
		0x36: {
			cycles:    6,
			zeroPageX: c.rol,
		},
		0x2E: {
			cycles:   6,
			absolute: c.rol,
		},
		0x3E: {
			cycles:    7,
			absoluteX: c.rol,
		},

		// LDX
		0xA2: {
			cycles:    2,
			immediate: c.ldx,
		},
		0xA6: {
			cycles:   3,
			zeroPage: c.ldx,
		},
		0xB6: {
			cycles:    4,
			zeroPageY: c.ldx,
		},
		0xAE: {
			cycles:   4,
			absolute: c.ldx,
		},
		0xBE: {
			cycles:    4,
			absoluteY: c.ldx,
		},

		// LDY
		0xA0: {
			cycles:    2,
			immediate: c.ldy,
		},
		0xA4: {
			cycles:   3,
			zeroPage: c.ldy,
		},
		0xB4: {
			cycles:    4,
			zeroPageY: c.ldy,
		},
		0xAC: {
			cycles:   4,
			absolute: c.ldy,
		},
		0xBC: {
			cycles:    4,
			absoluteY: c.ldy,
		},

		// PHA
		0x48: {
			cycles:  3,
			implied: c.pha,
		},

		// TRANSFER

		// NOP
		0xEA: {
			cycles: 2,
		},
	}
}
