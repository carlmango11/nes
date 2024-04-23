package nes

import (
	"github.com/carlmango11/nes/backend/nes/bus"
	"github.com/carlmango11/nes/backend/nes/cpu"
	"github.com/carlmango11/nes/backend/nes/ppu"
	"github.com/carlmango11/nes/backend/nes/rom"
	"log"
	"time"
)

type NES struct {
	tick int
	cpu  *cpu.CPU
	ppu  *ppu.PPU
}

func New(rom rom.ROM) *NES {
	p := ppu.New(rom)
	b := bus.New(rom, p)

	return &NES{
		cpu: cpu.New(b),
		ppu: p,
	}
}

func (n *NES) Run() {
	var i int64
	//for range time.Tick(time.Second / cpu.ClockSpeedHz) {
	for {
		n.Tick()

		//log.Println("small tick")
		i++
		if i%(1660000) == 0 {
			log.Println("tick", i)
			i = 0
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (n *NES) Tick() {
	if n.tick%4 == 0 {
		interrupt := n.ppu.Tick()

		if interrupt {
			n.cpu.Interrupt()
		}
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
