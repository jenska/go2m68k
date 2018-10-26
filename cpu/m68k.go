package cpu

type Address uint32

type AddressBus interface {
	Read(address Address, operand *Operand) (int, error)
	Write(address Address, operand *Operand, value int) error
}
