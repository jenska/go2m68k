package m68k

import (
	"fmt"
)

const (
	BusError           = 2
	AddressError       = 3
	IllegalOpcode      = 4
	PrivilegeViolation = 8
	LineA              = 10
	LineF              = 11
	None
)

type exception struct {
	xType   uint32
	message string
}

type memoryException struct {
	*exception
	faultAddress uint32
	rootCause    error
}

func (x *exception) Error() string {
	return fmt.Sprintf("exception %d", x.xType)
}

func (cpu *M68K) memoryException(xType, faultAddress uint32, rootCause error) error {
	if cpu.doubleFault {
		return &exception{0, "another group 0 exception during last one. cpu halted"}
	}

	cpu.doubleFault = true
	sr := cpu.SR.Get()

	var status uint32
	switch {
	case cpu.statusCode.read:
		status += 16
	case !cpu.statusCode.instruction:
		status += 8
	case cpu.statusCode.program:
		status += 2
	case !cpu.statusCode.program:
		status++
	case (sr & 0x2000) != 0:
		status += 4
	}

	cpu.SR.T = false
	cpu.SR.SetS(true)
	cpu.sync(8)
	cpu.pushSP(Word, status)
	cpu.pushSP(Long, faultAddress)
	cpu.pushSP(Word, uint32(cpu.IRD))
	cpu.pushSP(Word, sr)
	cpu.pushSP(Long, cpu.PC)
	cpu.sync(2)
	cpu.executeAt(BusError)
	cpu.doubleFault = false
	return &memoryException{&exception{xType, "memory exception"}, faultAddress, rootCause}
}

func (cpu *M68K) illegalException(xType uint32) error {
	sr := cpu.SR.Get()
	cpu.SR.SetS(true)
	cpu.SR.T = false

	cpu.sync(4)
	cpu.pushSP(Word, sr)
	cpu.pushSP(Long, cpu.PC)
	cpu.executeAt(xType)
	return &exception{xType, "illegal exception"}
}

func (cpu *M68K) executeAt(xType uint32) {
	cpu.PC = cpu.Read(Long, xType<<2)
	cpu.fullPrefetchFirstStep()
	cpu.sync(2)
	cpu.prefetch()
}
