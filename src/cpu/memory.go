package cpu

type Address int32
type AddressHandler interface {
	Mem8(a Address) (v uint8, ok bool)
	Mem16(a Address) (v uint16, ok bool)
	Mem32(a Address) (v uint32, ok bool)

	setMem8(a Address, v uint8) bool
	setMem16(a Address, v uint16) bool
	setMem32(a Address, v uint32) bool
}

type RAM []uint8
type ChipSets map[Address]AddressHandler

type MemoryHandler struct {
	ram      RAM
	chipsets ChipSets
}

func NewMemoryHandler(size int, sets ChipSets) *MemoryHandler {
	return &MemoryHandler{make([]uint8, size), sets}
}

func (mem *MemoryHandler) RegisterChipset(addresses []Address, handler AddressHandler) {
	for _, a := range addresses {
		if mem.chipsets[a] != nil || int(a) >= len(mem.ram) {
			panic("address space already in use")
		}
		mem.chipsets[a] = handler
	}
}

func (mem *MemoryHandler) Mem8(a Address) (uint8, bool) {
	if int(a) < len(mem.ram) {
		return mem.ram[a], true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.Mem8(a)
	} else {
		return 0, false
	}
}

func (mem *MemoryHandler) setMem8(a Address, v uint8) bool {
	if int(a) < len(mem.ram) {
		mem.ram[a] = v
		return true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.setMem8(a, v)
	}
	return false
}

func (mem *MemoryHandler) Mem16(a Address) (uint16, bool) {
	if int(a)+1 < len(mem.ram) {
		return uint16((mem.ram[a] << 8) | mem.ram[a+1]), true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.Mem16(a)
	} else {
		return 0, false
	}
}

func (mem *MemoryHandler) setMem16(a Address, v uint16) bool {
	if int(a)+1 < len(mem.ram) {
		mem.ram[a] = uint8(v >> 8)
		mem.ram[a+1] = uint8(v)
		return true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.setMem16(a, v)
	}
	return false
}

func (mem *MemoryHandler) Mem32(a Address) (uint32, bool) {
	if int(a)+3 < len(mem.ram) {
		return uint32((mem.ram[a] << 24) | (mem.ram[a+1] << 16) | (mem.ram[a+2] << 8) | mem.ram[a+3]), true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.Mem32(a)
	} else {
		return 0, false
	}
}

func (mem *MemoryHandler) setMem32(a Address, v uint32) bool {
	if int(a)+3 < len(mem.ram) {
		mem.ram[a] = uint8(v >> 24)
		mem.ram[a+1] = uint8(v >> 16)
		mem.ram[a+2] = uint8(v >> 8)
		mem.ram[a+3] = uint8(v)
		return true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.setMem32(a, v)
	}
	return false
}
