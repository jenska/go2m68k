package cpu

import (
	"fmt"
	glog "github.com/golang/glog"
)

type AddressHandler interface {
	Mem8(a uint32) (v uint8, ok bool)
	Mem16(a uint32) (v uint16, ok bool)
	Mem32(a uint32) (v uint32, ok bool)

	setMem8(a uint32, v uint8) bool
	setMem16(a uint32, v uint16) bool
	setMem32(a uint32, v uint32) bool
}

type RAM []uint8
type ChipSets map[uint32]AddressHandler

type MemoryHandler struct {
	ram      RAM
	chipsets ChipSets
}

func NewMemoryHandler(size int, sets ChipSets) *MemoryHandler {
	return &MemoryHandler{make([]uint8, size), sets}
}

func (mem *MemoryHandler) RegisterChipset(addresses []uint32, handler AddressHandler) {
	for _, a := range addresses {
		if mem.chipsets[a] != nil || int(a) >= len(mem.ram) {
			panic(fmt.Sprintf("address space %08x already in use", a))
		}
		mem.chipsets[a] = handler
	}
}

func (mem *MemoryHandler) Mem8(a uint32) (uint8, bool) {
	if int(a) < len(mem.ram) {
		return mem.ram[a], true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.Mem8(a)
	} else {
		return 0, false
	}
}

func (mem *MemoryHandler) setMem8(a uint32, v uint8) bool {
	if int(a) < len(mem.ram) {
		mem.ram[a] = v
		return true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.setMem8(a, v)
	}
	return false
}

func (mem *MemoryHandler) Mem16(a uint32) (uint16, bool) {
	if int(a)+1 < len(mem.ram) {
		glog.V(2).Infof("read %s at address $%08x", mem.ram[a:(a+1)], a)

		return uint16(mem.ram[a+1])<<8 | uint16(mem.ram[a]), true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.Mem16(a)
	} else {
		return 0, false
	}
}

func (mem *MemoryHandler) setMem16(a uint32, v uint16) bool {
	if int(a)+1 < len(mem.ram) {
		mem.ram[a] = uint8(v)
		mem.ram[a+1] = uint8(v >> 8)
		return true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.setMem16(a, v)
	}
	return false
}

func (mem *MemoryHandler) Mem32(a uint32) (uint32, bool) {
	if int(a)+3 < len(mem.ram) {
		r := uint32(mem.ram[a+3])<<24 | uint32(mem.ram[a+2])<<16 | uint32(mem.ram[a+1])<<8 | uint32(mem.ram[a])
		return r, true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.Mem32(a)
	} else {
		return 0, false
	}
}

func (mem *MemoryHandler) setMem32(a uint32, v uint32) bool {
	if int(a)+3 < len(mem.ram) {
		mem.ram[a] = uint8(v)
		mem.ram[a+1] = uint8(v >> 8)
		mem.ram[a+2] = uint8(v >> 16)
		mem.ram[a+3] = uint8(v >> 24)
		return true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.setMem32(a, v)
	}
	return false
}
