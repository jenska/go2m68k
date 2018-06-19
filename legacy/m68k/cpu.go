package m68k

type CPU struct {
	cpuType CPUType    /* CPU Type: 68000, 68010, 68EC020, or 68020 */
	dr      [8]Long    /* Data Registers */
	ar      [8]Address /* Address Registers */
	ppc     Address    /* Previous program counter */
	pc      Address    /* Program Counter */
	sp      [7]Address /* User, Interrupt, and Master Stack Pointers */
	vbr     Long       /* Vector Base Register (m68010+) */
	sfc     Long
	dfc     Long
	cacr    Long
	caar    Long
	ir      Long
	sr      StatusRegister
}

/* Pulse the RESET pin on the CPU.
 * You *MUST* reset the CPU at least once to initialize the emulation
 * Note: If you didn't call m68k_set_cpu_type() before resetting
 *       the CPU for the first time, the CPU will be set to
 *       CPUType68000.
 */
func (core *CPU) PulseReset() {

}

/* execute numCycles worth of instructions.  returns number of cycles used */
func (core *CPU) Execute(numCycles Long) Long {
	return 0
}

/* Number of cycles run so far */
func (core *CPU) Cycles() Long {
	return 0
}

func (core *CPU) CyclesRemaining() Long {
	return 0
}

func (core *CPU) ModifyTimeslice(cycles Long) {

}

func (core *CPU) EndTimeslice() {

}

/* Set the IPL0-IPL2 pins on the CPU (IRQ).
 * A transition from < 7 to 7 will cause a non-maskable interrupt (NMI).
 * Setting IRQ to 0 will clear an interrupt request.
 */
func (core *CPU) SetIrq(intLevel IRQ) {}

/* Halt the CPU as if you pulsed the HALT pin. */
func (core *CPU) PulseHalt() {}

/* Peek at the internals of a CPU context.  This can either be a context
 * retrieved using m68k_get_context() or the currently running context.
 * If context is NULL, the currently running CPU context will be used.
 */
func (core *CPU) GetReg(reg Register) {}

/* Poke values into the internals of the currently running CPU context */
func (core *CPU) SetReg(reg Register, value Long) {}

/* Check if an instruction is valid for the specified CPU type */
func (core *CPU) IsValidInstruction(instruction uint16, cpuType CPUType) bool {
	return false
}

/* Disassemble 1 instruction using the epecified CPU type at pc.  Stores
 * disassembly in str_buff and returns the size of the instruction in bytes.
 */
func (core *CPU) Disassemble(pc Address, cpuType CPUType) (line string, instructionSize int) {
	return "?", 0
}

/* IntAckCallback will be called with the interrupt level being acknowledged.
 * The host program must return either a vector from 0x02-0xff, or one of the
 * special interrupt acknowledge values specified earlier in this header.
 * If this is not implemented, the CPU will always assume an autovectored
 * interrupt, and will automatically clear the interrupt request when it
 * services the interrupt.
 */
func (core *CPU) SetIntAckCallback(callback func(level IRQ) IntVector) {

}

/* The CPU will call the callback with whatever was in the data field of the
 * BKPT instruction for 68020+, or 0 for 68010.
 */
func (core *CPU) SetBkptAckCallback(callback func(data Long)) {}

/* The CPU calls this callback every time it encounters a RESET instruction.
 */
func (core *CPU) SetResetInstrCallback(callback func()) {}

/* The CPU calls this callback with the function code before every memory
 * access to set the CPU's function code according to what kind of memory
 * access it is (supervisor/user, program/data and such).
 */
func (core *CPU) SetFCCallback(callback func(newFC Long)) {}

/* Set a callback for the instruction cycle of the CPU.
 * You must enable M68K_INSTRUCTION_HOOK in m68kconf.h.
 * The CPU calls this callback just before fetching the opcode in the
 * instruction cycle.
 */
func (core *CPU) SetInstrHookCallback(callback func()) {}

/* Used by shift & rotate instructions */
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
)
