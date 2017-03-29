package cpu

type EA interface {
	get(o *Operand) uint32
	set(o *Operand, value uint32)
	timing(o *Operand) int
	computedAddress() uint32
}

// 0 Dx
type EADataRegister struct {
	cpu *M68k
	register int
}
func (ea *EADataRegister) get(o *Operand) uint32 { return o.get(ea.cpu.D[ea.register])  }
func (ea *EADataRegister) set(o *Operand, value uint32) { o.set(&(ea.cpu.D[ea.register]), value)}
func (*EADataRegister) timing(o *Operand) int { return  0 }
func (*EADataRegister) computedAddress() uint32 { return 0 }

// 1 Ax
type EAAddressRegister struct {
	cpu *M68k
	register int
}
func (ea *EAAddressRegister) get(o *Operand) uint32 { return uint32(o.get(uint32(ea.cpu.A[ea.register])))  }
func (ea *EAAddressRegister) set(o *Operand, value uint32) { o.set(&(ea.cpu.A[ea.register]), value)}
func (*EAAddressRegister) timing(operand *Operand) int { return  0 }
func (*EAAddressRegister) computedAddress() uint32 { return 0 }

// 2 (Ax)
type EAAddressRegisterIndirect struct{
	cpu *M68k
	register int
}

func (ea *EAAddressRegisterIndirect) get(o *Operand) uint32 { return ea.cpu.mem(o, ea.cpu.A[ea.register])  }
func (ea *EAAddressRegisterIndirect) set(o *Operand, value uint32) { o.set(&(ea.cpu.A[ea.register]), value)}
func (*EAAddressRegisterIndirect) timing(o *Operand) int { if o == Long { return 8 } else { return 4 }}
func (ea *EAAddressRegisterIndirect) computedAddress() uint32 { return ea.cpu.A[ea.register] }



func (cpu *M68k) mem(o *Operand, address uint32) uint32 {
	address &= 0x00ffffff
	switch o {
	case Byte:
		if v, ok :=  cpu.memory.Mem8(address); ok {
			return uint32(v)
		}
	case Word:
		if v, ok := cpu.memory.Mem16(address); ok && (address&1)==0 {
			return uint32(v)
		}
	case Long:
		if v, ok := cpu.memory.Mem32(address); ok && (address&1)==0 {
			return v
		}
	default:
		// TODO raise exception
	}
	return 0
}

/*
func (cpu *M68k) push16(v uint16) {
	sp := cpu.A[7] - 2
	cpu.setMem16(sp, v)
	cpu.A[7] = sp
}

func (cpu *M68k) push32(v uint32) {
	sp := cpu.A[7] - 4
	cpu.setMem32(sp, v)
	cpu.A[7] = sp
}

func (cpu *M68k) setMem8(a Address, v uint8) {
	a &= 0xffffff
	if !cpu.memory.setMem8(a, v) {
		// TODO bus error
	}
}

func (cpu *M68k) setMem16(a Address, v uint16) {
	a &= 0xffffff
	if !cpu.memory.setMem16(a, v) {
		// TODO bus error
	}
}

func (cpu *M68k) setMem32(a Address, v uint32) {
	a &= 0xffffff
	if !cpu.memory.setMem32(a, v) {
		// TODO bus error
	}
}
*/
