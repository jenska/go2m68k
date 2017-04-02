package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryHandler_Mem8(t *testing.T) {
	mem := NewMemoryHandler(1000)
	mem.setMem(Byte,0, 1)
	if v, _ := mem.Mem(Byte,0); v != 1 {
		t.Error("failed to set byte value ")
	}
}

func TestMemoryHandler_Mem16(t *testing.T) {
	mem := NewMemoryHandler(1000)
	mem.setMem(Word,0, 1)
	if v, _ := mem.Mem(Word,0); v != 1 {
		t.Error("failed to set word value ")
	}
}

func TestMemoryHandler_Mem32(t *testing.T) {
	mem := NewMemoryHandler(1000)
	mem.setMem(Long,0, 1)
	if v, _ := mem.Mem(Long,0); v != 1 {
		t.Error("failed to set long value ")
	}
}

type TestMemoryController uint32

func (TestMemoryController) Mem(o *Operand, a uint32) (v uint32, ok bool) {
	return 123, true
}

func (TestMemoryController) setMem(o *Operand, a uint32, v uint32) bool {
	panic("implement me")
}

type TestSystemController uint32

func (TestSystemController) Mem(o *Operand, a uint32) (v uint32, ok bool) {
	return 123, true
}

func (TestSystemController) setMem(o *Operand, a uint32, v uint32) bool {
	panic("implement me")
}

func TestMemoryHandler_RegisterChipset(t *testing.T) {
	mem := NewMemoryHandler(1000)
	a := []uint32{0xffff8001}
	mem.RegisterChipset(a, TestMemoryController(0xffff8001))
	a = []uint32{0xffff8006}
	mem.RegisterChipset(a, TestSystemController(0xffff8006))

	if r, ok := mem.Mem(Byte,0xffff8001); ok == true {
		assert.Equal(t, uint32(123), r)
	} else {
		t.Error("failed to access chipset")
	}

	if r, ok := mem.Mem(Word,0xffff8006); ok == true {
		assert.Equal(t, uint32(123), r)
	} else {
		t.Error("failed to access chipset")
	}

}
