package cpu

const (
	eaModeIndirect      = 0b010000
	eaModePostIncrement = 0b011000
	eaModePreDecrement  = 0b100000
	eaModeAbsoluteShort = 0b111000
	eaModeAbsoluteLong  = 0b111001
	eaModeImmidiate     = 0b111100
	eaModeDisplacement  = 0b101000
)

type (
	ea interface {
		init(cpu *M68K, o *Size) modifier
		computedAddress() int32
		// cycles() int
	}

	modifier interface {
		computedAddress() int32
		read() int32
		write(int32)
	}

	eaRegister struct {
		reg  func(cpu *M68K) *int32
		cpu  *M68K
		size *Size
	}

	eaRegisterIndirect struct {
		eaRegister
		address int32
	}

	eaPostIncrement eaRegisterIndirect

	eaPreDecrement eaRegisterIndirect

	eaDisplacement eaRegisterIndirect

	eaIndirectIndex struct {
		eaRegisterIndirect
		index func(cpu *M68K, a int32) int32
	}

	eaAbsolute struct {
		cpu     *M68K
		eaSize  *Size
		size    *Size
		address int32
	}

	eaPCDisplacement eaDisplacement

	eaPCIndirectIndex eaIndirectIndex

	eaImmediate struct {
		value int32
	}

	eaStatusRegister struct {
		size *Size
		sr   *ssr
	}
)

var (
	eaSrc68000 = []ea{
		&eaRegister{reg: dy},
		&eaRegister{reg: ay},
		&eaRegisterIndirect{eaRegister{reg: ay}, 0},
		&eaPostIncrement{eaRegister{reg: ay}, 0},
		&eaPreDecrement{eaRegister{reg: ay}, 0},
		&eaDisplacement{eaRegister{reg: ay}, 0},
		&eaIndirectIndex{eaRegisterIndirect{eaRegister{reg: ay}, 0}, ix68000},
		&eaAbsolute{eaSize: Word},
		&eaAbsolute{eaSize: Long},
		&eaPCDisplacement{eaRegister{reg: nil}, 0},
		&eaPCIndirectIndex{eaRegisterIndirect{eaRegister{reg: nil}, 0}, ix68000},
		&eaImmediate{},
	}

	eaDst68000 = []ea{
		&eaRegister{reg: dy},
		&eaRegister{reg: ay},
		&eaRegisterIndirect{eaRegister{reg: ay}, 0},
		&eaPostIncrement{eaRegister{reg: ay}, 0},
		&eaPreDecrement{eaRegister{reg: ay}, 0},
		&eaDisplacement{eaRegister{reg: ay}, 0},
		&eaIndirectIndex{eaRegisterIndirect{eaRegister{reg: ay}, 0}, ix68000},
		&eaAbsolute{eaSize: Word},
		&eaAbsolute{eaSize: Long},
		&eaPCDisplacement{eaRegister{reg: nil}, 0},
		&eaPCIndirectIndex{eaRegisterIndirect{eaRegister{reg: nil}, 0}, ix68000},
		&eaStatusRegister{},
	}

	eaSrc68020 = []ea{
		&eaRegister{reg: dy},
		&eaRegister{reg: ay},
		&eaRegisterIndirect{eaRegister{reg: ay}, 0},
		&eaPostIncrement{eaRegister{reg: ay}, 0},
		&eaPreDecrement{eaRegister{reg: ay}, 0},
		&eaDisplacement{eaRegister{reg: ay}, 0},
		&eaIndirectIndex{eaRegisterIndirect{eaRegister{reg: ay}, 0}, ix68020},
		&eaAbsolute{eaSize: Word},
		&eaAbsolute{eaSize: Long},
		&eaPCDisplacement{eaRegister{reg: nil}, 0},
		&eaPCIndirectIndex{eaRegisterIndirect{eaRegister{reg: nil}, 0}, ix68020},
		&eaImmediate{},
	}

	eaDst68020 = []ea{
		&eaRegister{reg: dy},
		&eaRegister{reg: ay},
		&eaRegisterIndirect{eaRegister{reg: ay}, 0},
		&eaPostIncrement{eaRegister{reg: ay}, 0},
		&eaPreDecrement{eaRegister{reg: ay}, 0},
		&eaDisplacement{eaRegister{reg: ay}, 0},
		&eaIndirectIndex{eaRegisterIndirect{eaRegister{reg: ay}, 0}, ix68020},
		&eaAbsolute{eaSize: Word},
		&eaAbsolute{eaSize: Long},
		&eaPCDisplacement{eaRegister{reg: nil}, 0},
		&eaPCIndirectIndex{eaRegisterIndirect{eaRegister{reg: nil}, 0}, ix68020},
		&eaStatusRegister{},
	}
)

// TODO add cycles
func (cpu *M68K) resolveSrcEA(o *Size) modifier {
	mode := (cpu.ir >> 3) & 0x07
	if mode < 7 {
		return cpu.eaSrc[mode].init(cpu, o)
	}
	return cpu.eaSrc[mode+y(cpu.ir)].init(cpu, o)
}

// TODO add cycles
func (cpu *M68K) resolveDstEA(o *Size) modifier {
	mode := (cpu.ir >> 3) & 0x07
	if mode < 7 {
		return cpu.eaDst[mode].init(cpu, o)
	}
	return cpu.eaDst[mode+y(cpu.ir)].init(cpu, o)
}

func (cpu *M68K) push(s *Size, value int32) {
	cpu.a[7] -= s.size
	cpu.write(cpu.a[7], s, value)
}

func (cpu *M68K) pop(s *Size) int32 {
	res := cpu.read(cpu.a[7], s)
	cpu.a[7] += s.size // sometimes odd
	return res
}

func (cpu *M68K) popPc(s *Size) int32 {
	res := cpu.read(cpu.pc, s)
	cpu.pc += s.align // never odd
	return res
}

func x(ir uint16) uint16 { return (ir >> 9) & 0x7 }
func y(ir uint16) uint16 { return ir & 0x7 }

func dx(cpu *M68K) *int32 { return &cpu.d[x(cpu.ir)] }
func dy(cpu *M68K) *int32 { return &cpu.d[y(cpu.ir)] }

func ax(cpu *M68K) *int32 { return &cpu.a[x(cpu.ir)] }
func ay(cpu *M68K) *int32 { return &cpu.a[y(cpu.ir)] }

// -------------------------------------------------------------------
// Data register

func (ea *eaRegister) init(cpu *M68K, o *Size) modifier {
	ea.cpu, ea.size = cpu, o
	return ea
}

func (ea *eaRegister) read() int32 {
	return *ea.reg(ea.cpu) & int32(ea.size.mask)
}

func (ea *eaRegister) write(v int32) {
	ea.size.set(v, ea.reg(ea.cpu))
}

func (ea *eaRegister) computedAddress() int32 {
	panic("no address in register addressing mode")
}

// -------------------------------------------------------------------
// Address register indirect

func (ea *eaRegisterIndirect) init(cpu *M68K, o *Size) modifier {
	ea.cpu, ea.size, ea.address = cpu, o, *ea.reg(cpu)
	return ea
}

func (ea *eaRegisterIndirect) read() int32 {
	return ea.cpu.read(ea.address, ea.size)
}

func (ea *eaRegisterIndirect) write(v int32) {
	ea.cpu.write(ea.address, ea.size, v)
}

func (ea *eaRegisterIndirect) computedAddress() int32 {
	return ea.address
}

// -------------------------------------------------------------------
// Post increment

func (ea *eaPostIncrement) init(cpu *M68K, o *Size) modifier {
	ea.cpu, ea.size, ea.address = cpu, o, *ea.reg(cpu)
	*ea.reg(cpu) += ea.size.size
	return ea
}

func (ea *eaPostIncrement) read() int32 {
	return ea.cpu.read(ea.address, ea.size)
}

func (ea *eaPostIncrement) write(v int32) {
	ea.cpu.write(ea.address, ea.size, v)
}

// -------------------------------------------------------------------
// Pre decrement

func (ea *eaPreDecrement) init(cpu *M68K, o *Size) modifier {
	*ea.reg(cpu) -= o.size
	ea.cpu, ea.size, ea.address = cpu, o, *ea.reg(cpu)
	return ea
}

func (ea *eaPreDecrement) read() int32 {
	return ea.cpu.read(ea.address, ea.size)
}

func (ea *eaPreDecrement) write(v int32) {
	ea.cpu.write(ea.address, ea.size, v)
}

// -------------------------------------------------------------------
// Displacement

func (ea *eaDisplacement) init(cpu *M68K, o *Size) modifier {
	ea.cpu, ea.size = cpu, o
	ea.address = *ea.reg(cpu) + int32(int16(cpu.popPc(Word)))
	return ea
}

func (ea *eaDisplacement) computedAddress() int32 {
	return ea.address
}

func (ea *eaPCDisplacement) init(cpu *M68K, o *Size) modifier {
	ea.cpu, ea.size = cpu, o
	ea.address = cpu.pc + cpu.popPc(Word)
	return ea
}

func (ea *eaPCDisplacement) computedAddress() int32 {
	return ea.address
}

// -------------------------------------------------------------------
// Indirect + index

func (ea *eaIndirectIndex) init(cpu *M68K, o *Size) modifier {
	ea.cpu, ea.size = cpu, o
	ea.address = ea.index(cpu, *ea.reg(cpu))
	return ea
}

func (ea *eaPCIndirectIndex) init(cpu *M68K, o *Size) modifier {
	ea.cpu, ea.size = cpu, o
	ea.address = ea.index(cpu, cpu.pc)
	return ea
}

// -------------------------------------------------------------------
// absolute word and long

func (ea *eaAbsolute) init(cpu *M68K, o *Size) modifier {
	ea.cpu, ea.size = cpu, o
	ea.address = cpu.popPc(ea.eaSize)
	return ea
}

func (ea *eaAbsolute) read() int32 {
	return ea.cpu.read(ea.address, ea.size)
}

func (ea *eaAbsolute) write(v int32) {
	ea.cpu.write(ea.address, ea.size, v)
}

func (ea *eaAbsolute) computedAddress() int32 {
	return ea.address
}

// -------------------------------------------------------------------
// immediate

func (ea *eaImmediate) init(cpu *M68K, o *Size) modifier {
	ea.value = cpu.popPc(o)
	return ea
}

func (ea *eaImmediate) read() int32 {
	return ea.value
}

func (ea *eaImmediate) write(v int32) {
	panic("write on immediate addressing mode")
}

func (ea *eaImmediate) computedAddress() int32 {
	panic("no adress in immediate addressing mode")
}

// -------------------------------------------------------------------
// sr

func (ea *eaStatusRegister) init(cpu *M68K, o *Size) modifier {
	ea.sr = &cpu.sr
	return ea
}

func (ea *eaStatusRegister) read() int32 {
	if ea.size == Byte {
		return ea.sr.ccr()
	}
	return ea.sr.bits()
}

func (ea *eaStatusRegister) write(v int32) {
	if ea.size == Byte {
		ea.sr.setccr(v)
	} else {
		ea.sr.setbits(v)
	}
}

func (ea *eaStatusRegister) computedAddress() int32 {
	panic("no adress in status register addressing mode")
}

// -------------------------------------------------------------------
// Indexed addressing modes are encoded as follows:
//
// Base instruction format:
// F E D C B A 9 8 7 6 | 5 4 3 | 2 1 0
// x x x x x x x x x x | 1 1 0 | BASE REGISTER      (An)
//
// Base instruction format for destination EA in move instructions:
// F E D C | B A 9    | 8 7 6 | 5 4 3 2 1 0
// x x x x | BASE REG | 1 1 0 | X X X X X X       (An)
//
// Brief extension format:
//  F  |  E D C   |  B  |  A 9  | 8 | 7 6 5 4 3 2 1 0
// D/A | REGISTER | W/L | SCALE | 0 |  DISPLACEMENT
//
// Full extension format:
//  F     E D C      B     A 9    8   7    6    5 4       3   2 1 0
// D/A | REGISTER | W/L | SCALE | 1 | BS | IS | BD SIZE | 0 | I/IS
// BASE DISPLACEMENT (0, 16, 32 bit)                (bd)
// OUTER DISPLACEMENT (0, 16, 32 bit)               (od)
//
// D/A:     0 = Dn, 1 = An                          (Xn)
// W/L:     0 = W (sign extend), 1 = L              (.SIZE)
// SCALE:   00=1, 01=2, 10=4, 11=8                  (*SCALE)
// BS:      0=add base reg, 1=suppress base reg     (An suppressed)
// IS:      0=add index, 1=suppress index           (Xn suppressed)
// BD SIZE: 00=reserved, 01=NULL, 10=Word, 11=Long  (size of bd)
//
// IS I/IS Operation
// 0  000  No Memory Indirect
// 0  001  indir prex with null outer
// 0  010  indir prex with word outer
// 0  011  indir prex with long outer
// 0  100  reserved
// 0  101  indir postx with null outer
// 0  110  indir postx with word outer
// 0  111  indir postx with long outer
// 1  000  no memory indirect
// 1  001  mem indir with null outer
// 1  010  mem indir with word outer
// 1  011  mem indir with long outer
// 1  100-111  reserved
//
func ix68000(c *M68K, a int32) int32 {
	ext := c.popPc(Word)

	xn := c.da[ext>>12]
	if (ext & 0x800) == 0 {
		xn = int32(int16(ext))
	}
	return a + xn + int32(int8(ext))
}

func ix68020(c *M68K, a int32) int32 {
	ext := c.popPc(Word)
	var xn int32

	// brief extension format
	if (ext & 0x80) == 0 {

		xn = c.da[ext>>12]
		if (ext & 0x800) == 0 {
			xn = int32(int16(ext))
		}
		// Add scale
		xn <<= (ext >> 9) & 3
		return a + xn + int32(int8(ext))
	}

	// full extension format
	c.iclocks -= c.eaIdxCycles[ext&0x3f]
	if (ext & 0x40) != 0 { //BS
		a = 0
	}
	if (ext & 0x20) == 0 { // IS
		xn = c.da[ext>>12]
		if (ext & 0x800) == 0 {
			xn <<= (ext >> 9) & 3
		}
	}
	var bd int32
	if (ext & 0x10) != 0 {
		if (ext & 0x80) != 0 {
			bd = c.popPc(Long)
		} else {
			bd = c.popPc(Word)
		}
	}
	if (ext & 7) == 0 {
		return a + bd + xn
	}
	var od int32
	if (ext & 0x02) != 0 {
		if (ext & 0x01) != 0 {
			od = c.popPc(Long)
		} else {
			od = c.popPc(Word)
		}
	}
	if (ext & 0x04) != 0 {
		return c.read(a+bd, Long) + xn + od
	}
	return c.read(a+bd+xn, Long) + od
}
