package cpu

import "strconv"

const (
	flagLogical = iota
	flagCmp
	flagAdd
	flagSub
	flagAddx
	flagSubx
	flagZn
)

// SSR for M68000 cpu
type ssr struct {
	C, V, Z, N, X, S, T1, T0, M bool
	Interrupts                  int32
}

func appendFlag(flag bool, name string, str string) string {
	if flag {
		return str + name
	}
	return str + "-"
}

func (sr ssr) String() string {
	result := ""
	result = appendFlag(sr.C, "C", result)
	result = appendFlag(sr.V, "V", result)
	result = appendFlag(sr.Z, "Z", result)
	result = appendFlag(sr.N, "N", result)
	result = appendFlag(sr.X, "X", result)
	result = appendFlag(sr.S, "S", result)
	result = appendFlag(sr.T1, "T", result)
	result = appendFlag(sr.T0, "0", result)
	result = appendFlag(sr.N, "M", result)
	result += strconv.Itoa(int(sr.Interrupts))
	return result
}

// Get the status register as a bitmap
func (sr *ssr) bits() int32 {
	var result int32 = 0
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
func (sr *ssr) setbits(value int32) {
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

func (sr *ssr) ccr() int32 {
	return sr.bits() & 0xff
}

func (sr *ssr) setccr(value int32) {
	sr.C = (value & 1) != 0
	sr.V = (value & 2) != 0
	sr.Z = (value & 4) != 0
	sr.N = (value & 8) != 0
	sr.X = (value & 16) != 0
}

func (sr *ssr) x1() int32 {
	if sr.X {
		return 1
	}
	return 0
}

func (sr *ssr) setLogicalFlags(o *Size, data int32) {
	sr.N = o.IsNegative(data)
	sr.Z = o.mask&uint32(data) == 0
	sr.V = false
	sr.C = false
}

func (sr *ssr) setAddSubFlags(o *Size, src, dest, result int32) {
	sr.Z = (o.mask & uint32(result)) == 0
	sr.N = o.IsNegative(result)
	sr.setAddSubXFlags(o, src, dest, result)
}

func (sr *ssr) setAddSubXFlags(o *Size, src, dest, result int32) {
	sr.C = ((uint32(result) >> o.bits) & 1) != 0
	sr.X = sr.C
	srcN := o.IsNegative(src)
	destN := o.IsNegative(dest)
	sr.V = (srcN != destN) && (sr.N != destN)
}

// TODO: return func(result, src, dest int)
/*

func (sr *ssr) setFlags(opcode int, s *Size, result, src, dest int) {
	resN := s.IsNegative(result)
	destN := s.IsNegative(dest)
	srcN := s.IsNegative(src)

	switch opcode {
	case flagLogical:
		sr.V = false
		sr.C = false
		sr.Z = (s.mask & uint(result)) == 0
		sr.N = resN
	case flagCmp:
		sr.Z = (s.mask & uint(result)) == 0
		sr.N = resN
		sr.C = ((uint(result) >> s.bits) & 1) != 0
		sr.V = (srcN != destN) && (resN != destN)
	case flagSub, flagAdd:
		sr.Z = (s.mask & uint(result)) == 0
		sr.N = resN
		fallthrough
	case flagSubx, flagAddx:
		sr.C = ((uint(result) >> s.bits) & 1) != 0
		sr.X = sr.C
		sr.V = (srcN != destN) && (resN != destN)
	case flagZn:
		sr.Z = sr.Z && (s.mask&uint32(result)) == 0
		sr.N = resN
	}
}
*/

func (sr *ssr) testCC(code uint16) bool {
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
	return condition[code&0x0f]()
}
