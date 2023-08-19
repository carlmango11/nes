package nes

import (
	"Nes/bus"
	"Nes/cpu"
	"Nes/rom"
	"time"
)

type NES struct {
	cpu *cpu.CPU
}

func New(rom rom.ROM) *NES {
	b := bus.New(rom)

	return &NES{
		cpu: cpu.New(b),
	}
}

func (n *NES) Run() {
	for range time.Tick(time.Second / cpu.ClockSpeedHz) {
		n.cpu.Exec()
	}
}
