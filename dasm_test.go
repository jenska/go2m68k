package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisasm(t *testing.T) {
	assert.NotNil(t, opcodeInfo)
	page1 := NewRAMArea("TestRAM1", 1024*1024)
	io := NewIOManager(1024*1024, page1)
	Disassemble(0, io)
}
