package cpu

import (
	"Nes/bus"
	"Nes/ram"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBCD(t *testing.T) {
	tcs := map[byte]byte{
		0x21: 0x33,
		0x00: 0x00,
		0x01: 0x01,
		0x63: 0x99,
	}

	for hex, bcd := range tcs {
		assert.Equal(t, bcd, toBCD(hex))
		assert.Equal(t, hex, fromBCD(bcd))
	}
}

func TestSetFlags(t *testing.T) {
	c := New(bus.New(nil))
	c.p = 0b10101010

	c.setFlag(FlagC)
	assert.Equal(t, byte(0b10101011), c.p)

	c.clearFlag(FlagD)
	assert.Equal(t, byte(0b10100011), c.p)

	c.setFlag(FlagN)
	assert.Equal(t, byte(0b10100011), c.p)

	c.clearFlag(FlagN)
	assert.Equal(t, byte(0b00100011), c.p)
}

func TestFlagInstructions(t *testing.T) {
	r := ram.New()
	r.Write(0x00, 0x38) // set carry
	r.Write(0x01, 0x78) // set interrupt
	r.Write(0x02, 0xF8) // set decimal
	r.Write(0x03, 0x18) // clear carry
	r.Write(0x04, 0x58) // clear interrupt
	r.Write(0x05, 0xD8) // clear decimal

	c := New(bus.New(nil))

	c.Tick()
	assert.True(t, c.flagSet(FlagC))

	c.Tick()
	assert.True(t, c.flagSet(FlagI))

	c.Tick()
	assert.True(t, c.flagSet(FlagD))

	c.Tick()
	assert.False(t, c.flagSet(FlagC))

	c.Tick()
	assert.False(t, c.flagSet(FlagI))

	c.Tick()
	assert.False(t, c.flagSet(FlagD))
}
