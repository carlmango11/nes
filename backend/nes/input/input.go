package input

import (
	"fmt"
)

const (
	ButtonA      byte = 0
	ButtonB      byte = 1
	ButtonSelect byte = 2
	ButtonStart  byte = 3
	ButtonUp     byte = 4
	ButtonDown   byte = 5
	ButtonLeft   byte = 6
	ButtonRight  byte = 7
)

type Joypad struct {
	next    byte
	buttons [8]byte
}

type Input struct {
	register byte
	joypad1  *Joypad
	joypad2  *Joypad
}

func New() *Input {
	return &Input{
		joypad1: newJoypad(),
		joypad2: newJoypad(),
	}
}

func newJoypad() *Joypad {
	return &Joypad{
		buttons: [8]byte{},
	}
}

func (i *Input) Read(addr uint16) byte {
	switch addr {
	case 0x4016:
		return i.readJoypad(i.joypad1)
	case 0x4017:
		return i.readJoypad(i.joypad2)
	}

	panic(fmt.Sprintf("input: unhandled address %x", addr))
}

func (i *Input) Write(addr uint16, val byte) {
	if addr == 0x4016 {
		panic("omg")
		i.register = val
	}

	panic(fmt.Sprintf("input: unhandled address %x (%x)", addr, val))
}

func (i *Input) readJoypad(jp *Joypad) byte {
	val := jp.buttons[jp.next]

	jp.next++

	return val
}
