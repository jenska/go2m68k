package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	cpu := NewCPU().AttachBaseArea(0x1000, 1000, 3000).Go()

	assert.Equal(t, uint32(0x1000), cpu.readAddress(0))
	assert.Equal(t, 0x00, cpu.read(0, Byte))
	assert.Equal(t, 0x1000, cpu.read(0, Word))
	assert.Equal(t, 0x1000, cpu.read(0, Long))
	// bounds overflow
	assert.Equal(t, uint32(0x1000), cpu.readAddress(0x12000000))

	cpu.Halt <- true
}

func TestReadWrite(t *testing.T) {
	cpu := NewCPU().AttachBaseArea(2000, 1000, 3000).Go()

	assert.Panics(t, func() {
		cpu.write(0, Long, 400)
	})
	assert.Panics(t, func() {
		cpu.write(4, Long, 400)
	})
	assert.Equal(t, uint32(2000), cpu.readAddress(0))
	assert.Equal(t, uint32(1000), cpu.readAddress(4))

	cpu.write(100, Long, 3)
	assert.Equal(t, uint32(3), cpu.readAddress(100))
	assert.Equal(t, 3, cpu.read(100, Byte))
	assert.Equal(t, 3, cpu.read(100, Word))
	assert.Equal(t, 3, cpu.read(100, Long))

	cpu.Halt <- true
}

func TestReset(t *testing.T) {
	cpu := NewCPU().AttachBaseArea(2000, 1000, 3000).Go()

	cpu.Reset <- true
	cpu.trace <- func() {
		assert.True(t, cpu.SR.S)
		assert.Equal(t, uint32(1000), cpu.PC)
		assert.Equal(t, uint32(2000), cpu.SSP)
		assert.Equal(t, 0x2700, cpu.SR.ToBits())
	}

	cpu.Halt <- true
}

func TestHalt(t *testing.T) {
	cpu := NewCPU().AttachBaseArea(2000, 1000, 3000).Go()
	cpu.trace <- func() {
		assert.True(t, cpu.IsRunning)
	}
	cpu.Halt <- true
	assert.False(t, cpu.IsRunning)
}

func TestException_raiseException(t *testing.T) {
	cpu := NewCPU().AttachBaseArea(2000, 1000, 3000).Go()

	cpu.write(PrivilegeViolationError<<2, Long, 500)

	cpu.Reset <- true
	cpu.trace <- func() {
		cpu.raiseException(PrivilegeViolationError)
		assert.Equal(t, uint32(500), cpu.PC)
	}
	assert.True(t, cpu.IsRunning)
}

func TestInternalRead(t *testing.T) {
	cpu := NewCPU().AttachBaseArea(2000, 1000, 3000).Go()
	assert.NotNil(t, cpu.bus)
	cpu.bus.Read(0, Word)
	cpu.Halt <- true
}
