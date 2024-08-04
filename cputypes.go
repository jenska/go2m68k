package m68k

import (
	"github.com/jenska/m68k/internal/core"
)

func M68000(bc BusController) *core.Core {
	return core.Build(
		func(address uint32, o core.Operand) uint32 {
			address &= core.Address24Mask
			if o == core.Byte {
				return bc.as.Read(address, core.Byte)
			} else if address&1 == 0 {
				return bc.as.Read(address, o)
			}
			panic(core.AddressError(address))
		},
		func(address uint32, o core.Operand, v uint32) {
			address &= core.Address24Mask
			if o == core.Byte {
				bc.as.Write(address, core.Byte, v)
			} else if address&1 == 0 {
				bc.as.Write(address, o, v)
			} else {
				panic(core.AddressError(address))
			}
		},
		bc.as.Reset)(core.M68000InstructionSet)
}

func M68020(bc BusController) *core.Core {
	return core.Build(
		func(address uint32, o core.Operand) uint32 {
			return bc.as.Read(address, o)
		},
		func(address uint32, o core.Operand, v uint32) {
			bc.as.Write(address, o, v)
		},
		bc.as.Reset)(core.M68000InstructionSet)
}
