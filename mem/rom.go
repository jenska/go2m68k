package mem

import "github.com/jenska/atari2go/cpu"

func NewROM(start Address, rom []byte, size uint) AddressArea {
	end := start + Address(size)

	return AddressArea{
		start: start,
		end:   end,
		read: func(address Address, operand *cpu.Operand) (int, error) {
			if address >= start && address < end {
				return operand.Read(rom, uint(address-start)), nil
			}
			return 0, BusError(address)
		},
	}
}
