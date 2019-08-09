package mem

import "github.com/jenska/go2m68k/pkg/cpu"

func NewROM(start cpu.Address, rom []byte) AddressArea {
	end := start + cpu.Address(len(rom))

	return AddressArea{
		start: start,
		end:   end,
		read: func(address cpu.Address, operand *cpu.Operand) (int, error) {
			if address >= start && address < end {
				return operand.Read(rom[address-start:]), nil
			}
			return 0, BusError(address)
		},
	}
}
