package cpu

// Exception vectors handled by emulation
const (
	BusErrorVector              = 8
	AdressErrorVector           = 12
	IllegalInstructionVector    = 16
	ZeroDivideVector            = 20
	PrivilegeViolationError     = 32
	UnintializedInterruptVector = 60

	BusMask = 0x00ffffff
)

type (
	// Pin for CPU
	Pin       chan bool
	Exception chan uint32

	AddressBus interface {
		AttachCPU(cpu *M68000)
		Read(address uint32, s *Size) int
		Write(address uint32, s *Size, value int)
		Reset()
	}

	// M68000 CPU
	M68000 struct {
		Reset Pin
		Halt  Pin
		Error Exception

		io           AddressBus
		instructions []func()
		clock        int

		D [8]int32  // Data register
		A [8]uint32 // Address register (A[7] == USP)

		SSP uint32 // Supervisoro Stackpointer#
		PC  uint32 // Program counter
		SR  SSR    // Supervisor Status Register (+ CCR)
		IR  uint16 // Instruction Register
	}
)

func NewCPU() *M68000 {
	cpu := &M68000{
		Reset: make(Pin),
		Halt:  make(Pin),
		Error: make(Exception),
	}
	return cpu
}

func (cpu *M68000) AttachBus(bus AddressBus) {
	cpu.io = bus
	bus.AttachCPU(cpu)

	go func() {
		stop := false
		for !stop {
			select {
			case <-cpu.Reset:
				cpu.SR.Bits(0x2700)
				cpu.io.Reset()
				cpu.A[7], cpu.PC = cpu.readA(0), cpu.readA(4)

			case <-cpu.Halt:
				stop = true

			case x := <-cpu.Error:
				cpu.raiseException(x)

			default:
				cpu.IR = uint16(cpu.pop(Word))
			}
		}
	}()

	cpu.Reset <- true
}

func (cpu *M68000) raiseException(vector uint32) {
	oldSR := cpu.SR
	if !cpu.SR.S {
		cpu.SR.S = true
		cpu.SSP, cpu.A[7] = cpu.A[7], cpu.SSP
	}
	cpu.push(Long, int(cpu.PC))
	cpu.push(Word, int(oldSR.ToBits()))

	if xaddr := cpu.readA(vector); xaddr == 0 {
		if xaddr = cpu.readA(UnintializedInterruptVector); xaddr == 0 {
			cpu.Halt <- true
		}
	} else {
		cpu.PC = xaddr
	}
	cpu.clock += 34
}

func (cpu *M68000) readA(a uint32) uint32 {
	return uint32(cpu.io.Read(a&BusMask, Long))
}

func (cpu *M68000) read(a uint32, s *Size) int {
	if a&1 == 1 && s != Byte {
		cpu.Error <- AdressErrorVector
	}
	return cpu.io.Read(a&BusMask, s)
}

func (cpu *M68000) write(a uint32, s *Size, value int) {
	if a&1 == 1 && s != Byte {
		cpu.Error <- AdressErrorVector
	}
	cpu.io.Write(a&BusMask, s, value)
}

func (cpu *M68000) push(s *Size, value int) {
	cpu.A[7] -= s.size
	cpu.write(cpu.A[7], s, value)
}

func (cpu *M68000) pop(s *Size) int {
	result := cpu.read(cpu.A[7], s)
	cpu.A[7] += s.size
	return result
}
