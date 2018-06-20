package cpu

import (
	"github.com/jenska/atari2go/mem"
)

type m68000Base struct {
	hasFPU  bool // Indicates if a FPU is available (yes on 030, 040, may be on 020)
	cpuType CPUType

	executeMinCycles        uint32
	executeMaxCycles        uint32
	executeInputLines       uint32 // number of input lines
	executeDefaultIrqVector uint32

	dar          [16]uint32  // Data and Address Registers
	ppc          uint32      // Previous program counter
	pc           uint32      // Program Counter
	sp           [7]uint32   // User, Interrupt, and Master Stack Pointers
	vbr          uint32      // Vector Base Register (m68010+)
	sfc          uint32      // Source Function Code Register (m68010+)
	dfc          uint32      // Destination Function Code Register (m68010+)
	cacr         uint32      // Cache Control Register (m68020, unemulated)
	caar         uint32      // Cache Address Register (m68020, unemulated)
	ir           uint32      // Instruction Register
	fpr          [8]floatx80 // FPU Data Register (m68030/040)
	fpiar        uint32      // FPU Instruction Address Register (m68040)
	fpsr         uint32      // FPU Status Register (m68040)
	fpcr         uint32      // FPU Control Register (m68040)
	t1Flag       uint32      // Trace 1
	t0Flag       uint32      // Trace 0
	sFlag        uint32      // Supervisor
	mFlag        uint32      // Master/Interrupt state
	xFlag        uint32      // Extend
	nFlag        uint32      // Negative
	notZFlag     uint32      // Zero, inverted for speedups
	vFlag        uint32      // Overflow
	cFlag        uint32      // Carry
	intMask      uint32      // I0-I2
	intLevel     uint32      // State of interrupt pins IPL0-IPL2 -- ASG: changed from ints_pending
	stopped      uint32      // Stopped state
	prefAddr     uint32      // Last prefetch address
	prefData     uint32      // Data in the prefetch queue
	srMask       uint32      // Implemented status register bits
	instrMode    uint32      // Stores whether we are in instruction mode or group 0/1 exception mode
	runMode      uint32      // Stores whether we are processing a reset, bus error, address error, or something else
	hasPmmu      bool        // Indicates if a PMMU available (yes on 030, 040, no on EC030)
	pmmuEnabled  bool        // Indicates if the PMMU is enabled
	fpuJustReset bool        // Indicates the FPU was just reset

	// Clocks required for instructions / exceptions
	cyBccNotakeB  uint32
	cycBccNotakeE uint32
	cycDbccFNoexp uint32
	cycDbccFRxp   uint32
	cycSccRTrue   uint32
	cycMovemW     uint32
	cycMovemL     uint32
	cycShift      uint32
	cycReset      uint32

	initialCycles   int
	remainingCycles int // Number of clocks remaining
	resetCycles     int
	tracing         uint32

	addressError int

	aerrAddress   uint32
	aerrWriteMode uint32
	aerrFc        uint32

	// Virtual IRQ lines state
	virqState  uint32
	nmiPending uint32

	cycInstruction *uint8
	cycException   *uint8

	intAckCallback     func(*M68K, int) int                         // Interrupt Acknowledge
	bkptAckCallback    func(*mem.AddressSpace, int, uint32, uint32) // Breakpoint Acknowledge
	resetInstrCallback func(int)                                    // Called when a RESET instruction is encountered
	rteInstrCallback   func(int)                                    // Called when a RTE instruction is encountered
	cmpilInstrCallback func(*mem.AddressSpace, int, uint32, uint32) // Called when a CMPI.L #v, Dn instruction is encountered
	tasWriteCallback   func(*mem.AddressSpace, int, uint8, uint8)   // Called instead of normal write8 by the TAS instruction,

	program  *mem.AddressSpace
	oprogram *mem.AddressSpace
}

type M68K interface {
	presave()
	postload()
	clearAll()

	createDisassembler() *dasm // CPUType dependend disassembler

}

func InitM68K(cpuType CPUType) *M68K {

}
