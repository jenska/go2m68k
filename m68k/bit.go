package m68k

// bit manipulating instructions
func dynamicBchg(cpu *M68K, opcode uint32) {
	register := (opcode >> 9) & 7
	bit := cpu.D[register]
	var o *operand
	if ((opcode >> 3) & 7) == 0 {
		o = Long
	} else {
		o = Byte
	}
	bit &= (o.Bits - 1)
	m := cpu.loadEA(o, opcode&0x3f).compute()
	data := m.read() ^ (1 << bit)

	cpu.SR.Z = ((data & (1 << bit)) >> bit) != 0
	cpu.prefetch()
	if o == Long {
		cpu.sync(2)
		if bit > 15 {
			cpu.sync(2)
		}
	}
	m.write(data)
}

func staticBchg(cpu *M68K, opcode uint32) {
	bit := uint32(cpu.IRC)
	var o *operand
	if ((opcode >> 3) & 7) == 0 {
		o = Long
	} else {
		o = Byte
	}
	bit &= (o.Bits - 1)
	cpu.readExtensionWord()

	m := cpu.loadEA(o, opcode&0x3f).compute()
	data := m.read() ^ (1 << bit)

	cpu.SR.Z = ((data & (1 << bit)) >> bit) != 0
	cpu.prefetch()
	if o == Long {
		cpu.sync(2)
		if bit > 15 {
			cpu.sync(2)
		}
	}
	m.write(data)
}
