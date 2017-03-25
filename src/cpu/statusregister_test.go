package cpu

import (
	"testing"
)

func buildStatusRegister() StatusRegister {
	cpu := NewM68k(1000, nil)
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

	if !sr.X && !sr.C {
		t.Error("Flag operation failed")
	}
}

func TestStatusRegister_S(t *testing.T) {
	sr := buildStatusRegister()
	cpu := sr.cpu
	sr.setS(false)
	if cpu.A[7] != cpu.USP {
		t.Error("failed to switch to user mode")
	}

	sr.setS(true)
	if cpu.A[7] != cpu.SSP {
		t.Error("failed to switch to supervisor mode")
	}
}
