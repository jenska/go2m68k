package cpu

import (
	"container/list"
	"fmt"
)

const (
	XPT_SPR  = 0
	XPT_PCR  = 1
	XPT_BUS  = 2
	XPT_ADR  = 3
	XPT_ILL  = 4
	XPT_DBZ  = 5
	XPT_CHKN = 6
	XPT_TRV  = 7
	XPT_PRV  = 8
	XPT_TRC  = 9
	XPT_LNA  = 10
	XPT_LNF  = 11

	XPT_FPU = 13

	UNINITIALIZED_INTERRUPT_VECTOR = 15

	SPURIOUS_INTERRUPT           = 24
	LEVEL_1_INTERRUPT_AUTOVECTOR = 25
	LEVEL_2_INTERRUPT_AUTOVECTOR = 26
	LEVEL_3_INTERRUPT_AUTOVECTOR = 27
	LEVEL_4_INTERRUPT_AUTOVECTOR = 28
	LEVEL_5_INTERRUPT_AUTOVECTOR = 29
	LEVEL_6_INTERRUPT_AUTOVECTOR = 30
	LEVEL_7_INTERRUPT_AUTOVECTOR = 31

	TRAP_INSTRUCTION_VECTORS = 32
	USER_INTERRUPT_VECTORS   = 64
)

type IllegalAddressError struct {
	address uint32
}

func (e IllegalAddressError) Error() string {
	return fmt.Sprintf("Failed to access address %08x", e.address)
}

type AddressHandler interface {
	Read(o *Operand, a uint32) (v uint32, err error)
	Write(o *Operand, a, v uint32) error
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

func (mem *MemoryHandler) Read(o *Operand, a uint32) (uint32, error) {
	if int(a+o.Size) <= len(mem.ram) {
		r := uint32(mem.ram[a])
		switch o {
		case Long:
			r |= uint32(mem.ram[a+3])<<24 | uint32(mem.ram[a+2])<<16
			fallthrough
		case Word:
			r |= uint32(mem.ram[a+1]) << 8
		}
		return r, nil
	} else if handler := mem.lookup(a); handler != nil {
		return handler.Read(o, a)
	} else {
		return 0, IllegalAddressError{a}
	}
}

func (mem *MemoryHandler) Write(o *Operand, a, v uint32) error {
	if int(a+o.Size) <= len(mem.ram) {
		mem.ram[a] = uint8(v)
		switch o {
		case Long:
			mem.ram[a+3] = uint8(v >> 24)
			mem.ram[a+2] = uint8(v >> 16)
			fallthrough
		case Word:
			mem.ram[a+1] = uint8(v >> 8)
		}
		return nil
	} else if handler := mem.lookup(a); handler != nil {
		return handler.Write(o, a, v)
	}
	return IllegalAddressError{a}
}
