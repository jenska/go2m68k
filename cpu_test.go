package m68k

import (
	"testing"

	_ "github.com/jenska/m68k/internal/instructions"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	bc := NewBusController(BaseRAM(0x1000, 0xfc0000, 1024*1024), ROM(0xFC0000, nil))
	cpu := New(M68000, bc)
	assert.NotNil(t, cpu)
}

func TestReset(t *testing.T) {
	bc := NewBusController(BaseRAM(0x1000, 0xfc0000, 1024*1024), ROM(0xFC0000, nil))
	assert.Panics(t, func() {
		New(M68000, bc)
	})
	cpu := New(M68000, bc)
	assert.NotNil(t, cpu)
	signals := make(chan uint16)
	go cpu.Execute(signals)
	signals <- ResetSignal
}

func TestNewBusController(t *testing.T) {
	NewBusController(BaseRAM(0x1000, 0xfc0000, 1024*1024), ROM(0xFC0000, nil))
}
