package cpu

const (
	flagLogical = iota
	flagCmp
	flagAdd
	flagSub
	flagAddx
	flagSubx
	flagZn
)

// StatusRegister for M68000 cpu
type StatusRegister struct {
	C, V, Z, N, X, S, T1, T0, M bool
	Interrupts                  int
}

// Get the status register as a bitmap
func (sr *StatusRegister) Get() int {
	result := 0
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
	result += ((sr.Interrupts & 7) << 8)
	return result
}

// Set the status register as a bitmap
func (sr *StatusRegister) Set(value int) {
	sr.C = (value & 1) != 0
	sr.V = (value & 2) != 0
	sr.Z = (value & 4) != 0
	sr.N = (value & 8) != 0
	sr.X = (value & 16) != 0
	sr.M = (value & 0x1000) != 0
	sr.S = (value & 0x2000) != 0
	sr.T0 = (value & 0x4000) != 0
	sr.T1 = (value & 0x8000) != 0
	sr.Interrupts = (value & 0x0700) >> 8
}

func (sr *StatusRegister) GetCCR() int {
	return sr.Get() & 0xff
}

func (sr *StatusRegister) SetCCR(value int) {
	sr.Set(value & 0xff)
}

// TODO: return func(result, src, dest int)
func (sr *StatusRegister) setFlags(opcode int, o *Operand, result, src, dest int) {
	resN := o.IsNegative(result)
	destN := o.IsNegative(dest)
	srcN := o.IsNegative(src)

	switch opcode {
	case flagLogical:
		sr.V = false
		sr.C = false
		sr.Z = (o.mask & uint(result)) == 0
		sr.N = resN
	case flagCmp:
		sr.Z = (o.mask & uint(result)) == 0
		sr.N = resN
		sr.C = ((uint(result) >> o.bits) & 1) != 0
		sr.V = (srcN != destN) && (resN != destN)
	case flagSub, flagAdd:
		sr.Z = (o.mask & uint(result)) == 0
		sr.N = resN
		fallthrough
	case flagSubx, flagAddx:
		sr.C = ((uint(result) >> o.bits) & 1) != 0
		sr.X = sr.C
		sr.V = (srcN != destN) && (resN != destN)
	case flagZn:
		sr.Z = sr.Z && (o.mask&uint(result)) == 0
		sr.N = resN
	}
}

func (sr *StatusRegister) conditionalTest(code uint32) func() bool {
	var condition = []func() bool{
		func() bool { return true },
		func() bool { return false },
		func() bool { return !sr.C && !sr.Z },
		func() bool { return sr.C || sr.Z },
		func() bool { return !sr.C },
		func() bool { return sr.C },
		func() bool { return !sr.Z },
		func() bool { return sr.Z },
		func() bool { return !sr.V },
		func() bool { return sr.V },
		func() bool { return !sr.N },
		func() bool { return sr.N },
		func() bool { return !(sr.N != sr.V) },
		func() bool { return sr.N != sr.V },
		func() bool { return (sr.N && sr.V && !sr.Z) || (!sr.N && !sr.V && !sr.Z) },
		func() bool { return sr.Z || (sr.N && !sr.V) || (!sr.N && sr.V) },
	}
	return condition[code&0x0f]
}
