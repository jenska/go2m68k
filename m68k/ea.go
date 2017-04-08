package m68k

type operand struct {
	Size        uint32
	AlignedSize uint32
	Msb         uint32
	Mask        uint32
	Ext         string
	eaVecOffset int
	formatter   string
}

var Byte = &operand{1, 2, 0x80, 0xff, ".b", 0, "%02x"}
var Word = &operand{2, 2, 0x8000, 0xffff, ".w", 64, "%04x"}
var Long = &operand{4, 4, 0x80000000, 0xffffffff, ".l", 128, "%08x"}

func (o *operand) isNegative(value uint32) bool {
	return (o.Msb & value) != 0
}

func (o *operand) set(target *uint32, value uint32) {
	*target = (*target & ^o.Mask) | (value & o.Mask)
}

func (o *operand) getSigned(value uint32) int32 {
	v := uint32(value)
	if o.isNegative(v) {
		return int32(v | ^o.Mask)
	}
	return int32(v & o.Mask)
}

func (o *operand) get(value uint32) uint32 {
	return value & o.Mask
}

type ea interface {
	compute() modifier
	timing() int
}

type modifier interface {
	read() uint32
	write(value uint32)
}

// Helper for read and write of precomputed addresses
type addressModifier struct {
	cpu     *M68K
	o       *operand
	address uint32
	cycle   int
}

func (a *addressModifier) read() uint32       { return a.cpu.Read(a.o, a.address) }
func (a *addressModifier) write(value uint32) { a.cpu.Write(a.o, a.address, value) }
func (a *addressModifier) timing() int        { return a.cycle }

func newEAVectors(cpu *M68K) []ea {
	eaVec := make([]ea, 3*(1<<6))
	cyclesWord := []int{0, 0, 4, 6, 6, 8, 10, 8, 12, 8, 10, 8}
	cyclesLong := []int{0, 0, 8, 10, 10, 12, 14, 12, 16, 12, 14, 12}
	for _, o := range []*operand{Byte, Word, Long} {
		cycles := cyclesWord
		if o == Long {
			cycles = cyclesLong
		}
		for i := 0; i < 8; i++ {
			eaVec[i+o.eaVecOffset] = &eaDataRegister{cpu, o, i}
			eaVec[i+8+o.eaVecOffset] = &eaAddressRegister{cpu, o, i}
			eaVec[i+16+o.eaVecOffset] = &eaAddressRegisterIndirect{&addressModifier{cpu, o, 0, cycles[2]}, i}
			eaVec[i+24+o.eaVecOffset] = &eaAddressRegisterPostInc{&addressModifier{cpu, o, 0, cycles[3]}, i}
			eaVec[i+32+o.eaVecOffset] = &eaAddressRegisterPreDec{&addressModifier{cpu, o, 0, cycles[4]}, i}
			eaVec[i+40+o.eaVecOffset] = &eaAddressRegisterWithDisplacement{&addressModifier{cpu, o, 0, cycles[5]}, i}
			eaVec[i+48+o.eaVecOffset] = &eaAddressRegisterWithIndex{&addressModifier{cpu, o, 0, cycles[6]}, i}
		}
		eaVec[56+o.eaVecOffset] = &eaAbsoluteWord{&addressModifier{cpu, o, 0, cycles[7]}}
		eaVec[57+o.eaVecOffset] = &eaAbsoluteLong{&addressModifier{cpu, o, 0, cycles[8]}}
		eaVec[58+o.eaVecOffset] = &eaPCWithDisplacement{&addressModifier{cpu, o, 0, cycles[9]}}
		eaVec[59+o.eaVecOffset] = &eaPCWithIndex{&addressModifier{cpu, o, 0, cycles[10]}}
		eaVec[60+o.eaVecOffset] = &eaImmediate{&addressModifier{cpu, o, 0, cycles[11]}}
	}
	return eaVec
}

// 0 Dx
type eaDataRegister struct {
	cpu      *M68K
	o        *operand
	register int
}

func (ea *eaDataRegister) compute() modifier  { return ea }
func (ea *eaDataRegister) timing() int        { return 0 }
func (ea *eaDataRegister) read() uint32       { return ea.o.get(ea.cpu.D[ea.register]) }
func (ea *eaDataRegister) write(value uint32) { ea.o.set(&(ea.cpu.D[ea.register]), value) }

// 1 Ax
type eaAddressRegister eaDataRegister

func (ea *eaAddressRegister) compute() modifier  { return ea }
func (ea *eaAddressRegister) timing() int        { return 0 }
func (ea *eaAddressRegister) read() uint32       { return ea.o.get(ea.cpu.A[ea.register]) }
func (ea *eaAddressRegister) write(value uint32) { ea.o.set(&(ea.cpu.A[ea.register]), value) }

// 2 (Ax)
type eaAddressRegisterIndirect struct {
	*addressModifier
	register int
}

func (ea *eaAddressRegisterIndirect) compute() modifier {
	ea.address = ea.cpu.A[ea.register]
	return ea
}

// 3 (Ax)+
type eaAddressRegisterPostInc eaAddressRegisterIndirect

func (ea *eaAddressRegisterPostInc) compute() modifier {
	ea.address = ea.cpu.A[ea.register]
	ea.cpu.A[ea.register] += ea.o.Size
	return ea
}

// 4 -(Ax)
type eaAddressRegisterPreDec eaAddressRegisterIndirect

func (ea *eaAddressRegisterPreDec) init(cpu *M68K, o *operand, register int) {
	cpu.A[register] -= o.Size
	ea.cpu, ea.o, ea.address = cpu, o, cpu.A[register]
}
func (ea *eaAddressRegisterPreDec) compute() modifier {
	ea.cpu.A[ea.register] -= ea.o.Size
	ea.address = ea.cpu.A[ea.register]
	return ea
}

// 5 xxxx(Ax)
type eaAddressRegisterWithDisplacement eaAddressRegisterIndirect

func (ea *eaAddressRegisterWithDisplacement) compute() modifier {
	ea.address = uint32(int32(ea.cpu.A[ea.register]) + int32(int16(ea.cpu.popPC(Word))))
	return ea
}

// 5 xxxx(PC)
type eaPCWithDisplacement struct {
	*addressModifier
}

func (ea *eaPCWithDisplacement) compute() modifier {
	ea.address = uint32(int32(ea.cpu.PC) + int32(int16(ea.cpu.popPC(Word))))
	return ea
}

// 6 xx(Ax, Rx.w/.l)
type eaAddressRegisterWithIndex eaAddressRegisterIndirect

func (ea *eaAddressRegisterWithIndex) compute() modifier {
	ext := int(int16(ea.cpu.popPC(Word)))
	displacement := int(int8(ext))
	idxRegNumber := (ext >> 12) & 0x07
	idxSize := (ext & 0x0800) == 0x0800
	idxValue := 0
	if (ext & 0x8000) == 0x8000 { // address register
		if idxSize {
			idxValue = int(int16(ea.cpu.A[idxRegNumber]))
		} else {
			idxValue = int(ea.cpu.A[idxRegNumber])
		}
	} else { // data register
		if idxSize {
			idxValue = int(int16(ea.cpu.D[idxRegNumber]))
		} else {
			idxValue = int(ea.cpu.D[idxRegNumber])
		}
	}
	ea.address = uint32(int(ea.cpu.A[ea.register]) + idxValue + displacement)
	return ea
}

// 6 xx(PC, Rx.w/.l)
type eaPCWithIndex eaPCWithDisplacement

func (ea *eaPCWithIndex) compute() modifier {
	ext := int(int16(ea.cpu.popPC(Word)))
	displacement := int(int8(ext))
	idxRegNumber := (ext >> 12) & 0x07
	idxSize := (ext & 0x0800) == 0x0800
	idxValue := 0
	if (ext & 0x8000) == 0x8000 { // address register
		if idxSize {
			idxValue = int(int16(ea.cpu.A[idxRegNumber]))
		} else {
			idxValue = int(ea.cpu.A[idxRegNumber])
		}
	} else { // data register
		if idxSize {
			idxValue = int(int16(ea.cpu.D[idxRegNumber]))
		} else {
			idxValue = int(ea.cpu.D[idxRegNumber])
		}
	}
	ea.address = uint32(int(ea.cpu.PC) + idxValue + displacement)
	return ea
}

// 7. xxxx.w
type eaAbsoluteWord struct {
	*addressModifier
}

func (ea *eaAbsoluteWord) compute() modifier {
	ea.address = ea.cpu.popPC(Word)
	return ea
}

// 8. xxxx.l
type eaAbsoluteLong eaAbsoluteWord

func (ea *eaAbsoluteLong) compute() modifier {
	ea.address = ea.cpu.popPC(Long)
	return ea
}

// 9. #value
type eaImmediate struct {
	*addressModifier
}

func (ea *eaImmediate) compute() modifier {
	ea.address = uint32(ea.cpu.popPC(ea.o))
	return ea
}

func (ea *eaImmediate) read() uint32 {
	return ea.address
}

func (ea *eaImmediate) write(value uint32) {}
