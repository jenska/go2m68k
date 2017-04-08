package m68k

func registerControlInstructions(cpu *M68K) {
	cpu.registerInstruction(0x0074, ori2ccr)
	cpu.registerInstruction(0x00f4, ori2sr)
	cpu.registerInstruction(0x4e75, rts)
	cpu.registerInstruction(0x4e77, rtr)
	cpu.registerInstruction(0x4e73, rte)

	for reg := 0; reg < 8; reg++ {
		cpu.registerInstruction(0x4e50+reg, func(cpu *M68K) int {
			return link(cpu, &cpu.A[reg])
		})
		cpu.registerInstruction(0x4e58+reg, func(cpu *M68K) int {
			return ulink(cpu, &cpu.A[reg])
		})
	}

	for ea := range []int{2, 5, 6, 7} {
		for reg := 0; reg < 8; reg++ {
			if ea == 7 && reg > 3 {
				break
			}
			mode := (ea << 3) + reg
			a := cpu.eahandlers[mode+Long.eaOffset]
			cpu.registerInstruction(0x4e80+mode, func(cpu *M68K) int {
				return jsr(cpu, a)
			})
		}
	}
}

func ori2ccr(cpu *M68K) int {
	ccr := (cpu.popPC(Word) & 0xff) | cpu.SR.GetCCR()
	cpu.SR.SetCCR(ccr)
	return 8
}

func ori2sr(cpu *M68K) int {
	if cpu.SR.S() {
		sr := cpu.popPC(Word) | cpu.SR.Get()
		cpu.SR.Set(sr)
		return 8
	}
	return cpu.RaiseException(XptPrv)
}

func jsr(cpu *M68K, a ea) int {
	cpu.pushSP(Long, cpu.PC)
	cpu.PC = a.compute().read()
	return a.timing() + 8
}

func rte(cpu *M68K) int {
	if cpu.SR.S() {
		sr := cpu.popSP(Word)
		cpu.PC = cpu.popSP(Long)
		cpu.SR.Set(sr)
		return 20
	}
	return cpu.RaiseException(XptPrv)
}

func rtr(cpu *M68K) int {
	cpu.SR.SetCCR(cpu.popSP(Word))
	cpu.PC = cpu.popSP(Long)
	return 20
}

func rts(cpu *M68K) int {
	cpu.PC = cpu.popSP(Long)
	return 16
}

func link(cpu *M68K, target *uint32) int {
	displacement := int(uint16(cpu.popPC(Word)))
	cpu.pushSP(Long, *target)
	*target = cpu.A[7]
	cpu.A[7] = uint32(int(cpu.A[7]) + displacement)
	return 16
}

func ulink(cpu *M68K, target *uint32) int {
	cpu.A[7] = *target
	*target = cpu.popSP(Long)
	return 12
}
