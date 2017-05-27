package m68k

// TODO Read/Write for Long operands not cycle precise

func (cpu *M68K) fullPrefetch() {
	cpu.fullPrefetchFirstStep()
	cpu.prefetch()
}

func (cpu *M68K) fullPrefetchFirstStep() {
	cpu.IRC = uint16(cpu.Read(Word, cpu.PC))
}

func (cpu *M68K) prefetch() {
	cpu.IR = cpu.IRC
	cpu.IRC = uint16(cpu.Read(Word, cpu.PC+2))
	cpu.IRD = cpu.IR
}

func (cpu *M68K) readExtensionWord() {
	cpu.PC += Word.Size
	cpu.IRC = uint16(cpu.Read(Word, cpu.PC))
}

func (cpu *M68K) Read(o *operand, address uint32) uint32 {
	cpu.sync(2)
	if o.Size&1 == 0 && address&1 == 1 {
		cpu.raiseException0(AddressError, address)
	}

	v, err := cpu.memory.Read(o.Size, address&0xffffff)
	if err != nil {
		cpu.raiseException0(BusError, address)
	}
	cpu.sync(2)
	return v
}

func (cpu *M68K) Write(o *operand, address uint32, value uint32) {
	cpu.sync(2)
	if o.Size&1 == 0 && address&1 == 1 {
		cpu.raiseException0(AddressError, address)
	}
	if err := cpu.memory.Write(o.Size, address&0xffffff, value); err != nil {
		cpu.raiseException0(BusError, address)
	}
}

func (cpu *M68K) pushPC(o *operand, v uint32) {
	cpu.PC -= o.Size
	cpu.Write(o, cpu.PC, v)
}

func (cpu *M68K) popSP(o *operand) uint32 {
	result := cpu.Read(o, cpu.A[7])
	cpu.A[7] += o.Size
	return result
}

func (cpu *M68K) pushSP(o *operand, v uint32) {
	cpu.A[7] -= o.Size
	cpu.Write(o, cpu.A[7], v)
}
