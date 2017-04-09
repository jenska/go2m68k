package m68k

import (
	"container/list"
	"fmt"

	"github.com/golang/glog"
)

const (
	XptSpr  = 0
	XptPcr  = 1
	XptBus  = 2
	XptAdr  = 3
	XptIll  = 4
	XptDbz  = 5
	XptChkn = 6
	XptTrv  = 7
	XptPrv  = 8
	XptTrc  = 9
	XptLna  = 10
	XptLnf  = 11
	XptFPU  = 13

	XptUnitializedInterrupt = 15

	XptSpuriousInterrupt = 24
	XptLevel1Interrupt   = 25
	XptLevel2Interrupt   = 26
	XptLevel3Interrupt   = 27
	XptLevel4Interrupt   = 28
	XptLevel5Interrupt   = 29
	XptLevel6Interrupt   = 30
	XptLevel7Interrupt   = 31

	XptTrapInstruction = 32
	XptUserInterrupts  = 64
)

type IllegalAddressError struct {
	address uint32
}

func (e IllegalAddressError) Error() string {
	result := fmt.Sprintf("Failed to access address %08x", e.address)
	glog.Error(result)
	return result
}

type AddressHandler interface {
	Read(o *operand, a uint32) (v uint32, err error)
	Write(o *operand, a, v uint32) error
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

func (mem *MemoryHandler) Read(o *operand, a uint32) (uint32, error) {
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

func (mem *MemoryHandler) Write(o *operand, a, v uint32) error {
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
