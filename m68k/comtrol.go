package m68k

func registerControlInstructions(cpu *M68K) {
	cpu.registerInstruction(0x0074, ori2ccr)
	cpu.registerInstruction(0x00f4, ori2sr)

	cpu.registerInstruction(0x4e73, rte)
}

func ori2ccr(cpu *M68K) int {
	ccr := uint16(cpu.popPC(Word)&0xff) | cpu.SR.GetCCR()
	cpu.SR.SetCCR(ccr)
	return 8
}

func ori2sr(cpu *M68K) int {
	if cpu.SR.S() {
		sr := uint16(cpu.popPC(Word)) | cpu.SR.Get()
		cpu.SR.Set(sr)
		return 8
	}
	return cpu.RaiseException(XptPrv)
}

func rte(cpu *M68K) int {
	if cpu.SR.S() {
		sr := uint16(cpu.popPC(Word))
		cpu.PC = cpu.popPC(Long)
		cpu.SR.Set(sr)
		return 20
	}
	return cpu.RaiseException(XptPrv)
}
