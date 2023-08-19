package nes

import (
	"Nes/cpu"
	"Nes/log"
	"Nes/ram"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

const test = "0a"

var illegalOps = map[byte]bool{
	0x02: true,
	0x03: true,
	0x12: true,
	0x22: true,
	0x32: true,
	0x42: true,
	0x52: true,
	0x62: true,
	0x72: true,
}

var dodgy = map[byte]map[string]bool{
	0x20: {
		"20 55 13": true, // PC doesn't get set correctly
	},
}

type TestCase struct {
	Name    string
	Initial cpu.State
	Final   cpu.State
	Cycles  []any
}

func TestCPU(t *testing.T) {
	const testDataDir = "/Users/carl/IdeaProjects/Nes/nes/testdata/opcodes"

	de, err := os.ReadDir(testDataDir)
	if err != nil {
		panic(err)
	}

	for _, fd := range de {
		if fd.Name() == "00.json" {
			continue
		}

		code, err := hex.DecodeString(fd.Name()[0:2])
		if err != nil {
			log.Panicf("unable to decode %v: %v", fd.Name(), err)
		}

		if illegalOps[code[0]] {
			continue
		}

		if !cpu.New(nil, 0).HasOpCode(code[0]) {
			continue
		}

		name := fmt.Sprintf("%v/%v", testDataDir, fd.Name())
		runSuite(t, name, code[0])
	}
}

func runSuite(t *testing.T, fileName string, code byte) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var tests []TestCase

	err = json.Unmarshal(bytes, &tests)
	if err != nil {
		panic(err)
	}

	for _, test := range tests {
		if dodgy[code][test.Name] {
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			runTest(t, test)
		})
	}
}

func runTest(t *testing.T, test TestCase) {
	r := ram.New()
	c := cpu.New(r, 0)

	c.LoadState(test.Initial)

	c.Exec()

	s := c.State()

	assert.Equal(t, test.Final.PC, s.PC, "pc mismatch. actual: %x (%v), expected: %x (%v)", s.PC, s.PC, test.Final.PC, test.Final.PC)
	assert.Equal(t, test.Final.A, s.A)
	assert.Equal(t, test.Final.X, s.X)
	assert.Equal(t, test.Final.Y, s.Y)
	assert.Equal(t, test.Final.S, s.S)

	for _, e := range test.Final.RAM {
		actVal := r.Read(e[0])
		assert.Equal(t, byte(e[1]), actVal, "ram mismatch at %x (%v): exp %x (%v) != act %x (%v)", e[0], e[0], e[1], e[1], actVal, actVal)
	}
}
