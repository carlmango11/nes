package main

import (
	"Nes/nes"
	"Nes/rom"
	"io"
	"os"
)

func main() {
	f, err := os.Open("/Users/carl/IdeaProjects/Nes/roms/mario3.nes")
	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	n := nes.New(rom.FromBytes(bytes))
	n.Run()
}
