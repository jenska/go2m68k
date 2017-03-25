package cpu

import (
	"fmt"
)

type StatusRegister struct {
	C, V, Z, N, X, s, T bool
	Interrupts          uint16
	cpu                 *M68k
}

func NewStatusRegister(cpu *M68k) StatusRegister {
	return StatusRegister{cpu: cpu}
}

func (sr *StatusRegister) Get() uint16 {
	result := uint16(0)
	if sr.C {
		result += 1
	}
	if sr.V {
		result += 2
	}
	if sr.Z {
		result += 4
	}
	if sr.N {
		result += 8
	}
	if sr.X {
		result += 16
	}
	if sr.s {
		result += 0x2000
	}
	if sr.T {
		result += 0x8000
	}
	result += ((sr.Interrupts & 7) << 8)
	return result
}

func (sr *StatusRegister) Set(value uint16) {
	sr.C = (value & 1) != 0
	sr.V = (value & 2) != 0
	sr.Z = (value & 4) != 0
	sr.N = (value & 8) != 0
	sr.X = (value & 16) != 0
	sr.T = (value & 0x8000) != 0
	sr.Interrupts = (value & 0x0700) >> 8
	sr.setS((value & 0x2000) != 0)
}

func (sr *StatusRegister) GetCCR() uint16 {
	return sr.Get() & 0xff
}

func (sr *StatusRegister) SetCCR(value uint16) {
	sr.Set(value & 0xff)
}

func (sr *StatusRegister) S() bool {
	return sr.s
}

func (sr *StatusRegister) setS(value bool) {
	if sr.s {
		sr.cpu.SSP = sr.cpu.A[7]
	} else {
		sr.cpu.USP = sr.cpu.A[7]
	}
	sr.s = value
	if sr.s {
		sr.cpu.A[7] = sr.cpu.SSP
	} else {
		sr.cpu.A[7] = sr.cpu.USP
	}
}

func (sr StatusRegister) String() string {
	result := []byte{'-', '-', '-', '-', '-', '-', '-'}
	if sr.T {
		result[0] = 'T'
	}
	if sr.s {
		result[1] = 'S'
	}
	if sr.X {
		result[2] = 'X'
	}
	if sr.N {
		result[3] = 'N'
	}
	if sr.Z {
		result[4] = 'Z'
	}
	if sr.V {
		result[5] = 'V'
	}
	if sr.C {
		result[6] = 'C'
	}

	return fmt.Sprintf("%s-b%04b", result, sr.Interrupts&0x07)
}
