package ram

import (
	"fmt"
)

const size = 100000000

// Stack is $0100 - $01FF
type RAM struct {
	data [size]byte
}

func New() *RAM {
	return &RAM{
		data: [size]byte{},
	}
}

func (r *RAM) Read(addr uint16) byte {
	v := r.data[addr]
	fmt.Printf("\nram: read %x from %x", v, addr)

	return v
}

func (r *RAM) Write(addr uint16, v byte) {
	fmt.Printf("\nram: write %x to %x", v, addr)
	r.data[addr] = v
}

func (r *RAM) Load(bytes []byte) {
	for i, v := range bytes {
		r.data[i] = v
	}
}
