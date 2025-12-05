package cpu

import "testing"

func TestStepExecutesNOPFromMemory(t *testing.T) {
	c := New(1024)
	if err := c.LoadProgram(0x100, []byte{0x4e, 0x71}); err != nil {
		t.Fatalf("failed to load program: %v", err)
	}
	c.Reset(0x100)

	if err := c.Step(); err != nil {
		t.Fatalf("step failed: %v", err)
	}
	regs := c.Registers()
	if regs.PC != 0x102 {
		t.Fatalf("expected PC to advance to 0x102, got 0x%x", regs.PC)
	}
}

func TestRegisterInstructionUpdatesRegisters(t *testing.T) {
	c := New(256)
	if err := c.RegisterInstruction(0x7001, func(regs *Registers, _ *Memory) error {
		regs.D[0] = 1
		return nil
	}); err != nil {
		t.Fatalf("register instruction failed: %v", err)
	}

	if err := c.ExecuteInstruction(0x7001); err != nil {
		t.Fatalf("execute instruction failed: %v", err)
	}

	regs := c.Registers()
	if regs.D[0] != 1 {
		t.Fatalf("expected D0 to be set to 1, got %d", regs.D[0])
	}
}

func TestUnknownInstructionReturnsError(t *testing.T) {
	c := New(256)
	if err := c.ExecuteInstruction(0xffff); err == nil {
		t.Fatal("expected error for unknown opcode")
	}
}

func TestMemoryWordAccess(t *testing.T) {
	c := New(16)
	if err := c.memory.WriteWord(0, 0x1234); err != nil {
		t.Fatalf("write word failed: %v", err)
	}
	value, err := c.memory.ReadWord(0)
	if err != nil {
		t.Fatalf("read word failed: %v", err)
	}
	if value != 0x1234 {
		t.Fatalf("expected 0x1234, got 0x%x", value)
	}
}
