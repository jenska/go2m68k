package cpu

type (
	Address uint32
	Data    int32
	Opcode  int16

	AddressBus interface {
		Read(address Address, operand *Operand) (int, error)
		Write(address Address, operand *Operand, value int) error
	}

	M68K struct {
		// registers
		A       [8]Address
		D       [8]Data
		SR      StatusRegister
		PC      Address
		SP, USP Address
		opcode  Opcode
		bus     AddressBus
	}
)

func NewCPU(addressBus AddressBus) M68K {
	return M68K{bus: addressBus}
}
