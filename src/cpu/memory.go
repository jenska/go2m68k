package cpu

import (
	"fmt"
)

type AddressHandler interface {
	Mem(o *Operand, a uint32) (v uint32, ok bool)
	setMem(o *Operand, a, v uint32) bool
}

type RAM []uint8
type ChipSets map[uint32]AddressHandler

type MemoryHandler struct {
	ram      RAM
	chipsets ChipSets
}

func NewMemoryHandler(size int) *MemoryHandler {
	return &MemoryHandler{make([]uint8, size), ChipSets{}}
}

func (mem *MemoryHandler) RegisterChipset(addresses []uint32, handler AddressHandler) {
	for _, a := range addresses {
		if mem.chipsets[a] != nil || int(a) < len(mem.ram) {
			panic(fmt.Sprintf("address space %08x already in use", a))
		}
		mem.chipsets[a] = handler
	}
}

func (mem *MemoryHandler) Mem(o *Operand, a uint32) (uint32, bool) {
	if int(a+o.Size) <= len(mem.ram) {
		r := uint32(mem.ram[a])
		switch o {
		case Long:
			r |= uint32(mem.ram[a+3])<<24 | uint32(mem.ram[a+2])<<16
			fallthrough
		case Word:
			r |= uint32(mem.ram[a+1]) << 8
		}
		return r, true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.Mem(o,a)
	} else {
		return 0, false
	}
}

func (mem *MemoryHandler) setMem(o *Operand, a, v uint32) bool {
	if int(a+o.Size) <= len(mem.ram) {
		mem.ram[a] = uint8(v)
		switch o {
		case Long:
			mem.ram[a+3] = uint8(v >> 24)
			mem.ram[a+2] = uint8(v >> 16)
		case Word:
			mem.ram[a+1] = uint8(v >> 8)
		}
		return true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.setMem(o, a, v)
	}
	return false
}
