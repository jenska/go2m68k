package instructions

import (
	"testing"

	"github.com/jenska/go2m68k/internal/cpu"
)

func TestMoveQLoadsSignExtendedImmediate(t *testing.T) {
	c := cpu.New(128)
	if err := RegisterDefaults(c); err != nil {
		t.Fatalf("register defaults failed: %v", err)
	}

	if err := c.ExecuteInstruction(0x70ff); err != nil { // MOVEQ #0xFF, D0
		t.Fatalf("execute MOVEQ failed: %v", err)
	}

	regs := c.Registers()
	if regs.D[0] != 0xffffffff {
		t.Fatalf("expected D0 to contain 0xffffffff, got 0x%x", regs.D[0])
	}
}

func TestAddaAddsDataRegisterToAddressRegister(t *testing.T) {
	c := cpu.New(128)
	if err := RegisterDefaults(c); err != nil {
		t.Fatalf("register defaults failed: %v", err)
	}

	if err := c.RegisterInstruction(0x0001, func(regs *cpu.Registers, _ *cpu.Memory) error {
		regs.D[1] = 0x0003
		regs.A[2] = 0x0010
		return nil
	}); err != nil {
		t.Fatalf("setup instruction failed: %v", err)
	}

	if err := c.ExecuteInstruction(0x0001); err != nil {
		t.Fatalf("execute setup failed: %v", err)
	}

	if err := c.ExecuteInstruction(0xd4c1); err != nil { // ADDA.W D1, A2
		t.Fatalf("execute ADDA failed: %v", err)
	}

	regs := c.Registers()
	if regs.A[2] != 0x13 {
		t.Fatalf("expected A2 to contain 0x13, got 0x%x", regs.A[2])
	}
}

func TestSubaSubtractsDataRegisterFromAddressRegister(t *testing.T) {
	c := cpu.New(128)
	if err := RegisterDefaults(c); err != nil {
		t.Fatalf("register defaults failed: %v", err)
	}

	if err := c.RegisterInstruction(0x0002, func(regs *cpu.Registers, _ *cpu.Memory) error {
		regs.D[3] = 0x10
		regs.A[4] = 0x25
		return nil
	}); err != nil {
		t.Fatalf("setup instruction failed: %v", err)
	}

	if err := c.ExecuteInstruction(0x0002); err != nil {
		t.Fatalf("execute setup failed: %v", err)
	}

	if err := c.ExecuteInstruction(0x99c3); err != nil { // SUBA.L D3, A4
		t.Fatalf("execute SUBA failed: %v", err)
	}

	regs := c.Registers()
	if regs.A[4] != 0x15 {
		t.Fatalf("expected A4 to contain 0x15, got 0x%x", regs.A[4])
	}
}
