package m68k

import (
	"fmt"
)

type instruction func(cpu *M68K) int

func (cpu *M68K) init68000InstructionSet() {
	cpu.instructions = make([]instruction, 0x10000)
	registerMoveInstructions(cpu)
	registerControlInstructions(cpu)
}

func (cpu *M68K) registerInstruction(opcode int, i instruction) {
	if cpu.instructions[opcode] != nil {
		panic(fmt.Errorf("failed to set opcode $%04x with %s. already used by %s", opcode, cpu.instructions[opcode], i))
	}
	cpu.instructions[opcode] = i
}
