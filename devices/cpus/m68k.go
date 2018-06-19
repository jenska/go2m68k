package cpus

import (
	"github.com/jenska/atari2go/devices"
)

type (
	IRQ     uint8
	CPUType uint8
)

const (
	MMUATCEntries        = 22  // 68851 has 64, 030 has 22
	InstructionCacheSize = 128 // instruction cache constants
	LineBusError         = 16  // special input lines
)

/* There are 7 levels of interrupt to the 68K.
 * A transition from < 7 to 7 will cause a non-maskable interrupt (NMI).
 */
//go:generate stringer -type=IRQ
const (
	IRQNone IRQ = iota
	IRQ1
	IRQ2
	IRQ3
	IRQ4
	IRQ5
	IRQ6
	IRQ7
)

//go:generate stringer -type=CPUType
const (
	CPUTypeInvalid CPUType = iota
	CPUType68000
	CPUType68010
	CPUType68EC020
	CPUType68020
	CPUType68EC030
	CPUType68030
	CPUType68EC040
	CPUType68LC040
	CPUType68040
	CPUTypeSCC68070
	CPUTypeFSCPU32
	CPUTypeColdfire
)

// M68000Base structure for all 68000 CPU types
type M68000Base struct {
	devices.CPU_base
}
