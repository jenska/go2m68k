package cpu

import (
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
		return handler.Mem(o, a)
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
			fallthrough
		case Word:
			mem.ram[a+1] = uint8(v >> 8)
		}
		return true
	} else if handler := mem.chipsets[a]; handler != nil {
		return handler.setMem(o, a, v)
	}
	return false
}
