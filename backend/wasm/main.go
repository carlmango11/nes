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
		for x := range display {
			for y := range display[x] {
				result[y+(x*width)] = display[x][y]
			}
		}

		return result
	})

	js.Global().Set("getDisplay", getDisplayFunc)
}

func getROM() rom.ROM {
	f, err := roms.Open("roms/donkey.nes")
	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return rom.FromBytes(bytes)
}
