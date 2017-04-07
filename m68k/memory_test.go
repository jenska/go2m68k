package m68k

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryHandler_Mem8(t *testing.T) {
	mem := NewMemoryHandler(1000)
	mem.Write(Byte, 0, 1)
	if v, _ := mem.Read(Byte, 0); v != 1 {
		t.Error("failed to set byte value ")
	}
}

func TestMemoryHandler_Mem16(t *testing.T) {
	mem := NewMemoryHandler(1000)
	mem.Write(Word, 0, 1)
	if v, _ := mem.Read(Word, 0); v != 1 {
		t.Error("failed to set word value ")
	}
}

func TestMemoryHandler_Mem32(t *testing.T) {
	mem := NewMemoryHandler(1000)
	assert.NotNil(t, mem)
	mem.Write(Long, 0, 1)
	if v, _ := mem.Read(Long, 0); v != 1 {
		t.Error("failed to set long value ")
	}
}
