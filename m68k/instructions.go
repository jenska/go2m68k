package m68k

type Instruction interface {
	Execute(cpu *M68k) int
}

func (cpu *M68k) init68000InstructionSet() {
	cpu.instructions = []Instruction{}
}
