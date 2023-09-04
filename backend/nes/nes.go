package nes

import (
	"github.com/carlmango11/nes/backend/nes/bus"
	"github.com/carlmango11/nes/backend/nes/cpu"
	"github.com/carlmango11/nes/backend/nes/ppu"
	"github.com/carlmango11/nes/backend/nes/rom"
	"time"
)

type NES struct {
	tick int
	cpu  *cpu.CPU
	ppu  *ppu.PPU
}

func New(rom rom.ROM) *NES {
	p := ppu.New()
	b := bus.New(rom, p)

	return &NES{
		cpu: cpu.New(b),
		ppu: p,
	}
}

func (n *NES) Run() {
	for range time.Tick(time.Second / cpu.ClockSpeedHz) {
		n.Tick()
		//time.Sleep(time.Second)
	}
}

func (n *NES) Tick() {
	if n.tick%4 == 0 {
		n.ppu.Tick()
	}
	if n.tick%12 == 0 {
		n.cpu.Tick()
	}

	n.tick++
	if n.tick == 12 {
		n.tick = 0
	}
}

func (n *NES) Display() [240][256]byte {
	return n.ppu.State()
}
