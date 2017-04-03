package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryHandler_Mem8(t *testing.T) {
	mem := NewMemoryHandler(1000)
	mem.setMem(Byte, 0, 1)
	if v, _ := mem.Mem(Byte, 0); v != 1 {
		t.Error("failed to set byte value ")
	}
}

func TestMemoryHandler_Mem16(t *testing.T) {
	mem := NewMemoryHandler(1000)
	mem.setMem(Word, 0, 1)
	if v, _ := mem.Mem(Word, 0); v != 1 {
		t.Error("failed to set word value ")
	}
}

func TestMemoryHandler_Mem32(t *testing.T) {
	mem := NewMemoryHandler(1000)
	assert.NotNil(t, mem)
	mem.setMem(Long, 0, 1)
	if v, _ := mem.Mem(Long, 0); v != 1 {
		t.Error("failed to set long value ")
	}
}
