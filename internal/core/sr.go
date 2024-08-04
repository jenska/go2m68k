package core

import (
	"strconv"
)

type StatusRegister struct {
	T1  bool   // Trace flag
	T0  bool   // Trace flag         (68020 only)
	S   bool   // Supervisor flag
	M   bool   // Master flag        (68020 only)
	X   bool   // Extend flag
	N   bool   // Negative flag
	Z   bool   // Zero flag
	V   bool   // Overflow flag
	C   bool   // Carry flag
	Ipl uint16 // Required Interrupt Priority Level
}

func (sr StatusRegister) String() string {
	appendFlag := func(flag bool, name string, str string) string {
		if flag {
			return str + name
		}
		return str + "-"
	}

	result := ""
	result = appendFlag(sr.C, "C", result)
	result = appendFlag(sr.V, "V", result)
	result = appendFlag(sr.Z, "Z", result)
	result = appendFlag(sr.N, "N", result)
	result = appendFlag(sr.X, "X", result)
	result = appendFlag(sr.S, "S", result)
	result = appendFlag(sr.T1, "T", result)
	result = appendFlag(sr.T0, "0", result)
	result = appendFlag(sr.M, "M", result)
	result += strconv.Itoa(int(sr.Ipl))
	return result
}

// Get the status register as a word
func (sr StatusRegister) Word() uint16 {
	var result uint16 = 0

	if sr.C {
		result++
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
	if sr.M {
		result += 0x1000
	}
	if sr.S {
		result += 0x2000
	}
	if sr.T0 {
		result += 0x4000
	}
	if sr.T1 {
		result += 0x8000
	}
	result += ((sr.Ipl & 7) << 8)
	return result
}

// Set the status register as a bitmap
func NewStatusRegister(value uint16) StatusRegister {
	var sr StatusRegister
	sr.C = (value & 1) != 0
	sr.V = (value & 2) != 0
	sr.Z = (value & 4) != 0
	sr.N = (value & 8) != 0
	sr.X = (value & 16) != 0
	sr.M = (value & 0x1000) != 0
	sr.S = (value & 0x2000) != 0
	sr.T0 = (value & 0x4000) != 0
	sr.T1 = (value & 0x8000) != 0
	sr.Ipl = (value & 0x0700) >> 8
	return sr
}

func (sr *StatusRegister) SetFlags(o Operand, v uint32) {
	sr.N = o.MSB(v)
	sr.Z = o.UnsignedExtend(v) == 0
	sr.V = false
	sr.C = false
}

var conditionTable = []func(sr *StatusRegister) bool{
	func(sr *StatusRegister) bool { return true },
	func(sr *StatusRegister) bool { return false },
	func(sr *StatusRegister) bool { return !sr.C && !sr.Z },
	func(sr *StatusRegister) bool { return sr.C || sr.Z },
	func(sr *StatusRegister) bool { return !sr.C },
	func(sr *StatusRegister) bool { return sr.C },
	func(sr *StatusRegister) bool { return !sr.Z },
	func(sr *StatusRegister) bool { return sr.Z },
	func(sr *StatusRegister) bool { return !sr.V },
	func(sr *StatusRegister) bool { return sr.V },
	func(sr *StatusRegister) bool { return !sr.N },
	func(sr *StatusRegister) bool { return sr.N },
	func(sr *StatusRegister) bool { return !(sr.N != sr.V) },
	func(sr *StatusRegister) bool { return sr.N != sr.V },
	func(sr *StatusRegister) bool { return (sr.N && sr.V && !sr.Z) || (!sr.N && !sr.V && !sr.Z) },
	func(sr *StatusRegister) bool { return sr.Z || (sr.N && !sr.V) || (!sr.N && sr.V) },
}

func (sr *StatusRegister) TestCC(code uint8) bool {
	return conditionTable[code&0x0f](sr)
}
