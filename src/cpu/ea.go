package cpu

type Operand struct {
	Size        uint32
	AlignedSize uint32
	Msb         uint32
	Mask        uint32
	Ext         string
	eaVecOffset int
	formatter   string
}

var Byte = &Operand{1, 2, 0x80, 0xff, ".b", 0, "%02x"}
var Word = &Operand{2, 2, 0x8000, 0xffff, ".w", 64, "%04x"}
var Long = &Operand{4, 4, 0x80000000, 0xffffffff, ".l", 128, "%08x"}

func (o *Operand) isNegative(value uint32) bool {
	return (o.Msb & value) != 0
}

func (o *Operand) set(target *uint32, value uint32) {
	*target = (*target & ^o.Mask) | (value & o.Mask)
}

func (o *Operand) getSigned(value uint32) int32 {
	v := uint32(value)
	if o.isNegative(v) {
		return int32(v | ^o.Mask)
	}
	return int32(v & o.Mask)
}

func (o *Operand) get(value uint32) uint32 {
	return value & o.Mask
}

type EA interface {
	compute() Modifier
	timing() int
}

type Modifier interface {
	read() uint32
	write(value uint32)
}

// Helper for read and write of precomputed addresses
type addressModifier struct {
	cpu     *M68k
	o       *Operand
	address uint32
	cycle   int
}

func (a *addressModifier) read() uint32       { return a.cpu.read(a.o, a.address) }
func (a *addressModifier) write(value uint32) { a.cpu.write(a.o, a.address, value) }
func (a *addressModifier) timing() int        { return a.cycle }

func NewEAVectors(cpu *M68k) []EA {
	eaVec := make([]EA, 3*(1<<6))
	cyclesWord := []int{0, 0, 4, 6, 6, 8, 10, 8, 12, 8, 10}
	cyclesLong := []int{0, 0, 8, 10, 10, 12, 14, 12, 16, 12, 14}
	for _, o := range []*Operand{Byte, Word, Long} {
		cycles := cyclesWord
		if o == Long {
			cycles = cyclesLong
		}
		for i := 0; i < 8; i++ {
			eaVec[i+o.eaVecOffset] = &EADataRegister{cpu, o, i}
			eaVec[i+8+o.eaVecOffset] = &EAAddressRegister{cpu, o, i}
			eaVec[i+16+o.eaVecOffset] = &EAAddressRegisterIndirect{&addressModifier{cpu, o, 0, cycles[2]}, i}
			eaVec[i+24+o.eaVecOffset] = &EAAddressRegisterPostInc{&addressModifier{cpu, o, 0, cycles[3]}, i}
			eaVec[i+32+o.eaVecOffset] = &EAAddressRegisterPreDec{&addressModifier{cpu, o, 0, cycles[4]}, i}
			eaVec[i+40+o.eaVecOffset] = &EAAddressRegisterWithDisplacement{&addressModifier{cpu, o, 0, cycles[5]}, i}
			eaVec[i+48+o.eaVecOffset] = &EAAddressRegisterWithIndex{&addressModifier{cpu, o, 0, cycles[6]}, i}
		}
		eaVec[56+o.eaVecOffset] = &EAAbsoluteWord{&addressModifier{cpu, o, 0, cycles[7]}}
		eaVec[57+o.eaVecOffset] = &EAAbsoluteLong{&addressModifier{cpu, o, 0, cycles[8]}}
		eaVec[58+o.eaVecOffset] = &EAPCWithDisplacement{&addressModifier{cpu, o, 0, cycles[9]}}
		eaVec[59+o.eaVecOffset] = &EAPCWithIndex{&addressModifier{cpu, o, 0, cycles[10]}}
	}
	return eaVec
}

// 0 Dx
type EADataRegister struct {
	cpu      *M68k
	o        *Operand
	register int
}

func (ea *EADataRegister) compute() Modifier  { return ea }
func (ea *EADataRegister) timing() int        { return 0 }
func (ea *EADataRegister) read() uint32       { return ea.cpu.D[ea.register] }
func (ea *EADataRegister) write(value uint32) { ea.o.set(&(ea.cpu.D[ea.register]), value) }

// 1 Ax
type EAAddressRegister EADataRegister

func (ea *EAAddressRegister) read() uint32       { return uint32(ea.o.get(uint32(ea.cpu.A[ea.register]))) }
func (ea *EAAddressRegister) write(value uint32) { ea.o.set(&(ea.cpu.A[ea.register]), value) }
func (ea *EAAddressRegister) timing() int        { return 0 }
func (ea *EAAddressRegister) compute() Modifier  { return ea }

// 2 (Ax)
type EAAddressRegisterIndirect struct {
	*addressModifier
	register int
}

func (ea *EAAddressRegisterIndirect) compute() Modifier {
	ea.address = ea.cpu.A[ea.register]
	return ea
}

// 3 (Ax)+
type EAAddressRegisterPostInc EAAddressRegisterIndirect

func (ea *EAAddressRegisterPostInc) compute() Modifier {
	ea.address = ea.cpu.A[ea.register]
	ea.cpu.A[ea.register] += ea.o.Size
	return ea
}

// 4 -(Ax)
type EAAddressRegisterPreDec EAAddressRegisterIndirect

func (ea *EAAddressRegisterPreDec) init(cpu *M68k, o *Operand, register int) {
	cpu.A[register] -= o.Size
	ea.cpu, ea.o, ea.address = cpu, o, cpu.A[register]
}
func (ea *EAAddressRegisterPreDec) compute() Modifier {
	ea.cpu.A[ea.register] -= ea.o.Size
	ea.address = ea.cpu.A[ea.register]
	return ea
}

// 5 xxxx(Ax)
type EAAddressRegisterWithDisplacement EAAddressRegisterIndirect

func (ea *EAAddressRegisterWithDisplacement) compute() Modifier {
	ea.address = uint32(int32(ea.cpu.A[ea.register]) + int32(int16(ea.cpu.popPC(Word))))
	return ea
}

// 5 xxxx(PC)
type EAPCWithDisplacement struct {
	*addressModifier
}

func (ea *EAPCWithDisplacement) compute() Modifier {
	ea.address = uint32(int32(ea.cpu.PC) + int32(int16(ea.cpu.popPC(Word))))
	return ea
}

// 6 xx(Ax, Rx.w/.l)
type EAAddressRegisterWithIndex EAAddressRegisterIndirect

func (ea *EAAddressRegisterWithIndex) compute() Modifier {
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
type EAPCWithIndex EAPCWithDisplacement

func (ea *EAPCWithIndex) compute() Modifier {
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
type EAAbsoluteWord struct {
	*addressModifier
}

func (ea *EAAbsoluteWord) compute() Modifier {
	ea.address = ea.cpu.popPC(Word)
	return ea
}

// 8. xxxx.l
type EAAbsoluteLong EAAbsoluteWord

func (ea *EAAbsoluteLong) compute() Modifier {
	ea.address = ea.cpu.popPC(Long)
	return ea
}

func (cpu *M68k) read(o *Operand, address uint32) uint32 {
	address &= 0x00ffffff
	switch o {
	case Byte:
		if v, ok := cpu.memory.Mem8(address); ok {
			return uint32(v)
		}
	case Word:
		if v, ok := cpu.memory.Mem16(address); ok && (address&1) == 0 {
			return uint32(v)
		}
	case Long:
		if v, ok := cpu.memory.Mem32(address); ok && (address&1) == 0 {
			return v
		}
	}
	// TODO raise exception
	return 0
}

func (cpu *M68k) write(o *Operand, address uint32, value uint32) {
	address &= 0x00ffffff
	switch o {
	case Byte:
		if cpu.memory.setMem8(address, uint8(value)) {
			return
		}
	case Word:
		if (address&1) == 0 && cpu.memory.setMem16(address, uint16(value)) {
			return
		}
	case Long:
		if (address&1) == 0 && cpu.memory.setMem32(address, value) {
			return
		}
	}
	// TODO raise exception
}

func (cpu *M68k) popPC(o *Operand) uint32 {
	result := cpu.read(o, cpu.PC)
	cpu.PC += o.Size
	return result
}
