package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressSpace_NewArea(t *testing.T) {
	mem := make([]byte, 0x10000)
	r := func(address uint32, o Operand) uint32 {
		return o.Read(mem[address:])
	}
	w := func(address uint32, o Operand, v uint32) {
		o.Write(v, mem[address:])
	}
	reset := func() {
		for i := 0; i < len(mem); i++ {
			mem[i] = 0
		}
	}

	as := &AddressSpace{table: make([]*AddressArea, 0x10)}
	tests := []struct {
		name   string
		offset uint32
		size   uint32
		reader Reader
		writer Writer
		reset  func()
		check  func()
	}{
		{"offset 0, ram area", 0, 0x10000, r, w, reset, func() {
			as.Write(0x100, Long, 2208)
			assert.Equal(t, uint32(2208), as.Read(0x100, Long))
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewArea(tt.offset, tt.size, tt.reader, tt.writer, tt.reset)
			tt.check()
		})
	}
}
