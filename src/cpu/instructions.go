package cpu




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