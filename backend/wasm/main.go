package main

import (
	"embed"
	_ "embed"
	"github.com/carlmango11/nes/backend/nes"
	"github.com/carlmango11/nes/backend/nes/log"
	"github.com/carlmango11/nes/backend/nes/rom"
	"io"
	"syscall/js"
)

//go:embed roms/*.nes
var roms embed.FS

func main() {
	createBindings()

	waitC := make(chan bool)
	<-waitC
}

func createBindings() {
	var n *nes.NES

	log.Debug = false

	n = nes.New(getROM())
	go n.Run()

	getDisplayFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if n == nil {
			return 1
		}

		display := n.Display()

		height := len(display)
		width := len(display[0])

		result := make([]any, height*width)

		// doesn't support normal 2D typed arrays, only []any
		for y := range display {
			for x := range display[y] {
				if display[y][x] > 0 {
					log.Printf("OMG zero zero %v", display[y][x])
				}
				result[x+(y*width)] = display[y][x]
			}
		}

		return result
	})

	js.Global().Set("getDisplay", getDisplayFunc)
}

func getROM() rom.ROM {
	const romName = "roms/donkey.nes"

	f, err := roms.Open(romName)
	if err != nil {
		panic(err)
	}

	log.Printf("loaded: %v", romName)

	bytes, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return rom.FromBytes(bytes)
}
