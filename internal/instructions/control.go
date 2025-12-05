package instructions

import "github.com/jenska/go2m68k/internal/cpu"

func registerControl(c *cpu.CPU) error {
	if err := c.RegisterInstruction(0x4e71, func(_ *cpu.Registers, _ *cpu.Memory) error { // NOP
		return nil
	}); err != nil {
		return err
	}

	return c.RegisterInstruction(0x4e72, func(regs *cpu.Registers, _ *cpu.Memory) error { // STOP
		regs.SR &^= 0x0003 // clear T0/T1
		return nil
	})
}
