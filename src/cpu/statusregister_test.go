package cpu

import (
	"testing"
	"mem"
)


func buildStatusRegister() StatusRegister {
	mem := mem.NewPhysicalAddressSpace(1000)
	cpu := NewM68k(mem)
	return NewStatusRegister(cpu)
}

func TestNewStatusRegister(t *testing.T) {
	sr := buildStatusRegister()
	if sr.Get() != 0 {
		t.Error("sr must be 0 when initialized")
	}

}

func TestStatusRegister_Get(t *testing.T) {
	sr := buildStatusRegister()
	sr.C = true
	sr.X = true

	bitmap := sr.Get()
	sr.Set(bitmap)

	if( !sr.X && !sr.C ) {
		t.Error("Flag operation failed")
	}
}