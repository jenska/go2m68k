package cpu


func (cpu *M68k) initInstructionSet() {
	for i := 0; i<8; i++ {
		cpu.instructions[0x4840 + i] = func() int {
			return swap(cpu, i)
		}
	}

}

func swap(cpu *M68k, reg int) int {
	v := cpu.D[reg]
	vh := v >> 16
	v = (v << 16) | vh
	cpu.D[reg] = v

	sr := cpu.SR
	sr.Z, sr.N, sr.V, sr.C = v == 0, Long.isNegative(v), false, false

	return 4
}

func exg_dd(cpu *M68k, rx, ry int) int {
	cpu.D[rx], cpu.D[ry] = cpu.D[ry], cpu.D[rx]
	return 6
}