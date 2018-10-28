package cpu

type (
	// Address ...
	Address uint32
	Data    int32
	Opcode  int16

	AddressBus interface {
		Read(address Address, operand *Operand) (int, error)
		Write(address Address, operand *Operand, value int) error
		Reset()
	}

	M68K struct {
		// registers
		A        [8]Address
		D        [8]Data
		SR       StatusRegister
		PC       Address
		SSP, USP Address
		opcode   Opcode
		bus      AddressBus
	}
)

func NewCPU(addressBus AddressBus) M68K {
	result := M68K{bus: addressBus}
	result.Reset()
	return result
}

func (c *M68K) readA(address Address) (Address, error) {
	a, err := c.bus.Read(0, Long)
	return Address(a), err
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

func (c *M68K) Halt() {
	panic("halt not implemented")
}
