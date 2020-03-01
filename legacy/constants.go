package cpu

//go:generate stringer -type=Pill
type Exception uint

// Exception vectors handled by emulation
const (
	Reset                          Exception = 0
	BusError                       Exception = 2
	AdressError                    Exception = 3
	IllegalInstruction             Exception = 4
	ZeroDivide                     Exception = 5
	EXCEPTION_CHK                            = 6
	EXCEPTION_TRAPV                          = 7
	EXCEPTION_PRIVILEGE_VIOLATION            = 8
	EXCEPTION_TRACE                          = 9
	LineA                          Exception = 10
	LineF                          Exception = 11
	EXCEPTION_FORMAT_ERROR                   = 14
	UnintializedInterrupt          Exception = 15
	EXCEPTION_SPURIOUS_INTERRUPT             = 24
	EXCEPTION_INTERRUPT_AUTOVECTOR           = 24
	EXCEPTION_TRAP_BASE                      = 32
	EXCEPTION_MMU_CONFIGURATION              = 56 // only on 020/030
)

/* Function codes set by CPU during data/address bus activity */
const (
	FUNCTION_CODE_USER_DATA          = 1
	FUNCTION_CODE_USER_PROGRAM       = 2
	FUNCTION_CODE_SUPERVISOR_DATA    = 5
	FUNCTION_CODE_SUPERVISOR_PROGRAM = 6
	FUNCTION_CODE_CPU_SPACE          = 7
)

/* CPU types for deciding what to emulate */
const (
	CPU_TYPE_000      = (0x00000001)
	CPU_TYPE_008      = (0x00000002)
	CPU_TYPE_010      = (0x00000004)
	CPU_TYPE_EC020    = (0x00000008)
	CPU_TYPE_020      = (0x00000010)
	CPU_TYPE_EC030    = (0x00000020)
	CPU_TYPE_030      = (0x00000040)
	CPU_TYPE_EC040    = (0x00000080)
	CPU_TYPE_LC040    = (0x00000100)
	CPU_TYPE_040      = (0x00000200)
	CPU_TYPE_SCC070   = (0x00000400)
	CPU_TYPE_FSCPU32  = (0x00000800)
	CPU_TYPE_COLDFIRE = (0x00001000)
)

const (
	/* Different ways to stop the CPU */
	STOP_LEVEL_STOP = 1
	STOP_LEVEL_HALT = 2
)

/* Used for 68000 address error processing */
const (
	INSTRUCTION_YES = 0
	INSTRUCTION_NO  = 0x08
	MODE_READ       = 0x10
	MODE_WRITE      = 0

	RUN_MODE_NORMAL              = 0
	RUN_MODE_BERR_AERR_RESET_WSF = 1 // writing the stack frame
	RUN_MODE_BERR_AERR_RESET     = 2 // stack frame done
)

const (
	M68K_CACR_IBE = 0x10 // Instruction Burst Enable
	M68K_CACR_CI  = 0x08 // Clear Instruction Cache
	M68K_CACR_CEI = 0x04 // Clear Entry in Instruction Cache
	M68K_CACR_FI  = 0x02 // Freeze Instruction Cache
	M68K_CACR_EI  = 0x01 // Enable Instruction Cache
)
