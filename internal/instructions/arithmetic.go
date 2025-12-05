package instructions

import "github.com/jenska/go2m68k/internal/cpu"

func registerArithmetic(c *cpu.CPU) error {
	for dest := 0; dest < 8; dest++ {
		for src := 0; src < 8; src++ {
			if err := registerAdda(c, dest, src); err != nil {
				return err
			}
			if err := registerSuba(c, dest, src); err != nil {
				return err
			}
		}
	}
	return nil
}

func registerAdda(c *cpu.CPU, dest, src int) error {
	wordOpcode := uint16(0xd0c0 | dest<<9 | src)
	if err := c.RegisterInstruction(wordOpcode, func(regs *cpu.Registers, _ *cpu.Memory) error {
		value := uint32(int16(regs.D[src] & 0xffff))
		regs.A[dest] = uint32(int32(regs.A[dest]) + int32(value))
		return nil
	}); err != nil {
		return err
	}

	longOpcode := uint16(0xd1c0 | dest<<9 | src)
	return c.RegisterInstruction(longOpcode, func(regs *cpu.Registers, _ *cpu.Memory) error {
		regs.A[dest] = regs.A[dest] + uint32(regs.D[src])
		return nil
	})
}

func registerSuba(c *cpu.CPU, dest, src int) error {
	wordOpcode := uint16(0x90c0 | dest<<9 | src)
	if err := c.RegisterInstruction(wordOpcode, func(regs *cpu.Registers, _ *cpu.Memory) error {
		value := uint32(int16(regs.D[src] & 0xffff))
		regs.A[dest] = uint32(int32(regs.A[dest]) - int32(value))
		return nil
	}); err != nil {
		return err
	}

	longOpcode := uint16(0x91c0 | dest<<9 | src)
	return c.RegisterInstruction(longOpcode, func(regs *cpu.Registers, _ *cpu.Memory) error {
		regs.A[dest] = regs.A[dest] - uint32(regs.D[src])
		return nil
	})
}
