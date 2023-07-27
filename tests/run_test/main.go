package main

import (
	"Nes/cpu"
	"Nes/ram"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	f, err := os.Open("/Users/carl/IdeaProjects/Nes/roms/test.bin")
	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	r := ram.New()
	r.Load(bytes)

	c := cpu.New(r, 0x0400)

	for range time.Tick(time.Millisecond * 70) {
		c.Exec()
		c.PrintState()
		fmt.Println("\n--------------------")
	}
}
