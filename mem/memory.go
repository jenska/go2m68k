package mem

import (
	"container/list"
	"fmt"
)

type IllegalAddressError struct {
	address uint32
}

func (e IllegalAddressError) Error() string {
	return fmt.Sprintf("Failed to access address %08x", e.address)
}

type AddressHandler interface {
	Read(size, address uint32) (value uint32, err error)
	Write(size, address, value uint32) error
	Start() uint32
	End() uint32
}

type MemoryHandler struct {
	ram      []uint8
	chipsets *list.List
}

func NewMemoryHandler(size int) *MemoryHandler {
	return &MemoryHandler{make([]uint8, size), list.New()}
}

func (mem *MemoryHandler) RegisterChipset(addressHandler AddressHandler) {
	s1, e1 := addressHandler.Start(), addressHandler.End()
	for e := mem.chipsets.Front(); e != nil; e = e.Next() {
		handler := e.Value.(AddressHandler)
		s2, e2 := handler.Start(), handler.End()
		if (s1 >= s2 && s1 <= e2) || (e1 <= e2 && e1 >= s2) {
			panic(fmt.Errorf("address range $%08x - $%08x already allocated by %s ($%08x - $%08x)", s1, e1, handler, s2, e2))
		}
	}
	mem.chipsets.PushFront(addressHandler)
}

func (mem *MemoryHandler) Start() uint32 { return 0 }
func (mem *MemoryHandler) End() uint32   { return uint32(len(mem.ram) - 1) }

func (mem *MemoryHandler) lookup(address uint32) AddressHandler {
	for e := mem.chipsets.Front(); e != nil; e = e.Next() {
		handler := e.Value.(AddressHandler)
		if handler.Start() >= address || handler.End() <= address {
			return handler
		}
	}
	return nil
}

func (mem *MemoryHandler) Read(size, a uint32) (uint32, error) {
	if int(a+size) <= len(mem.ram) {
		r := uint32(mem.ram[a])
		switch size {
		case 4:
			r |= uint32(mem.ram[a+3])<<24 | uint32(mem.ram[a+2])<<16
			fallthrough
		case 2:
			r |= uint32(mem.ram[a+1]) << 8
		}
		return r, nil
	} else if handler := mem.lookup(a); handler != nil {
		return handler.Read(size, a)
	} else {
		return 0, IllegalAddressError{a}
	}
}

func (mem *MemoryHandler) Write(size, a, v uint32) error {
	if int(a+size) <= len(mem.ram) {
		mem.ram[a] = uint8(v)
		switch size {
		case 4:
			mem.ram[a+3] = uint8(v >> 24)
			mem.ram[a+2] = uint8(v >> 16)
			fallthrough
		case 2:
			mem.ram[a+1] = uint8(v >> 8)
		}
		return nil
	} else if handler := mem.lookup(a); handler != nil {
		return handler.Write(size, a, v)
	}
	return IllegalAddressError{a}
}
