package m68k

import (
	"fmt"
	"strings"

	"github.com/jenska/m68k/internal/core"
	_ "github.com/jenska/m68k/internal/instructions"
)

const (
	M68000 Model = 1 << iota
	M68010

	HaltSignal = iota
	ResetSignal
	Int1Signal
	Int2Signal
	Int3Signal
	Int4Signal
	Int5Signal
	Int6Singal
	Int7Signal
)

type (
	Model uint16 // CPU model
	FC    uint16 // Function code

	// CPU Interface
	CPU interface {
		Execute(signals <-chan uint16)
	}

	// Disassembles a single instruction and returns the instruction size.
	Disassemble func(addr uint32, b strings.Builder) uint16

	Reader interface {
		Read8(address uint32) uint8
		Read16(address uint32) uint16
		Read32(address uint32) uint32
	}

	Writer interface {
		Write8(address uint32, value uint8)
		Write16(address uint32, value uint16)
		Write32(address uint32, value uint32)
	}

	MemoryArea func() *core.AddressArea

	// BusController encapsulates the memory areas.
	BusController struct {
		as *core.AddressSpace
	}
)

// AddArea adds a chip memory area to the bus controllers scope
func ChipArea(offset, size uint32, reader Reader, writer Writer, reset func()) MemoryArea {
	return func() *core.AddressArea {
		return core.NewArea(offset, size,
			func(offset uint32, o core.Operand) uint32 {
				switch o {
				case core.Byte:
					return uint32(reader.Read8(offset))
				case core.Word:
					return uint32(reader.Read16(offset))
				case core.Long:
					return reader.Read32(offset)
				default:
					panic("invalid operand size")
				}
			},
			func(offset uint32, o core.Operand, v uint32) {
				switch o {
				case core.Byte:
					writer.Write8(offset, uint8(v))
				case core.Word:
					writer.Write16(offset, uint16(v))
				case core.Long:
					writer.Write32(offset, v)
				default:
					panic("invalid operand size")
				}
			},
			reset)
	}
}

func BaseRAM(ssp, pc, size uint32) MemoryArea {
	return func() *core.AddressArea {
		return core.NewBasePage(ssp, pc, size)
	}
}

func RAM(offset, size uint32) MemoryArea {
	return func() *core.AddressArea {
		return core.NewRAM(offset, size)
	}
}

func ROM(offset uint32, rom []byte) MemoryArea {
	return func() *core.AddressArea {
		return core.NewROM(offset, rom)
	}
}

func NewBusController(memoryMap ...MemoryArea) BusController {
	bc := BusController{core.NewAddressSpace()}
	for _, area := range memoryMap {
		core.Allocate(bc.as, area())
	}
	return bc
}

func New(m Model, bc BusController) CPU {
	switch m {
	case M68000:
		return core.Build(
			func(address uint32, o core.Operand) uint32 {
				address &= core.Address24Mask
				if address&1 == 0 {
					return bc.as.Read(address, o)
				} else if o == core.Byte {
					return bc.as.Read(address, core.Byte)
				}
				panic(core.AddressError(address))
			},
			func(address uint32, o core.Operand, v uint32) {
				address &= core.Address24Mask
				if address&1 == 0 {
					bc.as.Write(address, o, v)
				} else if o == core.Byte {
					bc.as.Write(address, core.Byte, v)
				} else {
					panic(core.AddressError(address))
				}
			},
			bc.as.Reset)(core.M68000InstructionSet)
	default:
		panic(fmt.Errorf("unsupported CPU model"))
	}
}
