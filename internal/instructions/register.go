package instructions

import "github.com/jenska/go2m68k/internal/cpu"

// RegisterDefaults installs the built-in instruction set into the provided CPU.
func RegisterDefaults(c *cpu.CPU) error {
	if err := registerControl(c); err != nil {
		return err
	}
	if err := registerArithmetic(c); err != nil {
		return err
	}
	if err := registerMove(c); err != nil {
		return err
	}
	return nil
}
