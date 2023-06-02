package core

import (
	"fmt"
)

const (
	// effective addressing modes.
	ModeDN   = 0  //  0: Dn
	ModeAN   = 1  //  1: An
	ModeAI   = 2  //  2: (An)
	ModePI   = 3  //  3: (An)+
	ModePD   = 4  //  4: -(An)
	ModeDI   = 5  //  5: (d,An)
	ModeIX   = 6  //  6: (d,An,Xi)
	ModeAW   = 7  //  7: (####).w
	ModeAL   = 8  //  8: (####).l
	ModeDIPC = 9  //  9: (d,PC)
	ModeIXPC = 10 // 10: (d,PC,Xi)
	ModeIM   = 11 // 11: ####
)

type (
	EA interface {
		Read() uint32
		Write(value uint32)
		Address() uint32
	}

	eaRegister struct {
		reg *uint32
		o   Operand
	}

	eaAddress struct {
		address uint32
		c       *Core
		o       Operand
	}

	eaImmediate struct {
		value uint32
		o     Operand
	}
)

func (c *Core) FetchEA(mode, xn uint16, o Operand) EA {
	switch mode {
	case ModeDN: // data register
		return eaRegister{reg: &c.D[xn], o: o}
	case ModeAN: // address register
		return eaRegister{reg: &c.A[xn], o: o}
	case ModeAI: // address register indirect
		return eaAddress{address: c.A[xn], c: c, o: o}
	case ModePI: // address register indirect with post increment
		ea := eaAddress{address: c.A[xn], c: c, o: o}
		c.A[xn] += o.Size()
		return ea
	case ModePD: // address register indirect with pre decrement
		c.A[xn] -= o.Size()
		return eaAddress{address: c.A[xn], c: c, o: o}
	case ModeDI: // address register indirect with displacement
		displacement := Word.SignedExtend(c.PopPc(Word))
		return eaAddress{address: uint32(int32(c.A[xn]) + displacement), c: c, o: o}
	case ModeIX:
		return eaAddress{address: uint32(int32(c.A[xn]) + c.extend()), c: c, o: o}
	case ModeAW:
		switch mode + xn {
		case ModeAW:
			return eaAddress{address: c.PopPc(Word), c: c, o: o}
		case ModeAL:
			return eaAddress{address: c.PopPc(Long), c: c, o: o}
		case ModeDIPC:
			displacement := Word.SignedExtend(c.PopPc(Word))
			return eaAddress{address: uint32(int32(c.PC0) + displacement), c: c, o: o}
		case ModeIXPC:
			return eaAddress{address: uint32(int32(c.PC0) + c.extend()), c: c, o: o}
		case ModeIM:
			return eaImmediate{value: c.PopPc(o), o: o}
		}
	}
	panic("invalid ea")
}

// M68000 only extend calculation
func (c Core) extend() int32 {
	e := c.PopPc(Word)
	yn := int32(c.Regs[e>>12])
	if e&0x800 == 0 {
		yn = Word.SignedExtend(uint32(yn))
	}
	return yn + Byte.SignedExtend(e)
}

func (ea eaRegister) Read() uint32 {
	return ea.o.UnsignedExtend(*ea.reg)
}

func (ea eaRegister) Write(value uint32) {
	ea.o.WriteToLong(value, ea.reg)
}

func (ea eaRegister) Address() uint32 {
	panic("unsupported operation")
}

func (ea eaRegister) String() string {
	return fmt.Sprintf("reg %s.%s", ea.o.HexString(*ea.reg), ea.o)
}

func (ea eaAddress) Read() uint32 {
	return ea.c.readRaw(ea.address, ea.o)
}

func (ea eaAddress) Write(value uint32) {
	ea.c.writeRaw(ea.address, ea.o, value)
}

func (ea eaAddress) Address() uint32 {
	return ea.address
}

func (ea eaAddress) String() string {
	return fmt.Sprintf("adr (%s).%s", ea.o.HexString(ea.address), ea.o)
}

func (ea eaImmediate) Read() uint32 {
	return ea.value
}

func (ea eaImmediate) Write(value uint32) {
	panic(fmt.Errorf("write not supported for immediate value %s", ea.o.HexString(ea.value)))
}

func (ea eaImmediate) Address() uint32 {
	panic("unsupported operation")
}

func (ea eaImmediate) String() string {
	return fmt.Sprintf("#%s.%s", ea.o.HexString(ea.value), ea.o)
}

func (c *Core) PopPc(o Operand) uint32 {
	res := c.readRaw(c.PC, o)
	c.PC += o.Align()
	return res
}

func (c *Core) Pop(o Operand) uint32 {
	res := c.readRaw(*c.SP, o)
	*c.SP += o.Size()
	return res
}

func (c *Core) Push(o Operand, v uint32) {
	*c.SP -= o.Size()
	c.writeRaw(*c.SP, o, v)
}
