package instructions

import "github.com/jenska/go2m68k/internal/cpu"

func registerMove(c *cpu.CPU) error {
	for dest := 0; dest < 8; dest++ {
		for immediate := 0; immediate <= 0xff; immediate++ {
			opcode := uint16(0x7000 | dest<<9 | immediate)
			value := uint32(int32(int8(immediate)))
			if err := c.RegisterInstruction(opcode, func(regs *cpu.Registers, _ *cpu.Memory) error {
				regs.D[dest] = value
				return nil
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
