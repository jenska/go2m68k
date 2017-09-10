package m68k

type (
	Address uint32
	Long    uint32
	Word    uint16
	Byte    uint8

	CPUType   uint8
	IRQ       uint8
	IntVector uint32
	Register  uint8
)

//go:generate golang/x/tools/cmd/stringer -type=CPUType
/* CPU types for use in setCpuType() */
const (
	CPUTypeInvalid CPUType = iota
	CPUType68000
	CPUType68010
	CPUType68EC020
	CPUType680020
)

/* There are 7 levels of interrupt to the 68K.
 * A transition from < 7 to 7 will cause a non-maskable interrupt (NMI).
 */
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

/* Registers used by getReg() and setReg() */
const (
	/* Real registers */
	RegD0 Register = iota /* Data registers */
	RegD1
	RegD2
	RegD3
	RegD4
	RegD5
	RegD6
	RegD7
	RegA0 /* Address registers */
	RegA1
	RegA2
	RegA3
	RegA4
	RegA5
	RegA6
	RegA7
	RegPC   /* Program Counter */
	RegSR   /* Status Register */
	RegSP   /* The current Stack Pointer (located in A7) */
	RegUSP  /* User Stack Pointer */
	RegISP  /* Interrupt Stack Pointer */
	RegMSP  /* Master Stack Pointer */
	RegSFC  /* Source Function Code */
	RegDFC  /* Destination Function Code */
	RegVBR  /* Vector Base Register */
	RegCACR /* Cache Control Register */
	RegCAAR /* Cache Address Register */

	/* Assumed registers */
	/* These are cheat registers which emulate the 1-longword prefetch
	 * present in the 68000 and 68010.
	 */
	RegPrefAddr /* Last prefetch address */
	RegPrefData /* Last prefetch data */

	/* Convenience registers */
	RegPPC     /* Previous value in the program counter */
	RegIR      /* Instruction register */
	RegCPUType /* Type of CPU being run */
)

/* Special interrupt acknowledge values.
 * Use these as special returns from the interrupt acknowledge callback
 * (specified later in this header).
 */
const (
	/* Causes an interrupt autovector (0x18 + interrupt level) to be taken.
	 * This happens in a real 68K if VPA or AVEC is asserted during an interrupt
	 * acknowledge cycle instead of DTACK.
	 */
	IntAckAutovector IntVector = 0xffffffff

	/* Causes the spurious interrupt vector (0x18) to be taken
	 * This happens in a real 68K if BERR is asserted during the interrupt
	 * acknowledge cycle (i.e. no devices responded to the acknowledge).
	 */
	IntAckSpurious IntVector = 0xfffffffe
)

type operand struct {
	size        uint32
	alignedSize uint32
	bits        uint32
	msb         uint32
	mask        uint32
	ext         string
	formatter   string
}

var (
	byteOp = &operand{1, 2, 8, 0x80, 0xff, ".b", "%02x"}
	wordOp = &operand{2, 2, 16, 0x8000, 0xffff, ".w", "%04x"}
	longOp = &operand{4, 4, 32, 0x80000000, 0xffffffff, ".l", "%08x"}
)

func (o *operand) isNegative(value uint32) bool {
	return (o.msb & value) != 0
}

func (o *operand) set(target *uint32, value uint32) {
	*target = (*target & ^o.mask) | (value & o.mask)
}

func (o *operand) getSigned(value uint32) int32 {
	v := uint32(value)
	if o.isNegative(v) {
		return int32(v | ^o.mask)
	}
	return int32(v & o.mask)
}

func (o *operand) get(value uint32) uint32 {
	return value & o.mask
}
