package cpu

type (
	// Address ...
	Address uint32
	Data    int32

	AddressBus interface {
		Read(address Address, operand *Operand) (int, error)
		Write(address Address, operand *Operand, value int) error
		Reset()
		SetSuperVisorFlag(*bool)
	}

	M68K struct {
		// registers
		A        [8]Address
		D        [8]Data
		SR       StatusRegister
		PC       Address
		SSP, USP Address

		opcodeAddr Address
		opcode     int

		bus AddressBus

		instructions []Instruction
	}
)

func NewCPU(addressBus AddressBus) M68K {
	result := M68K{bus: addressBus}
	result.instructions = make([]Instruction, 0x10000)
	init68000InstructionSet(&result)
	result.Reset()
	return result
}

func (c *M68K) Step() int {
	c.opcodeAddr = c.PC
	if opcode, err := c.bus.Read(c.PC, Word); err == nil {
		c.opcode = opcode
		c.PC += Address(Word.size)

		if instruction := c.instructions[opcode>>6]; instruction != nil {
			return instruction()
		} else {
			return c.RaiseException(IllegalInstruction)
		}
	} else {
		panic("bus error")
	}
}

func (c *M68K) Reset() {
	c.SR.Set(0x2700)
	c.bus.Reset()

	sp, err1 := c.readA(0)
	pc, err2 := c.readA(4)
	if err1 != nil || err2 != nil {
		c.Halt()
	}
	c.SSP, c.PC = sp, pc
}

func (c *M68K) RaiseException(x Exception) int {
	oldSR := c.SR
	if !c.SR.S {
		c.SR.S = true
		c.USP = c.A[7]
		c.A[7] = c.SSP
	}
	c.pushA(c.PC)
	c.push(Word, oldSR.Get())

	if xaddr, err := c.readA(Address(x << 2)); err != nil || xaddr == 0 {
		if xaddr, err = c.readA(Address(UnintializedInterrupt << 2)); err != nil || xaddr == 0 {
			c.Halt()
		}
	} else {
		c.PC = xaddr
	}
	return 34
}

func (c *M68K) Halt() {
	panic("halt not implemented")
}

func (c *M68K) readA(address Address) (Address, error) {
	a, err := c.bus.Read(0, Long)
	return Address(a), err
}

func (c *M68K) push(o *Operand, value int) {
	c.A[7] -= Address(o.size)
	if c.bus.Write(c.A[7], o, value) != nil {
		panic("bus error")
	}
}

func (c *M68K) pushA(address Address) {
	c.push(Long, int(address))
}
