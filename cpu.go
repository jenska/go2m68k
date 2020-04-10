package cpu

import (
	"strconv"
)

//go:generate stringer -type=CPUError

// Exceptions handled by emulation
const (
	BusError                CPUError = 2
	AdressError                      = 3
	IllegalInstruction               = 4
	ZeroDivideError                  = 5
	PrivilegeViolationError          = 8
	UnintializedInterrupt            = 15

	// Address Bus Mask for 68000 CPU
	BusMask = 0x00ffffff
)

type (
	CPUError uint32

	AddressBus interface {
		Read(address uint32, s *Size) int
		Write(address uint32, s *Size, value int)
		Reset()
	}

	CPUBuilder interface {
		AttachBus(AddressBus) CPUBuilder
		AttachBaseArea(ssp, pc, size uint32) CPUBuilder
		InitISA68000() CPUBuilder
		Build() *M68K
		Go() *M68K
	}

	// M68K CPU
	M68K struct {
		IsRunning bool
		Reset     chan bool
		Halt      chan bool

		instructions [0x10000]func()
		bus          AddressBus
		trace        chan func()

		D [8]int32  // Data register
		A [8]uint32 // Address register (A[7] == USP)

		SSP, USP uint32 // Supervisoro Stackpointer#
		PC       uint32 // Program counter
		SR       SSR    // Supervisor Status Register (+ CCR)
		IR       uint16 // Instruction Register
	}
)

func (e CPUError) Error() string {
	return strconv.Itoa(int(e))
}

func NewCPU() CPUBuilder {
	return &M68K{SR: SSR{S: true}}
}

func (cpu *M68K) AttachBus(bus AddressBus) CPUBuilder {
	cpu.bus = bus
	return cpu
}

func (cpu *M68K) AttachBaseArea(ssp, pc, size uint32) CPUBuilder {
	cpu.bus = NewIOManager(size, NewBaseArea("BaseArea", cpu, ssp, pc, size))
	return cpu
}

// Go routine starter for M68K
func (cpu *M68K) Build() *M68K {
	if cpu.bus == nil {
		panic("no bus attached")
	}
	cpu.Reset = make(chan bool)
	cpu.Halt = make(chan bool)
	cpu.trace = make(chan func())
	return cpu
}

func (cpu *M68K) Go() *M68K {
	cpu = cpu.Build()
	go func() {
		cpu.IsRunning = true
		defer func() {
			close(cpu.Reset)
			close(cpu.Halt)
			close(cpu.trace)
			cpu.IsRunning = false
		}()

		for {
			select {
			case <-cpu.Reset:
				cpu.bus.Reset()
				cpu.SR.Bits(0x2700)
				cpu.SSP, cpu.PC = cpu.readAddress(0), cpu.readAddress(4)

			case <-cpu.Halt:
				return

			case fn := <-cpu.trace:
				fn()

			default:
				cpu.step()
			}
		}

	}()
	return cpu
}

func (cpu *M68K) step() *M68K {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(CPUError); ok {
				cpu.raiseException(err)
			}
		}
	}()

	cpu.IR = uint16(cpu.read(cpu.PC, Word))

	return cpu
}

func (cpu *M68K) raiseException(err CPUError) {
	oldSR := cpu.SR
	if !cpu.SR.S {
		cpu.SR.S = true
	}
	cpu.USP = cpu.A[7]
	cpu.A[7] = cpu.SSP
	cpu.push(Long, int(cpu.PC))
	cpu.push(Word, int(oldSR.ToBits()))

	if xaddr := cpu.readAddress(uint32(err) << 2); xaddr == 0 {
		if xaddr = cpu.readAddress(UnintializedInterrupt << 2); xaddr == 0 {
			cpu.Halt <- true
		}
	} else {
		cpu.PC = xaddr
	}
}

func (cpu *M68K) readAddress(a uint32) uint32 {
	return uint32(cpu.bus.Read(a&BusMask, Long))
}

func (cpu *M68K) read(a uint32, s *Size) int {
	if a&1 == 1 && s != Byte {
		panic(AdressError)
	}
	return cpu.bus.Read(a&BusMask, s)
}

func (cpu *M68K) write(a uint32, s *Size, value int) {
	if a&1 == 1 && s != Byte {
		panic(AdressError)
	}
	cpu.bus.Write(a&BusMask, s, value)
}

func (cpu *M68K) push(s *Size, value int) {
	cpu.A[7] -= s.size
	cpu.write(cpu.A[7], s, value)
}

func (cpu *M68K) pop(s *Size) int {
	result := cpu.read(cpu.A[7], s)
	cpu.A[7] += s.size
	return result
}
