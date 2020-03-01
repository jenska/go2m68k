package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReset(t *testing.T) {
	cpu := NewCPU()
	page1 := ProtectPage(NewRAMPage())
	page2 := NewRAMPage()
	cpu.AttachBus(NewIOManager([]*Page{page1, page2}))
	cpu.Reset <- true
	assert.True(t, cpu.SR.S)
}

func TestHalt(t *testing.T) {
	cpu := NewCPU()
	page1 := ProtectPage(NewRAMPage())
	page2 := NewRAMPage()
	cpu.AttachBus(NewIOManager([]*Page{page1, page2}))
	cpu.Halt <- true
}

func TestException(t *testing.T) {
	cpu := NewCPU()
	page1 := ProtectPage(NewRAMPage())
	page2 := NewRAMPage()
	cpu.AttachBus(NewIOManager([]*Page{page1, page2}))
}

func TestInternalRead(t *testing.T) {
	cpu := NewCPU()
	page1 := ProtectPage(NewRAMPage())
	page2 := NewRAMPage()
	cpu.AttachBus(NewIOManager([]*Page{page1, page2}))
	assert.NotNil(t, cpu.io)
	cpu.io.Read(0, Word)
}

func TestRead(t *testing.T) {
	cpu := NewCPU()
	page1 := ProtectPage(NewRAMPage())
	page2 := NewRAMPage()
	cpu.AttachBus(NewIOManager([]*Page{page1, page2}))
	assert.Equal(t, uint32(0), cpu.readA(0))
	assert.Equal(t, 0, cpu.read(0, Byte))
	assert.Equal(t, 0, cpu.read(0, Word))
	assert.Equal(t, 0, cpu.read(0, Long))
	// bounds
	assert.Equal(t, uint32(0), cpu.readA(0x12000000))
	// empty page
}
