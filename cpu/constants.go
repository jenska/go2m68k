package cpu

type (
	Address uint32
	Long    uint32
	Word    uint16
	Byte    uint8

	CPUType   uint32
	IRQ       uint8
	IntVector uint32
	Register  uint8
)

const (
	CPUTypeInvalid CPUType = 1 << iota
	CPUType68000
	CPUType68008
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

const (
	maskAll        = (CPUType68000 | CPUType68008 | CPUType68010 | CPUType68EC020 | CPUType68020 | CPUType68EC030 | CPUType68030 | CPUType68EC040 | CPUType68040 | CPUTypeFSCPU32)
	mask24BitSpace = (CPUType68000 | CPUType68008 | CPUType68010 | CPUType68EC020)
	mask32BitSpace = (CPUType68020 | CPUType68EC030 | CPUType68030 | CPUType68EC040 | CPUType68040 | CPUTypeFSCPU32)
	mask010orLater = (CPUType68010 | CPUType68EC020 | CPUType68020 | CPUType68030 | CPUType68EC030 | CPUType68040 | CPUType68EC040 | CPUTypeFSCPU32)
	mask020orLater = (CPUType68EC020 | CPUType68020 | CPUType68EC030 | CPUType68030 | CPUType68EC040 | CPUType68040 | CPUTypeFSCPU32)
	mask030orLater = (CPUType68030 | CPUType68EC030 | CPUType68040 | CPUType68EC040)
	mask040orLater = (CPUType68040 | CPUType68EC040)
)

// Registers used by getReg() and setReg()
//go:generate stringer -type=Register
const (
	// Real registers
	RegD0 Register = iota // Data registers
	RegD1
	RegD2
	RegD3
	RegD4
	RegD5
	RegD6
	RegD7
	RegA0 // Address registers
	RegA1
	RegA2
	RegA3
	RegA4
	RegA5
	RegA6
	RegA7
	RegPC   // Program Counter
	RegSR   // Status Register
	RegSP   // The current Stack Pointer (located in A7)
	RegUSP  // User Stack Pointer
	RegISP  // Interrupt Stack Pointer
	RegMSP  // Master Stack Pointer
	RegSFC  // Source Function Code
	RegDFC  // Destination Function Code
	RegVBR  // Vector Base Register
	RegCACR // Cache Control Register
	RegCAAR // Cache Address Register

	// Assumed registers
	// These are cheat registers which emulate the 1-longword prefetch
	// present in the 68000 and 68010.
	RegPrefAddr // Last prefetch address
	RegPrefData // Last prefetch data

	// Convenience registers
	RegPPC     // Previous value in the program counter
	RegIR      // Instruction register
	RegCPUType // Type of CPU being run
)

// Special interrupt acknowledge values.
// Use these as special returns from the interrupt acknowledge callback
// (specified later in this header).

//go:generate stringer -type=IntVector
const (
	// Causes an interrupt autovector (0x18 + interrupt level) to be taken.
	// This happens in a real 68K if VPA or AVEC is asserted during an interrupt
	// acknowledge cycle instead of DTACK.

	IntAckAutovector IntVector = 0xffffffff

	// Causes the spurious interrupt vector (0x18) to be taken
	// This happens in a real 68K if BERR is asserted during the interrupt
	// acknowledge cycle (i.e. no devices responded to the acknowledge).

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

// Used by shift & rotate instructions
var (
	shiftByteTable = []Byte{
		0x00, 0x80, 0xc0, 0xe0, 0xf0, 0xf8, 0xfc, 0xfe, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff}
	shiftWordTable = []Word{
		0x0000, 0x8000, 0xc000, 0xe000, 0xf000, 0xf800, 0xfc00, 0xfe00, 0xff00,
		0xff80, 0xffc0, 0xffe0, 0xfff0, 0xfff8, 0xfffc, 0xfffe, 0xffff, 0xffff,
		0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff,
		0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff,
		0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff,
		0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff,
		0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff,
		0xffff, 0xffff}
	shiftLongTable = []Long{
		0x00000000, 0x80000000, 0xc0000000, 0xe0000000, 0xf0000000, 0xf8000000,
		0xfc000000, 0xfe000000, 0xff000000, 0xff800000, 0xffc00000, 0xffe00000,
		0xfff00000, 0xfff80000, 0xfffc0000, 0xfffe0000, 0xffff0000, 0xffff8000,
		0xffffc000, 0xffffe000, 0xfffff000, 0xfffff800, 0xfffffc00, 0xfffffe00,
		0xffffff00, 0xffffff80, 0xffffffc0, 0xffffffe0, 0xfffffff0, 0xfffffff8,
		0xfffffffc, 0xfffffffe, 0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff,
		0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff,
		0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff,
		0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff,
		0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff,
		0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff}

	// Number of clock cycles to use for exception processing.
	m68kiExceptionCycleTable = [7][256]uint8{
		{ // 000
			40, //  0: Reset - Initial Stack Pointer
			4,  //  1: Reset - Initial Program Counter
			50, //  2: Bus Error                             (unemulated)
			50, //  3: Address Error                         (unemulated)
			34, //  4: Illegal Instruction
			38, //  5: Divide by Zero
			40, //  6: CHK
			34, //  7: TRAPV
			34, //  8: Privilege Violation
			34, //  9: Trace
			4,  // 10: 1010
			4,  // 11: 1111
			4,  // 12: RESERVED
			4,  // 13: Coprocessor Protocol Violation        (unemulated)
			4,  // 14: Format Error
			44, // 15: Uninitialized Interrupt
			4,  // 16: RESERVED
			4,  // 17: RESERVED
			4,  // 18: RESERVED
			4,  // 19: RESERVED
			4,  // 20: RESERVED
			4,  // 21: RESERVED
			4,  // 22: RESERVED
			4,  // 23: RESERVED
			44, // 24: Spurious Interrupt
			44, // 25: Level 1 Interrupt Autovector
			44, // 26: Level 2 Interrupt Autovector
			44, // 27: Level 3 Interrupt Autovector
			44, // 28: Level 4 Interrupt Autovector
			44, // 29: Level 5 Interrupt Autovector
			44, // 30: Level 6 Interrupt Autovector
			44, // 31: Level 7 Interrupt Autovector
			34, // 32: TRAP #0
			34, // 33: TRAP #1
			34, // 34: TRAP #2
			34, // 35: TRAP #3
			34, // 36: TRAP #4
			34, // 37: TRAP #5
			34, // 38: TRAP #6
			34, // 39: TRAP #7
			34, // 40: TRAP #8
			34, // 41: TRAP #9
			34, // 42: TRAP #10
			34, // 43: TRAP #11
			34, // 44: TRAP #12
			34, // 45: TRAP #13
			34, // 46: TRAP #14
			34, // 47: TRAP #15
			4,  // 48: FP Branch or Set on Unknown Condition (unemulated)
			4,  // 49: FP Inexact Result                     (unemulated)
			4,  // 50: FP Divide by Zero                     (unemulated)
			4,  // 51: FP Underflow                          (unemulated)
			4,  // 52: FP Operand Error                      (unemulated)
			4,  // 53: FP Overflow                           (unemulated)
			4,  // 54: FP Signaling NAN                      (unemulated)
			4,  // 55: FP Unimplemented Data Type            (unemulated)
			4,  // 56: MMU Configuration Error               (unemulated)
			4,  // 57: MMU Illegal Operation Error           (unemulated)
			4,  // 58: MMU Access Level Violation Error      (unemulated)
			4,  // 59: RESERVED
			4,  // 60: RESERVED
			4,  // 61: RESERVED
			4,  // 62: RESERVED
			4,  // 63: RESERVED
			// 64-255: User Defined
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
		},
		{ // 010
			40,  //  0: Reset - Initial Stack Pointer
			4,   //  1: Reset - Initial Program Counter
			126, //  2: Bus Error                             (unemulated)
			126, //  3: Address Error                         (unemulated)
			38,  //  4: Illegal Instruction
			44,  //  5: Divide by Zero
			44,  //  6: CHK
			34,  //  7: TRAPV
			38,  //  8: Privilege Violation
			38,  //  9: Trace
			4,   // 10: 1010
			4,   // 11: 1111
			4,   // 12: RESERVED
			4,   // 13: Coprocessor Protocol Violation        (unemulated)
			4,   // 14: Format Error
			44,  // 15: Uninitialized Interrupt
			4,   // 16: RESERVED
			4,   // 17: RESERVED
			4,   // 18: RESERVED
			4,   // 19: RESERVED
			4,   // 20: RESERVED
			4,   // 21: RESERVED
			4,   // 22: RESERVED
			4,   // 23: RESERVED
			46,  // 24: Spurious Interrupt
			46,  // 25: Level 1 Interrupt Autovector
			46,  // 26: Level 2 Interrupt Autovector
			46,  // 27: Level 3 Interrupt Autovector
			46,  // 28: Level 4 Interrupt Autovector
			46,  // 29: Level 5 Interrupt Autovector
			46,  // 30: Level 6 Interrupt Autovector
			46,  // 31: Level 7 Interrupt Autovector
			38,  // 32: TRAP #0
			38,  // 33: TRAP #1
			38,  // 34: TRAP #2
			38,  // 35: TRAP #3
			38,  // 36: TRAP #4
			38,  // 37: TRAP #5
			38,  // 38: TRAP #6
			38,  // 39: TRAP #7
			38,  // 40: TRAP #8
			38,  // 41: TRAP #9
			38,  // 42: TRAP #10
			38,  // 43: TRAP #11
			38,  // 44: TRAP #12
			38,  // 45: TRAP #13
			38,  // 46: TRAP #14
			38,  // 47: TRAP #15
			4,   // 48: FP Branch or Set on Unknown Condition (unemulated)
			4,   // 49: FP Inexact Result                     (unemulated)
			4,   // 50: FP Divide by Zero                     (unemulated)
			4,   // 51: FP Underflow                          (unemulated)
			4,   // 52: FP Operand Error                      (unemulated)
			4,   // 53: FP Overflow                           (unemulated)
			4,   // 54: FP Signaling NAN                      (unemulated)
			4,   // 55: FP Unimplemented Data Type            (unemulated)
			4,   // 56: MMU Configuration Error               (unemulated)
			4,   // 57: MMU Illegal Operation Error           (unemulated)
			4,   // 58: MMU Access Level Violation Error      (unemulated)
			4,   // 59: RESERVED
			4,   // 60: RESERVED
			4,   // 61: RESERVED
			4,   // 62: RESERVED
			4,   // 63: RESERVED
			// 64-255: User Defined
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
		},
		{ // 020
			4,  //  0: Reset - Initial Stack Pointer
			4,  //  1: Reset - Initial Program Counter
			50, //  2: Bus Error                             (unemulated)
			50, //  3: Address Error                         (unemulated)
			20, //  4: Illegal Instruction
			38, //  5: Divide by Zero
			40, //  6: CHK
			20, //  7: TRAPV
			34, //  8: Privilege Violation
			25, //  9: Trace
			20, // 10: 1010
			20, // 11: 1111
			4,  // 12: RESERVED
			4,  // 13: Coprocessor Protocol Violation        (unemulated)
			4,  // 14: Format Error
			30, // 15: Uninitialized Interrupt
			4,  // 16: RESERVED
			4,  // 17: RESERVED
			4,  // 18: RESERVED
			4,  // 19: RESERVED
			4,  // 20: RESERVED
			4,  // 21: RESERVED
			4,  // 22: RESERVED
			4,  // 23: RESERVED
			30, // 24: Spurious Interrupt
			30, // 25: Level 1 Interrupt Autovector
			30, // 26: Level 2 Interrupt Autovector
			30, // 27: Level 3 Interrupt Autovector
			30, // 28: Level 4 Interrupt Autovector
			30, // 29: Level 5 Interrupt Autovector
			30, // 30: Level 6 Interrupt Autovector
			30, // 31: Level 7 Interrupt Autovector
			20, // 32: TRAP #0
			20, // 33: TRAP #1
			20, // 34: TRAP #2
			20, // 35: TRAP #3
			20, // 36: TRAP #4
			20, // 37: TRAP #5
			20, // 38: TRAP #6
			20, // 39: TRAP #7
			20, // 40: TRAP #8
			20, // 41: TRAP #9
			20, // 42: TRAP #10
			20, // 43: TRAP #11
			20, // 44: TRAP #12
			20, // 45: TRAP #13
			20, // 46: TRAP #14
			20, // 47: TRAP #15
			4,  // 48: FP Branch or Set on Unknown Condition (unemulated)
			4,  // 49: FP Inexact Result                     (unemulated)
			4,  // 50: FP Divide by Zero                     (unemulated)
			4,  // 51: FP Underflow                          (unemulated)
			4,  // 52: FP Operand Error                      (unemulated)
			4,  // 53: FP Overflow                           (unemulated)
			4,  // 54: FP Signaling NAN                      (unemulated)
			4,  // 55: FP Unimplemented Data Type            (unemulated)
			4,  // 56: MMU Configuration Error               (unemulated)
			4,  // 57: MMU Illegal Operation Error           (unemulated)
			4,  // 58: MMU Access Level Violation Error      (unemulated)
			4,  // 59: RESERVED
			4,  // 60: RESERVED
			4,  // 61: RESERVED
			4,  // 62: RESERVED
			4,  // 63: RESERVED
			// 64-255: User Defined
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
		},
		{ // 030 - not correct
			4,  //  0: Reset - Initial Stack Pointer
			4,  //  1: Reset - Initial Program Counter
			50, //  2: Bus Error                             (unemulated)
			50, //  3: Address Error                         (unemulated)
			20, //  4: Illegal Instruction
			38, //  5: Divide by Zero
			40, //  6: CHK
			20, //  7: TRAPV
			34, //  8: Privilege Violation
			25, //  9: Trace
			20, // 10: 1010
			20, // 11: 1111
			4,  // 12: RESERVED
			4,  // 13: Coprocessor Protocol Violation        (unemulated)
			4,  // 14: Format Error
			30, // 15: Uninitialized Interrupt
			4,  // 16: RESERVED
			4,  // 17: RESERVED
			4,  // 18: RESERVED
			4,  // 19: RESERVED
			4,  // 20: RESERVED
			4,  // 21: RESERVED
			4,  // 22: RESERVED
			4,  // 23: RESERVED
			30, // 24: Spurious Interrupt
			30, // 25: Level 1 Interrupt Autovector
			30, // 26: Level 2 Interrupt Autovector
			30, // 27: Level 3 Interrupt Autovector
			30, // 28: Level 4 Interrupt Autovector
			30, // 29: Level 5 Interrupt Autovector
			30, // 30: Level 6 Interrupt Autovector
			30, // 31: Level 7 Interrupt Autovector
			20, // 32: TRAP #0
			20, // 33: TRAP #1
			20, // 34: TRAP #2
			20, // 35: TRAP #3
			20, // 36: TRAP #4
			20, // 37: TRAP #5
			20, // 38: TRAP #6
			20, // 39: TRAP #7
			20, // 40: TRAP #8
			20, // 41: TRAP #9
			20, // 42: TRAP #10
			20, // 43: TRAP #11
			20, // 44: TRAP #12
			20, // 45: TRAP #13
			20, // 46: TRAP #14
			20, // 47: TRAP #15
			4,  // 48: FP Branch or Set on Unknown Condition (unemulated)
			4,  // 49: FP Inexact Result                     (unemulated)
			4,  // 50: FP Divide by Zero                     (unemulated)
			4,  // 51: FP Underflow                          (unemulated)
			4,  // 52: FP Operand Error                      (unemulated)
			4,  // 53: FP Overflow                           (unemulated)
			4,  // 54: FP Signaling NAN                      (unemulated)
			4,  // 55: FP Unimplemented Data Type            (unemulated)
			4,  // 56: MMU Configuration Error               (unemulated)
			4,  // 57: MMU Illegal Operation Error           (unemulated)
			4,  // 58: MMU Access Level Violation Error      (unemulated)
			4,  // 59: RESERVED
			4,  // 60: RESERVED
			4,  // 61: RESERVED
			4,  // 62: RESERVED
			4,  // 63: RESERVED
			// 64-255: User Defined
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
		},
		{ // 040  // TODO: these values are not correct
			4,  //  0: Reset - Initial Stack Pointer
			4,  //  1: Reset - Initial Program Counter
			50, //  2: Bus Error                             (unemulated)
			50, //  3: Address Error                         (unemulated)
			20, //  4: Illegal Instruction
			38, //  5: Divide by Zero
			40, //  6: CHK
			20, //  7: TRAPV
			34, //  8: Privilege Violation
			25, //  9: Trace
			20, // 10: 1010
			20, // 11: 1111
			4,  // 12: RESERVED
			4,  // 13: Coprocessor Protocol Violation        (unemulated)
			4,  // 14: Format Error
			30, // 15: Uninitialized Interrupt
			4,  // 16: RESERVED
			4,  // 17: RESERVED
			4,  // 18: RESERVED
			4,  // 19: RESERVED
			4,  // 20: RESERVED
			4,  // 21: RESERVED
			4,  // 22: RESERVED
			4,  // 23: RESERVED
			30, // 24: Spurious Interrupt
			30, // 25: Level 1 Interrupt Autovector
			30, // 26: Level 2 Interrupt Autovector
			30, // 27: Level 3 Interrupt Autovector
			30, // 28: Level 4 Interrupt Autovector
			30, // 29: Level 5 Interrupt Autovector
			30, // 30: Level 6 Interrupt Autovector
			30, // 31: Level 7 Interrupt Autovector
			20, // 32: TRAP #0
			20, // 33: TRAP #1
			20, // 34: TRAP #2
			20, // 35: TRAP #3
			20, // 36: TRAP #4
			20, // 37: TRAP #5
			20, // 38: TRAP #6
			20, // 39: TRAP #7
			20, // 40: TRAP #8
			20, // 41: TRAP #9
			20, // 42: TRAP #10
			20, // 43: TRAP #11
			20, // 44: TRAP #12
			20, // 45: TRAP #13
			20, // 46: TRAP #14
			20, // 47: TRAP #15
			4,  // 48: FP Branch or Set on Unknown Condition (unemulated)
			4,  // 49: FP Inexact Result                     (unemulated)
			4,  // 50: FP Divide by Zero                     (unemulated)
			4,  // 51: FP Underflow                          (unemulated)
			4,  // 52: FP Operand Error                      (unemulated)
			4,  // 53: FP Overflow                           (unemulated)
			4,  // 54: FP Signaling NAN                      (unemulated)
			4,  // 55: FP Unimplemented Data Type            (unemulated)
			4,  // 56: MMU Configuration Error               (unemulated)
			4,  // 57: MMU Illegal Operation Error           (unemulated)
			4,  // 58: MMU Access Level Violation Error      (unemulated)
			4,  // 59: RESERVED
			4,  // 60: RESERVED
			4,  // 61: RESERVED
			4,  // 62: RESERVED
			4,  // 63: RESERVED
			// 64-255: User Defined
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
		},
		{ // CPU32
			4,  //  0: Reset - Initial Stack Pointer
			4,  //  1: Reset - Initial Program Counter
			50, //  2: Bus Error                             (unemulated)
			50, //  3: Address Error                         (unemulated)
			20, //  4: Illegal Instruction
			38, //  5: Divide by Zero
			40, //  6: CHK
			20, //  7: TRAPV
			34, //  8: Privilege Violation
			25, //  9: Trace
			20, // 10: 1010
			20, // 11: 1111
			4,  // 12: RESERVED
			4,  // 13: Coprocessor Protocol Violation        (unemulated)
			4,  // 14: Format Error
			30, // 15: Uninitialized Interrupt
			4,  // 16: RESERVED
			4,  // 17: RESERVED
			4,  // 18: RESERVED
			4,  // 19: RESERVED
			4,  // 20: RESERVED
			4,  // 21: RESERVED
			4,  // 22: RESERVED
			4,  // 23: RESERVED
			30, // 24: Spurious Interrupt
			30, // 25: Level 1 Interrupt Autovector
			30, // 26: Level 2 Interrupt Autovector
			30, // 27: Level 3 Interrupt Autovector
			30, // 28: Level 4 Interrupt Autovector
			30, // 29: Level 5 Interrupt Autovector
			30, // 30: Level 6 Interrupt Autovector
			30, // 31: Level 7 Interrupt Autovector
			20, // 32: TRAP #0
			20, // 33: TRAP #1
			20, // 34: TRAP #2
			20, // 35: TRAP #3
			20, // 36: TRAP #4
			20, // 37: TRAP #5
			20, // 38: TRAP #6
			20, // 39: TRAP #7
			20, // 40: TRAP #8
			20, // 41: TRAP #9
			20, // 42: TRAP #10
			20, // 43: TRAP #11
			20, // 44: TRAP #12
			20, // 45: TRAP #13
			20, // 46: TRAP #14
			20, // 47: TRAP #15
			4,  // 48: FP Branch or Set on Unknown Condition (unemulated)
			4,  // 49: FP Inexact Result                     (unemulated)
			4,  // 50: FP Divide by Zero                     (unemulated)
			4,  // 51: FP Underflow                          (unemulated)
			4,  // 52: FP Operand Error                      (unemulated)
			4,  // 53: FP Overflow                           (unemulated)
			4,  // 54: FP Signaling NAN                      (unemulated)
			4,  // 55: FP Unimplemented Data Type            (unemulated)
			4,  // 56: MMU Configuration Error               (unemulated)
			4,  // 57: MMU Illegal Operation Error           (unemulated)
			4,  // 58: MMU Access Level Violation Error      (unemulated)
			4,  // 59: RESERVED
			4,  // 60: RESERVED
			4,  // 61: RESERVED
			4,  // 62: RESERVED
			4,  // 63: RESERVED
			// 64-255: User Defined
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
		},
		{ // ColdFire - not correct
			4,  //  0: Reset - Initial Stack Pointer
			4,  //  1: Reset - Initial Program Counter
			50, //  2: Bus Error                             (unemulated)
			50, //  3: Address Error                         (unemulated)
			20, //  4: Illegal Instruction
			38, //  5: Divide by Zero
			40, //  6: CHK
			20, //  7: TRAPV
			34, //  8: Privilege Violation
			25, //  9: Trace
			20, // 10: 1010
			20, // 11: 1111
			4,  // 12: RESERVED
			4,  // 13: Coprocessor Protocol Violation        (unemulated)
			4,  // 14: Format Error
			30, // 15: Uninitialized Interrupt
			4,  // 16: RESERVED
			4,  // 17: RESERVED
			4,  // 18: RESERVED
			4,  // 19: RESERVED
			4,  // 20: RESERVED
			4,  // 21: RESERVED
			4,  // 22: RESERVED
			4,  // 23: RESERVED
			30, // 24: Spurious Interrupt
			30, // 25: Level 1 Interrupt Autovector
			30, // 26: Level 2 Interrupt Autovector
			30, // 27: Level 3 Interrupt Autovector
			30, // 28: Level 4 Interrupt Autovector
			30, // 29: Level 5 Interrupt Autovector
			30, // 30: Level 6 Interrupt Autovector
			30, // 31: Level 7 Interrupt Autovector
			20, // 32: TRAP #0
			20, // 33: TRAP #1
			20, // 34: TRAP #2
			20, // 35: TRAP #3
			20, // 36: TRAP #4
			20, // 37: TRAP #5
			20, // 38: TRAP #6
			20, // 39: TRAP #7
			20, // 40: TRAP #8
			20, // 41: TRAP #9
			20, // 42: TRAP #10
			20, // 43: TRAP #11
			20, // 44: TRAP #12
			20, // 45: TRAP #13
			20, // 46: TRAP #14
			20, // 47: TRAP #15
			4,  // 48: FP Branch or Set on Unknown Condition (unemulated)
			4,  // 49: FP Inexact Result                     (unemulated)
			4,  // 50: FP Divide by Zero                     (unemulated)
			4,  // 51: FP Underflow                          (unemulated)
			4,  // 52: FP Operand Error                      (unemulated)
			4,  // 53: FP Overflow                           (unemulated)
			4,  // 54: FP Signaling NAN                      (unemulated)
			4,  // 55: FP Unimplemented Data Type            (unemulated)
			4,  // 56: MMU Configuration Error               (unemulated)
			4,  // 57: MMU Illegal Operation Error           (unemulated)
			4,  // 58: MMU Access Level Violation Error      (unemulated)
			4,  // 59: RESERVED
			4,  // 60: RESERVED
			4,  // 61: RESERVED
			4,  // 62: RESERVED
			4,  // 63: RESERVED
			// 64-255: User Defined
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
		},
	}

	m68kiEAIdxCycleTable = [64]uint8{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, // ..01.000 no memory indirect, base nullptr
		5, // ..01..01 memory indirect,    base nullptr, outer nullptr
		7, // ..01..10 memory indirect,    base nullptr, outer 16
		7, // ..01..11 memory indirect,    base nullptr, outer 32
		0, 5, 7, 7, 0, 5, 7, 7, 0, 5, 7, 7,
		2, // ..10.000 no memory indirect, base 16
		7, // ..10..01 memory indirect,    base 16,   outer nullptr
		9, // ..10..10 memory indirect,    base 16,   outer 16
		9, // ..10..11 memory indirect,    base 16,   outer 32
		0, 7, 9, 9, 0, 7, 9, 9, 0, 7, 9, 9,
		6,  // ..11.000 no memory indirect, base 32
		11, // ..11..01 memory indirect,    base 32,   outer nullptr
		13, // ..11..10 memory indirect,    base 32,   outer 16
		13, // ..11..11 memory indirect,    base 32,   outer 32
		0, 11, 13, 13, 0, 11, 13, 13, 0, 11, 13, 13,
	}
)
