package core

import (
	"fmt"
	"sync"
)

type (
	AddressArea struct {
		read   Reader
		write  Writer
		reset  func()
		offset uint32
		size   uint32
	}

	AddressSpace struct {
		table []*AddressArea
	}
)

func NewBasePage(ssp, pc, size uint32) *AddressArea {
	// initialize base page physical memory
	mem := make([]byte, size)
	Long.Write(ssp, mem[0:])
	Long.Write(pc, mem[4:])
	// create address area
	return NewArea(0, size,
		func(offset uint32, o Operand) uint32 {
			return o.Read(mem[offset:])
		},
		func(offset uint32, o Operand, v uint32) {
			if offset < 8 {
				panic(BusError(offset)) // do not allow modification of first 2 longs
			}
			o.Write(v, mem[offset:])
		},
		func() {
			clear(mem[8:])
		})
}

func NewRAM(offset, size uint32) *AddressArea {
	mem := make([]byte, size)

	return NewArea(offset, size,
		func(offset uint32, o Operand) uint32 {
			return o.Read(mem[offset:])
		},
		func(offset uint32, o Operand, v uint32) {
			o.Write(v, mem[offset:])
		},
		func() {
			clear(mem)
		})
}

func NewROM(offset uint32, rom []byte) *AddressArea {
	return NewArea(offset, uint32(len(rom)),
		func(offset uint32, o Operand) uint32 {
			return o.Read(rom[offset:])
		},
		nil, nil)
}

func NewArea(offset, size uint32, reader Reader, writer Writer, reset func()) *AddressArea {
	return &AddressArea{reader, writer, reset, offset, size}
}

// Read returns a value at address or panics otherwise with a BusError
func (as *AddressSpace) Read(address uint32, o Operand) uint32 {
	if area := as.table[address>>PageShift]; area != nil {
		if read := area.read; read != nil {
			return read(address-area.offset, o)
		}
	}
	panic(BusError(address))
}

// Write writes a value to address or panics a BusError
func (as *AddressSpace) Write(address uint32, o Operand, value uint32) {
	if area := as.table[address>>PageShift]; area != nil {
		if write := area.write; write != nil {
			write(address-area.offset, o, value)
			return
		}
	}
	panic(BusError(address))
}

// Reset all areas
func (as *AddressSpace) Reset() {
	var wg sync.WaitGroup
	var prevArea *AddressArea = nil
	for _, area := range as.table {
		if area != nil && area != prevArea {
			prevArea = area
			if area.reset != nil {
				wg.Add(1)
				go func(a *AddressArea) {
					defer wg.Done()
					a.reset()
				}(area)
			}
		}
	}
	wg.Wait()
}

func NewAddressSpace() *AddressSpace {
	return &AddressSpace{table: make([]*AddressArea, 1<<(32-PageShift))}
}

func Allocate(as *AddressSpace, aa *AddressArea) *AddressSpace {
	if aa.offset&(PageSize-1) != 0 {
		panic(fmt.Errorf("offset $%x must be multiple of %d", aa.offset, PageSize))
	}
	if aa.size == 0 {
		panic(fmt.Errorf("size must not be 0"))
	}
	if aa.size&(PageSize-1) != 0 {
		panic(fmt.Errorf("size $%x must be multiple of %d", aa.size, PageSize))
	}

	for i := aa.offset >> PageShift; i < (aa.offset+aa.size)>>PageShift; i++ {
		if as.table[i] != nil {
			panic(fmt.Errorf("memory at offset $%x already allocated", aa.offset))
		}
		as.table[i] = aa
	}
	return as
}

func clear(slice []byte) {
	for i := range slice {
		slice[i] = 0
	}
}
