package m68k2go

import (
	"fmt"
)

type (
	ea interface {
		compute() modifier
		timing() int
	}

	modifier interface {
		read() int
		write(value int)
	}

	// 0 Dx
	eaDataRegister struct {
		o        *operand
		register *int
	}
	// 1 Ax
	eaAddressRegister eaDataRegister
	// 2 (Ax)
	eaAddressRegisterIndirect struct {
		addressModifier
		register *int
	}
	// 3 (Ax)+
	eaAddressRegisterPostInc eaAddressRegisterIndirect
	// 4 -(Ax)
	eaAddressRegisterPreDec eaAddressRegisterIndirect
	// 5 xxxx(Ax)
	eaAddressRegisterWithDisplacement eaAddressRegisterIndirect
	// 5 xxxx(PC)
	eaPCWithDisplacement struct {
		addressModifier
	}
	// 6 xx(Ax, Rx.w/.l)
	eaAddressRegisterWithIndex eaAddressRegisterIndirect
	// 6 xx(PC, Rx.w/.l)
	eaPCWithIndex eaPCWithDisplacement
	// 7. xxxx.w
	eaAbsoluteWord struct {
		addressModifier
	}
	// 8. xxxx.l
	eaAbsoluteLong eaAbsoluteWord
	// 9. #value
	eaImmediate struct {
		addressModifier
	}

	// Helper for read and write of precomputed addresses
	addressModifier struct {
		cpu     *m68k
		o       *operand
		address int
		cycles  int
	}
)

func (c *m68k) readImm(o *operand) int {
	result := c.Read(c.pc, o)
	c.pc += o.align
	return result
}

func (c *m68k) operandY() *operand {
	ops := []*operand{Byte, Word, Long}
	return ops[(c.ir>>6)&0x3]
}

func (c *m68k) eaY(o *operand) ea {
	switch (c.ir >> 3) & 0x7 {
	case 0:
		return &eaDataRegister{o, c.dy()}
	case 1:
		return &eaAddressRegister{o, c.ay()}
	case 2:
		return &eaAddressRegisterIndirect{addressModifier{c, o, 0, o.cycles(4, 8)}, c.ay()}
	case 3:
		return &eaAddressRegisterPostInc{addressModifier{c, o, 0, o.cycles(6, 10)}, c.ay()}
	case 4:
		return &eaAddressRegisterPreDec{addressModifier{c, o, 0, o.cycles(6, 10)}, c.ay()}
	case 5:
		return &eaAddressRegisterWithDisplacement{addressModifier{c, o, 0, o.cycles(8, 12)}, c.ay()}
	case 6:
		return &eaAddressRegisterWithIndex{addressModifier{c, o, 0, o.cycles(10, 14)}, c.ay()}
	case 7:
		switch c.ir & 0x7 {
		case 0:
			return &eaAbsoluteWord{addressModifier{c, o, 0, o.cycles(8, 12)}}
		case 1:
			return &eaAbsoluteLong{addressModifier{c, o, 0, o.cycles(12, 16)}}
		case 2:
			return &eaPCWithDisplacement{addressModifier{c, o, 0, o.cycles(8, 12)}}
		case 3:
			return &eaPCWithIndex{addressModifier{c, o, 0, o.cycles(10, 14)}}
		case 4:
			return &eaImmediate{addressModifier{c, o, 0, o.cycles(8, 12)}}
		}
	}
	panic(fmt.Sprintf("illegal adressing mode %d for %s", c.ir&0xf, c.Name()))
}

func (c *m68k) modOp() (modifier, *operand) {
	o := c.operandY()
	return c.eaY(o).compute(), o
}

func (c *m68k) immOp() (int, *operand) {
	o := c.operandY()
	return c.readImm(o), o
}

func (c *m68k) dx() *int { return &c.d[(c.ir>>9)&0x7] }
func (c *m68k) dy() *int { return &c.d[c.ir&0x7] }

func (c *m68k) ax() *int { return &c.a[(c.ir>>9)&0x7] }
func (c *m68k) ay() *int { return &c.a[c.ir&0x7] }

func (a *addressModifier) read() int       { return a.cpu.Read(a.address, a.o) }
func (a *addressModifier) write(value int) { a.cpu.Write(a.address, a.o, value) }
func (a *addressModifier) timing() int     { return a.cycles }

func (ea *eaDataRegister) compute() modifier { return ea }
func (ea *eaDataRegister) timing() int       { return 0 }
func (ea *eaDataRegister) read() int         { return ea.o.get(*ea.register) }
func (ea *eaDataRegister) write(value int)   { ea.o.set(value, ea.register) }

func (ea *eaAddressRegister) compute() modifier { return ea }
func (ea *eaAddressRegister) timing() int       { return 0 }
func (ea *eaAddressRegister) read() int         { return ea.o.get(*ea.register) }
func (ea *eaAddressRegister) write(value int)   { ea.o.set(value, ea.register) }

func (ea *eaAddressRegisterIndirect) compute() modifier {
	ea.address = *ea.register
	return ea
}

func (ea *eaAddressRegisterPostInc) compute() modifier {
	ea.address = *ea.register
	*ea.register += ea.o.size
	return ea
}

func (ea *eaAddressRegisterPreDec) compute() modifier {
	*ea.register -= ea.o.size
	ea.address = *ea.register
	return ea
}

func (ea *eaAddressRegisterWithDisplacement) compute() modifier {
	ea.address = *ea.register + ea.cpu.readImm(Word)
	return ea
}

func (ea *eaPCWithDisplacement) compute() modifier {
	ea.address = ea.cpu.pc + ea.cpu.readImm(Word)
	return ea
}

func (ea *eaAddressRegisterWithIndex) compute() modifier {
	ext := ea.cpu.readImm(Word)
	displacement := int(int8(ext))
	idxRegNumber := (ext >> 12) & 0x07
	idxValue := 0
	if (ext & 0x8000) == 0x8000 { // address register
		idxValue = ea.cpu.a[idxRegNumber]
	} else { // data register
		idxValue = ea.cpu.d[idxRegNumber]
	}
	ea.address = *ea.register + idxValue + displacement
	return ea
}

func (ea *eaPCWithIndex) compute() modifier {
	ext := ea.cpu.readImm(Word)
	displacement := int(int8(ext))
	idxRegNumber := (ext >> 12) & 0x07
	idxValue := 0
	if (ext & 0x8000) == 0x8000 { // address register
		idxValue = int(ea.cpu.a[idxRegNumber])
	} else { // data register
		idxValue = int(ea.cpu.d[idxRegNumber])
	}
	if (ext & 0x0800) == 0 {
		idxValue = int(int16(idxValue))
	}
	ea.address = ea.cpu.pc + idxValue + displacement
	return ea
}

func (ea *eaAbsoluteWord) compute() modifier {
	ea.address = ea.cpu.readImm(Word)
	return ea
}

func (ea *eaAbsoluteLong) compute() modifier {
	ea.address = ea.cpu.readImm(Long)
	return ea
}

func (ea *eaImmediate) compute() modifier {
	ea.address = ea.cpu.readImm(ea.o)
	return ea
}

func (ea *eaImmediate) read() int       { return ea.address }
func (ea *eaImmediate) write(value int) {}
