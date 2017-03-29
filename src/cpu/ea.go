package cpu

type Operand struct {
	Size        uint32
	AlignedSize uint32
	Msb         uint32
	Mask        uint32
	Ext         string
	formatter   string
}

var Byte = &Operand{1, 2, 0x80, 0xff, ".b", "%02x"}
var Word = &Operand{2, 2, 0x8000, 0xffff, ".w", "%04x"}
var Long = &Operand{4, 4, 0x80000000, 0xffffffff, ".l", "%08x"}

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
	init(cpu *M68k, o *Operand, param int)
	get() uint32
	set(value uint32)
	timing() int
	computedAddress() uint32
}

// 0 Dx
type EADataRegister struct {
	cpu      *M68k
	o        *Operand
	register int
}

func (ea *EADataRegister) init(cpu *M68k, o *Operand, register int) {
	ea.cpu, ea.o, ea.register = cpu, o, register
}
func (ea *EADataRegister) get() uint32             { return ea.o.get(ea.cpu.D[ea.register]) }
func (ea *EADataRegister) set(value uint32)        { ea.o.set(&(ea.cpu.D[ea.register]), value) }
func (ea *EADataRegister) timing() int             { return 0 }
func (ea *EADataRegister) computedAddress() uint32 { return 0 }

// 1 Ax
type EAAddressRegister EADataRegister

func (ea *EAAddressRegister) init(cpu *M68k, o *Operand, register int) {
	ea.cpu, ea.o, ea.register = cpu, o, register
}
func (ea *EAAddressRegister) get() uint32          { return uint32(ea.o.get(uint32(ea.cpu.A[ea.register]))) }
func (ea *EAAddressRegister) set(value uint32)     { ea.o.set(&(ea.cpu.A[ea.register]), value) }
func (*EAAddressRegister) timing() int             { return 0 }
func (*EAAddressRegister) computedAddress() uint32 { return 0 }

// 2 (Ax)
type EAAddressRegisterIndirect struct {
	cpu     *M68k
	o       *Operand
	address uint32
}

func (ea *EAAddressRegisterIndirect) init(cpu *M68k, o *Operand, register int) {
	ea.cpu, ea.o, ea.address = cpu, o, cpu.A[register]
}
func (ea *EAAddressRegisterIndirect) get() uint32      { return ea.cpu.read(ea.o, ea.address) }
func (ea *EAAddressRegisterIndirect) set(value uint32) { ea.cpu.write(ea.o, ea.address, value) }
func (ea *EAAddressRegisterIndirect) timing() int {
	if ea.o == Long {
		return 8
	} else {
		return 4
	}
}
func (ea *EAAddressRegisterIndirect) computedAddress() uint32 { return ea.address }

// 3 (Ax)+
type EAAddressRegisterPostInc EAAddressRegisterIndirect

func (ea *EAAddressRegisterPostInc) init(cpu *M68k, o *Operand, register int) {
	ea.cpu, ea.o, ea.address = cpu, o, cpu.A[register]
	cpu.A[register] += o.Size
}
func (ea *EAAddressRegisterPostInc) get() uint32      { return ea.cpu.read(ea.o, ea.address) }
func (ea *EAAddressRegisterPostInc) set(value uint32) { ea.cpu.write(ea.o, ea.address, value) }
func (ea *EAAddressRegisterPostInc) timing() int {
	if ea.o == Long {
		return 10
	} else {
		return 6
	}
}
func (ea *EAAddressRegisterPostInc) computedAddress() uint32 { return ea.address }

// 4 -(Ax)
type EAAddressRegisterPreDec EAAddressRegisterIndirect

func (ea *EAAddressRegisterPreDec) init(cpu *M68k, o *Operand, register int) {
	cpu.A[register] -= o.Size
	ea.cpu, ea.o, ea.address = cpu, o, cpu.A[register]
}
func (ea *EAAddressRegisterPreDec) get() uint32      { return ea.cpu.read(ea.o, ea.address) }
func (ea *EAAddressRegisterPreDec) set(value uint32) { ea.cpu.write(ea.o, ea.address, value) }
func (ea *EAAddressRegisterPreDec) timing() int {
	if ea.o == Long {
		return 10
	} else {
		return 6
	}
}
func (ea *EAAddressRegisterPreDec) computedAddress() uint32 { return ea.address }

// 5 xxxx(Ax)
type EAAddressRegisterWithDisplacement EAAddressRegisterIndirect

func (ea *EAAddressRegisterWithDisplacement) init(cpu *M68k, o *Operand, register int) {
	ea.cpu, ea.o, ea.address = cpu, o, uint32(int32(cpu.A[register])+int32(int16(cpu.popPC(Word))))
}
func (ea *EAAddressRegisterWithDisplacement) get() uint32      { return ea.cpu.read(ea.o, ea.address) }
func (ea *EAAddressRegisterWithDisplacement) set(value uint32) { ea.cpu.write(ea.o, ea.address, value) }
func (ea *EAAddressRegisterWithDisplacement) timing() int {
	if ea.o == Long {
		return 12
	} else {
		return 8
	}
}
func (ea *EAAddressRegisterWithDisplacement) computedAddress() uint32 { return ea.address }

// 5 xxxx(PC)
type EAPCWithDisplacement EAAddressRegisterIndirect

func (ea *EAPCWithDisplacement) init(cpu *M68k, o *Operand, _ int) {
	ea.cpu, ea.o, ea.address = cpu, o, uint32(int32(cpu.PC)+int32(int16(cpu.popPC(Word))))
}
func (ea *EAPCWithDisplacement) get() uint32      { return ea.cpu.read(ea.o, ea.address) }
func (ea *EAPCWithDisplacement) set(value uint32) { ea.cpu.write(ea.o, ea.address, value) }
func (ea *EAPCWithDisplacement) timing() int {
	if ea.o == Long {
		return 12
	} else {
		return 8
	}
}
func (ea *EAPCWithDisplacement) computedAddress() uint32 { return ea.address }

// 6 xx(Ax, Rx.w/.l)
type EAAddressRegisterWithIndex EAAddressRegisterIndirect

func (ea *EAAddressRegisterWithIndex) init(cpu *M68k, o *Operand, register int) {
	ea.cpu, ea.o = cpu, o
	ext := int(int16(cpu.popPC(Word)))
	displacement := int(int8(ext))
	idxRegNumber := (ext >> 12) & 0x07
	idxSize := (ext & 0x0800) == 0x0800
	idxValue := 0
	if (ext & 0x8000) == 0x8000 { // address register
		if idxSize {
			idxValue = int(int16(cpu.A[idxRegNumber]))
		} else {
			idxValue = int(cpu.A[idxRegNumber])
		}
	} else { // data register
		if idxSize {
			idxValue = int(int16(cpu.D[idxRegNumber]))
		} else {
			idxValue = int(cpu.D[idxRegNumber])
		}
	}
	ea.address = uint32(int(cpu.A[register]) + idxValue + displacement)
}
func (ea *EAAddressRegisterWithIndex) get() uint32      { return ea.cpu.read(ea.o, ea.address) }
func (ea *EAAddressRegisterWithIndex) set(value uint32) { ea.cpu.write(ea.o, ea.address, value) }
func (ea *EAAddressRegisterWithIndex) timing() int {
	if ea.o == Long {
		return 14
	} else {
		return 10
	}
}
func (ea *EAAddressRegisterWithIndex) computedAddress() uint32 { return ea.address }

// 6 xx(PC, Rx.w/.l)
type EAPCWithIndex EAAddressRegisterIndirect

func (ea *EAPCWithIndex) init(cpu *M68k, o *Operand, _ int) {
	ea.cpu, ea.o = cpu, o
	ext := int(int16(cpu.popPC(Word)))
	displacement := int(int8(ext))
	idxRegNumber := (ext >> 12) & 0x07
	idxSize := (ext & 0x0800) == 0x0800
	idxValue := 0
	if (ext & 0x8000) == 0x8000 { // address register
		if idxSize {
			idxValue = int(int16(cpu.A[idxRegNumber]))
		} else {
			idxValue = int(cpu.A[idxRegNumber])
		}
	} else { // data register
		if idxSize {
			idxValue = int(int16(cpu.D[idxRegNumber]))
		} else {
			idxValue = int(cpu.D[idxRegNumber])
		}
	}
	ea.address = uint32(int(cpu.PC) + idxValue + displacement)
}
func (ea *EAPCWithIndex) get() uint32      { return ea.cpu.read(ea.o, ea.address) }
func (ea *EAPCWithIndex) set(value uint32) { ea.cpu.write(ea.o, ea.address, value) }
func (ea *EAPCWithIndex) timing() int {
	if ea.o == Long {
		return 14
	} else {
		return 10
	}
}
func (ea *EAPCWithIndex) computedAddress() uint32 { return ea.address }

// 7. xxxx.w
type EAAbsoluteWord EAAddressRegisterIndirect

func (ea *EAAbsoluteWord) init(cpu *M68k, o *Operand, r int) {
	ea.cpu, ea.o, ea.address = cpu, o, cpu.popPC(Word)
}
func (ea *EAAbsoluteWord) get() uint32      { return ea.cpu.read(ea.o, ea.address) }
func (ea *EAAbsoluteWord) set(value uint32) { ea.cpu.write(ea.o, ea.address, value) }
func (ea *EAAbsoluteWord) timing() int {
	if ea.o == Long {
		return 12
	} else {
		return 8
	}
}
func (ea *EAAbsoluteWord) computedAddress() uint32 { return ea.address }

// 8. xxxx.l
type EAAbsoluteLong EAAddressRegisterIndirect

func (ea *EAAbsoluteLong) init(cpu *M68k, o *Operand, r int) {
	ea.cpu, ea.o, ea.address = cpu, o, cpu.popPC(Long)
}
func (ea *EAAbsoluteLong) get() uint32      { return ea.cpu.read(ea.o, ea.address) }
func (ea *EAAbsoluteLong) set(value uint32) { ea.cpu.write(ea.o, ea.address, value) }
func (ea *EAAbsoluteLong) timing() int {
	if ea.o == Long {
		return 16
	} else {
		return 12
	}
}
func (ea *EAAbsoluteLong) computedAddress() uint32 { return ea.address }

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
