package cpu

import "testing"

func TestMemoryHandler_Mem8(t *testing.T) {
	mem := NewMemoryHandler(1000, nil)
	mem.setMem8(0, 1)
	if v, _ := mem.Mem8(0); v != 1 {
		t.Error("failed to set byte value ")
	}
}

func TestMemoryHandler_Mem16(t *testing.T) {
	mem := NewMemoryHandler(1000, nil)
	mem.setMem16(0, 1)
	if v, _ := mem.Mem16(0); v != 1 {
		t.Error("failed to set word value ")
	}
}

func TestMemoryHandler_Mem32(t *testing.T) {
	mem := NewMemoryHandler(1000, nil)
	mem.setMem32(0, 1)
	if v, _ := mem.Mem32(0); v != 1 {
		t.Error("failed to set long value ")
	}
}
