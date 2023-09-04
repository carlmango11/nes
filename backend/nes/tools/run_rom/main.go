package main

import (
	"github.com/carlmango11/nes/backend/nes"
	"github.com/carlmango11/nes/backend/nes/rom"
	"io"
	"os"
)

func main() {
	//log.Debug = false

	f, err := os.Open("/Users/carl/IdeaProjects/Nes/backend/wasm/roms/color_test.nes")
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
