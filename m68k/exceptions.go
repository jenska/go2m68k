package m68k

type group0exception uint32

const (
	BusError               group0exception = 2
	AddressError           group0exception = 3
	IllegalOpcode                          = 4
	PrivilegeViolation                     = 8
	LineA                                  = 10
	LineF                                  = 11
	UninitializedInterrupt                 = 15
	SpuriousInterrupt                      = 2
	TrapInstruction                        = 32
)

func (cpu *M68K) raiseException(x uint32) {
	sr := uint32(cpu.SR.Get())
	cpu.SR.SetS(true)
	cpu.SR.T = false // exceptions unset the trace flag

	cpu.sync(6)
	cpu.pushSP(Long, cpu.PC)
	cpu.pushSP(Word, sr)

	address := cpu.Read(Long, x<<2)
	if address == 0 {
		// interrupt vector is uninitialised
		// raise a uninitialised interrupt vector exception instead
		address = cpu.Read(Long, UninitializedInterrupt<<2)
		if address == 0 {
			cpu.cpuHalted(cpu)
		}
	}
	cpu.sync(8)
	cpu.PC = address
}

func (cpu *M68K) raiseIterrupt(priority uint32) {
	if priority != 0 {
		priority &= 0x07
		if priority > cpu.SR.Interrupts {
			cpu.raiseException(SpuriousInterrupt + priority)
			cpu.SR.Interrupts = priority
		}
	}
}

func (cpu *M68K) raiseException0(x group0exception, faultAddress uint32) {
	if cpu.doubleFault {
		cpu.cpuHalted(cpu)
	}
	cpu.doubleFault = true
	sr := uint32(cpu.SR.Get())
	cpu.SR.SetS(true)
	cpu.SR.T = false // exceptions unset the trace flag

	cpu.sync(8)
	cpu.pushSP(Long, cpu.PC)
	cpu.pushSP(Word, sr)
	cpu.pushSP(Long, faultAddress)
	cpu.pushSP(Word, 2) // TODO stack frame

	address := cpu.Read(Long, uint32(x)<<2)
	if address == 0 {
		// interrupt vector is uninitialised
		// raise a uninitialised interrupt vector exception instead
		address = cpu.Read(Long, UninitializedInterrupt<<2)
		if address == 0 {
			cpu.cpuHalted(cpu)
		}
	}
	cpu.sync(2)
	cpu.PC = address
	cpu.doubleFault = false
	panic(x)
}
