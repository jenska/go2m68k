package m68k

type operand struct {
	Size        uint32
	AlignedSize uint32
	Msb         uint32
	Mask        uint32
	Ext         string
	eaOffset    int
	formatter   string
}

var (
	Byte = &operand{1, 2, 0x80, 0xff, ".b", 0, "%02x"}
	Word = &operand{2, 2, 0x8000, 0xffff, ".w", 64, "%04x"}
	Long = &operand{4, 4, 0x80000000, 0xffffffff, ".l", 128, "%08x"}
)

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
}

type modifier interface {
	read() uint32
	write(value uint32)
}

func (cpu *M68K) initEATable() []ea {
	eaTable := make([]ea, 3*(1<<6))
	for _, o := range []*operand{Byte, Word, Long} {
		for i := 0; i < 8; i++ {
			eaTable[i+o.eaOffset] = &eaDataRegister{cpu, o, i}
			eaTable[i+8+o.eaOffset] = &eaAddressRegister{cpu, o, i}
			eaTable[i+16+o.eaOffset] = &eaAddressRegisterIndirect{&addressModifier{cpu, o, 0}, i}
			eaTable[i+24+o.eaOffset] = &eaAddressRegisterPostInc{&addressModifier{cpu, o, 0}, i}
			eaTable[i+32+o.eaOffset] = &eaAddressRegisterPreDec{&addressModifier{cpu, o, 0}, i}
			eaTable[i+40+o.eaOffset] = &eaAddressRegisterWithDisplacement{&addressModifier{cpu, o, 0}, i}
			eaTable[i+48+o.eaOffset] = &eaAddressRegisterWithIndex{&addressModifier{cpu, o, 0}, i}
		}
		eaTable[56+o.eaOffset] = &eaAbsoluteWord{&addressModifier{cpu, o, 0}}
		eaTable[57+o.eaOffset] = &eaAbsoluteLong{&addressModifier{cpu, o, 0}}
		eaTable[58+o.eaOffset] = &eaPCWithDisplacement{&addressModifier{cpu, o, 0}}
		eaTable[59+o.eaOffset] = &eaPCWithIndex{&addressModifier{cpu, o, 0}}
		eaTable[60+o.eaOffset] = &eaImmediate{&addressModifier{cpu, o, 0}}
	}
	return eaTable
}

func (cpu *M68K) loadEA(o *operand, eaType int) ea {
	return cpu.eaTable[eaType+o.eaOffset]
}

// Helper for read and write of precomputed addresses
type addressModifier struct {
	cpu     *M68K
	o       *operand
	address uint32
}

func (a *addressModifier) read() uint32       { return a.cpu.Read(a.o, a.address) }
func (a *addressModifier) write(value uint32) { a.cpu.Write(a.o, a.address, value) }

// 0 Dx
type eaDataRegister struct {
	cpu      *M68K
	o        *operand
	register int
}

func (ea *eaDataRegister) compute() modifier  { return ea }
func (ea *eaDataRegister) read() uint32       { return ea.o.get(ea.cpu.D[ea.register]) }
func (ea *eaDataRegister) write(value uint32) { ea.o.set(&(ea.cpu.D[ea.register]), value) }

// 1 Ax
type eaAddressRegister eaDataRegister

func (ea *eaAddressRegister) compute() modifier  { return ea }
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

func (ea *eaAddressRegisterPreDec) compute() modifier {
	if ea.register == 7 {
		ea.cpu.A[ea.register] -= ea.o.AlignedSize
	} else {
		ea.cpu.A[ea.register] -= ea.o.Size
	}
	ea.address = ea.cpu.A[ea.register]
	return ea
}

// 5 xxxx(Ax)
type eaAddressRegisterWithDisplacement eaAddressRegisterIndirect

func (ea *eaAddressRegisterWithDisplacement) compute() modifier {
	ea.address = uint32(int32(ea.cpu.A[ea.register]) + int32(ea.cpu.IRC))
	ea.cpu.readExtensionWord()
	return ea
}

// 5 xxxx(PC)
type eaPCWithDisplacement struct{ *addressModifier }

func (ea *eaPCWithDisplacement) compute() modifier {
	ea.address = uint32(int32(ea.cpu.PC) + int32(ea.cpu.IRC))
	ea.cpu.readExtensionWord()
	return ea
}

// 6 xx(Ax, Rx.w/.l)
type eaAddressRegisterWithIndex eaAddressRegisterIndirect

func (ea *eaAddressRegisterWithIndex) compute() modifier {
	ext := int(ea.cpu.IRC)
	displacement := ext & 0xff
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
	ea.cpu.sync(2)
	ea.cpu.readExtensionWord()
	return ea
}

// 6 xx(PC, Rx.w/.l)
type eaPCWithIndex eaPCWithDisplacement

func (ea *eaPCWithIndex) compute() modifier {
	ext := int(ea.cpu.IRC)
	displacement := ext & 0xff
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
	ea.cpu.sync(2)
	ea.cpu.readExtensionWord()
	return ea
}

// 7. xxxx.w
type eaAbsoluteWord struct{ *addressModifier }

func (ea *eaAbsoluteWord) compute() modifier {
	ea.address = uint32(ea.cpu.IRC)
	ea.cpu.readExtensionWord()
	return ea
}

// 8. xxxx.l
type eaAbsoluteLong eaAbsoluteWord

func (ea *eaAbsoluteLong) compute() modifier {
	ea.address = uint32(ea.cpu.IRC << 16)
	ea.cpu.readExtensionWord()
	ea.address |= uint32(ea.cpu.IRC)
	ea.cpu.readExtensionWord()
	return ea
}

// 9. #value
type eaImmediate struct{ *addressModifier }

func (ea *eaImmediate) compute() modifier {
	ea.address = uint32(ea.cpu.IRC) & ea.o.Mask
	if ea.o == Long {
		ea.address <<= 16
		ea.cpu.readExtensionWord()
		ea.address |= uint32(ea.cpu.IRC)
	}
	return ea
}

func (ea *eaImmediate) read() uint32       { return ea.address }
func (ea *eaImmediate) write(value uint32) {}
