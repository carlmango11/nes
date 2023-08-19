package rom

import (
	"Nes/log"
	"fmt"
)

// CPU $6000-$7FFF: 8 KB PRG RAM bank (optional)
// CPU $8000-$9FFF (or $C000-$DFFF): 8 KB switchable PRG ROM bank
// CPU $A000-$BFFF: 8 KB switchable PRG ROM bank
// CPU $C000-$DFFF (or $8000-$9FFF): 8 KB PRG ROM bank, fixed to the second-last bank
// CPU $E000-$FFFF: 8 KB PRG ROM bank, fixed to the last bank

// PPU $0000-$07FF (or $1000-$17FF): 2 KB switchable CHR bank
// PPU $0800-$0FFF (or $1800-$1FFF): 2 KB switchable CHR bank
// PPU $1000-$13FF (or $0000-$03FF): 1 KB switchable CHR bank
// PPU $1400-$17FF (or $0400-$07FF): 1 KB switchable CHR bank
// PPU $1800-$1BFF (or $0800-$0BFF): 1 KB switchable CHR bank
// PPU $1C00-$1FFF (or $0C00-$0FFF): 1 KB switchable CHR bank
type ROM4 struct {
	data *romData
}

type Metadata struct {
	prgSize byte
	chrSize byte
}

func newROM4(data *romData) *ROM4 {
	log.Debugf("banks prg: %v", len(data.prg))

	return &ROM4{
		data: data,
	}
}

func (r *ROM4) Read(addr uint16) byte {
	switch {
	case addr < 0xE000:
		panic(fmt.Sprintf("rom4: unhandler addr: %v", addr))
	case addr < 0xFFFF:
		// fixed to last bank
		translated := addr - 0xE000
		val := r.data.prg[len(r.data.prg)-1][translated]

		log.Debugf("reading last bank: %x (=> %x): %x", addr, translated, val)

		return val
	}

	panic(fmt.Sprintf("rom4: unhandler addr: %v", addr))
}
