package main

import (
	_ "embed"
	"github.com/carlmango11/nes/backend/nes"
	"github.com/carlmango11/nes/backend/nes/log"
	"github.com/carlmango11/nes/backend/nes/rom"
	"time"
)

//go:embed donkey.nes
var donkeyRom []byte

//go:embed color_test.nes
var colourRom []byte

func main() {
	log.Debug = false

	n := nes.New(rom.FromBytes(donkeyRom))
	go n.Run()

	for range time.Tick(time.Second) {
		n.Display()
	}
}
