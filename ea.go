package cpu

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
		o        *Size
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
		cpu     *M68K
		o       *Size
		address uint32
		// cycles  int
	}
)

func (c *M68K) readImm(o *Size) int {
	result := c.read(c.PC, o)
	c.PC += o.align
	return result
}

func (c *M68K) operandY() *Size {
	ops := []*Size{Byte, Word, Long}
	return ops[(c.IR>>6)&0x3]
}

/*
func (c *M68K) eaY(o *Size) ea {
	switch (c.IR >> 3) & 0x7 {
	case 0:
		return eaDataRegister{o, c.dy()}
	case 1:
		return &eaAddressRegister{o, c.ay()}
	case 2:
		return &eaAddressRegisterIndirect{addressModifier{c, o, 0}, c.ay()}
	case 3:
		return &eaAddressRegisterPostInc{addressModifier{c, o, 0}, c.ay()}
	case 4:
		return &eaAddressRegisterPreDec{addressModifier{c, o, 0}, c.ay()}
	case 5:
		return &eaAddressRegisterWithDisplacement{addressModifier{c, o, 0}, c.ay()}
	case 6:
		return &eaAddressRegisterWithIndex{addressModifier{c, o, 0}, c.ay()}
	case 7:
		switch c.IR & 0x7 {
		case 0:
			return &eaAbsoluteWord{addressModifier{c, o, 0}}
		case 1:
			return &eaAbsoluteLong{addressModifier{c, o, 0}}
		case 2:
			return &eaPCWithDisplacement{addressModifier{c, o, 0}}
		case 3:
			return &eaPCWithIndex{addressModifier{c, o, 0}}
		case 4:
			return &eaImmediate{addressModifier{c, o, 0}}
		}
	}
	panic(fmt.Sprintf("illegal adressing mode %d", c.IR&0xf))
}
func (c *M68K) modOp() (modifier, *Size) {
	o := c.operandY()
	return c.eaY(o).compute(), o
}

func (c *M68K) immOp() (int, *Size) {
	o := c.operandY()
	return c.readImm(o), o
}

func (c *M68K) dx() int32 { return c.D[(c.IR>>9)&0x7] }
func (c *M68K) dy() int32 { return c.D[c.IR&0x7] }

func (c *M68K) ax() int32 { return c.D[(c.IR>>9)&0x7] }
func (c *M68K) ay() int32 { return c.D[c.IR&0x7] }

func (a *addressModifier) read() int       { return a.cpu.read(a.address, a.o) }
func (a *addressModifier) write(value int) { a.cpu.write(a.address, a.o, value) }
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
*/
