package mem

import (
	"github.com/jenska/go2m68k/pkg/cpu"
)

/// addresses

type (
	MemoryReader func(cpu.Address, *cpu.Operand) (int, error)
	MemoryWriter func(cpu.Address, *cpu.Operand, int) error
	ResetHandler func()

	AddressArea struct {
		start  cpu.Address
		end    cpu.Address
		read   MemoryReader
		write  MemoryWriter
		reset  ResetHandler
		parent *addressMap
	}

	addressMap struct {
		areas []AddressArea
		cache *AddressArea
		sv    *bool
	}
)

func (a *addressMap) findAddressArea(address cpu.Address) *AddressArea {
	if address >= a.cache.start && address < a.cache.end {
		return a.cache
	}
	for _, area := range a.areas {
		if address >= area.start && address < area.end {
			a.cache = &area
			return &area
		}
	}
	return nil
}

func (a *addressMap) Read(address cpu.Address, operand *cpu.Operand) (int, error) {
	if area := a.findAddressArea(address); area != nil {
		if read := area.read; read != nil {
			return area.read(address, operand)
		} else {
			return 0, AddressError(address)
		}
	}
	return 0, BusError(address)
}

func (a *addressMap) Write(address cpu.Address, operand *cpu.Operand, value int) error {
	if area := a.findAddressArea(address); area != nil {
		if write := area.write; write != nil {
			return area.write(address, operand, value)
		} else {
			return AddressError(address)
		}
	}
	return BusError(address)
}

func (a *addressMap) SetSuperVisorFlag(flag *bool) {
	a.sv = flag
}

func (a *addressMap) Reset() {
	for _, area := range a.areas {
		if area.reset != nil {
			area.reset()
		}
	}
}

func NewAddressBus(areas ...AddressArea) cpu.AddressBus {
	result := &addressMap{areas: areas, cache: &areas[0]}
	for _, area := range areas {
		area.parent = result
	}
	return result
}
