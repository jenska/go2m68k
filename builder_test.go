package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddArea(t *testing.T) {
	builder := NewAddressBusBuilder()
	assert.Panics(t, func() {
		builder.AddArea(0, 0, NewRAMArea(100))
	})
	assert.Panics(t, func() {
		builder.AddArea(-1, 1, NewRAMArea(100))
	})
	assert.Panics(t, func() {
		builder.AddArea(0, 1, nil)
	})

	builder.AddArea(0, 0x10000, NewRAMArea(0x10000))
	assert.Panics(t, func() {
		builder.AddArea(0, 0x10000, NewRAMArea(0x10000))
	})
	builder.AddArea(0x10000, 0x10000, NewRAMArea(0x10000))
}

func TestNewAddressArea(t *testing.T) {
	ram := make([]byte, 1000)
	assert.Panics(t, func() {
		NewAddressArea(
			nil,
			func(offset int32, s *Size, value int32) {
				s.write(ram[offset:], value)
			},
			func() {
				for i := range ram {
					ram[i] = 0
				}
			},
		)
	})

	assert.Panics(t, func() {
		NewROMArea(nil)
	})

	assert.Panics(t, func() {
		NewBaseArea(0, 0, 0)
	})

}

func TestSetBus(t *testing.T) {
	builder := NewBuilder()
	assert.Panics(t, func() {
		builder.SetBus(nil)
	})

	assert.Panics(t, func() {
		builder.Build()
	})
}

func TestISA68000(t *testing.T) {
	// if tcpu.instructions[0xf300] != nil {
	// 	t.Fail()
	// }
}
