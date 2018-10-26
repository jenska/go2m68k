package mem

import "github.com/jenska/atari2go/cpu"

func NewRAM(start cpu.Address, size uint) AddressArea {
	ram := make([]byte, size)
	end := start + cpu.Address(size)

	return AddressArea{
		start: start,
		end:   end,
		write: func(address cpu.Address, operand *cpu.Operand, value int) error {
			if address >= start && address < end {
				operand.Write(ram[:], uint(address-start), value)
				return nil
			}
			return BusError(address)
		},
		read: func(address cpu.Address, operand *cpu.Operand) (int, error) {
			if address >= start && address < end {
				return operand.Read(ram[:], uint(address-start)), nil
			}
			return 0, BusError(address)
		},
	}
}

func NewProtectedRAM(start cpu.Address, size uint, sr *cpu.StatusRegister) AddressArea {
	area := NewRAM(start, size)
	protectedWrite := area.write
	area.write = func(address cpu.Address, operand *cpu.Operand, value int) error {
		if sr.S {
			return protectedWrite(address, operand, value)
		}
		return cpu.SuperVisorException(address)
	}
	return area
}
