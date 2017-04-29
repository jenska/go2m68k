package m68k

type instruction func(cpu *M68K, opcode uint16)

func (cpu *M68K) init68000InstructionSet() {
	cpu.irqMode = AutoVectorInterrut

	cpu.instructions = make([]instruction, 0x10000)
	for i := range cpu.instructions {
		cpu.instructions[i] = illegal
	}

}

func illegal(cpu *M68K, opcode uint16) {
	switch {
	case opcode&0xa000 == 0xa000:
		cpu.illegalException(LineA)
	case opcode&0xf000 == 0xf000:
		cpu.illegalException(LineF)
	default:
		cpu.illegalException(IllegalOpcode)
	}
}
