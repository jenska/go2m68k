package cpu

import "fmt"

func (cpu *M68k) init68000InstructionSet() {
	registerMove(cpu)
}

func addInstruction(cpu *M68k, opcode uint16, instruction Instruction) {
	if cpu.instructions[opcode] != nil {
		panic(fmt.Errorf("Opcode $%04x already set"))
	}
	cpu.instructions[opcode] = instruction
}
