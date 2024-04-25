package main

import (
	"fmt"
	"github.com/carlmango11/nes/backend/nes/rom"
	"io"
	"os"
)

func main() {
	f, err := os.Open("/Users/carl/IdeaProjects/Nes/backend/wasm/roms/donkey.nes")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	r3 := rom.FromBytes(bytes)

	start := uint16(0x12d0)

	//for i := start; i < 1; i++ {
	spriteBytes := readData(r3, start, 16)

	s := readSprite(spriteBytes)
	printSprite(s)

	fmt.Printf("\n---\n")
	//}
}

func readData(r rom.ROM, from uint16, c int) []byte {
	var data []byte

	for i := 0; i < c; i++ {
		data = append(data, r.Read(from))
		from++
	}

	return data
}

func printSprite(s [8][8]byte) {
	for y := range s {
		for x := range s[y] {
			switch s[y][x] {
			case 0:
				fmt.Printf("⬜")
			case 1:
				fmt.Printf("🟩")
			case 2:
				fmt.Printf("🟥")
			case 3:
				fmt.Printf("🟧")
			}
		}
		fmt.Printf("\n")
	}
}

func readSprite(b []byte) [8][8]byte {
	state := [8][8]byte{}

	for row := uint16(0); row < 8; row++ {
		p0 := b[row]
		p1 := b[row+8]

		for col := byte(0); col < 8; col++ {
			b0 := p0 & 0x1
			b1 := p1 & 0x1

			p0 >>= 1
			p1 >>= 1

			xOffset := 7 - col
			v := calcColourIndex(b0, b1)
			state[row][xOffset] = v

			//fmt.Printf("\nwrote %v to %v/%v", v, row, xOffset)
		}
	}

	return state
}

func calcColourIndex(b0, b1 byte) byte {
	if b0 == b1 {
		if b1 == 0 {
			return 0
		}

		return 3
	}

	if b0 == 1 {
		return 1
	}

	return 2
}
