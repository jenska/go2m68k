package m68k

import (
	"testing"

	"github.com/jenska/atari2go/mem"
)

func buildStatusRegister() StatusRegister {
	cpu := NewM68k(mem.NewMemoryHandler(1000))
	return newStatusRegister(cpu)
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
	sr.SetS(false)
	if cpu.A[7] != cpu.USP {
		t.Error("failed to switch to user mode")
	}

	sr.SetS(true)
	if cpu.A[7] != cpu.SSP {
		t.Error("failed to switch to supervisor mode")
	}
}
