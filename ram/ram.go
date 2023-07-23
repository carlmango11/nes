package ram

const size = 100

type RAM struct {
	data [size]byte
}

func New() *RAM {
	return &RAM{
		data: [size]byte{},
	}
}

func (r *RAM) Read(addr uint16) byte {
	return r.data[addr]
}

func (r *RAM) Write(addr uint16, v byte) {
	r.data[addr] = v
}
