package cpu

/*
goos: linux
goarch: amd64
pkg: github.com/jenska/go2m68k
BenchmarkDbra-4   	    2056	    526254 ns/op
PASS

Showing nodes accounting for 1050ms, 90.52% of 1160ms total
Showing top 10 nodes out of 25
      flat  flat%   sum%        cum   cum%
     180ms 15.52% 15.52%      710ms 61.21%  github.com/jenska/go2m68k.(*M68K).SetISA68000.func1
     150ms 12.93% 28.45%      530ms 45.69%  github.com/jenska/go2m68k.(*addressAreaQueue).read
     140ms 12.07% 40.52%      280ms 24.14%  github.com/jenska/go2m68k.NewBaseArea.func1
     110ms  9.48% 50.00%      960ms 82.76%  github.com/jenska/go2m68k.(*M68K).step
     100ms  8.62% 58.62%      100ms  8.62%  github.com/jenska/go2m68k.(*addressAreaQueue).findArea
      90ms  7.76% 66.38%      800ms 68.97%  github.com/jenska/go2m68k.(*M68K).popPC (inline)
      90ms  7.76% 74.14%      140ms 12.07%  github.com/jenska/go2m68k.glob..func4
      90ms  7.76% 81.90%       90ms  7.76%  runtime.chanrecv
      50ms  4.31% 86.21%       50ms  4.31%  encoding/binary.bigEndian.Uint16
      50ms  4.31% 90.52%      340ms 29.31%  github.com/jenska/go2m68k.dbra
*/
import (
	"fmt"
	"log"
)

// TODO:
//  add cpu cycles
//  add tracing (t0,t1)
// Exceptions handled by emulation
const (
	M68K_CPU_TYPE_68000 Type = iota
	M68K_CPU_TYPE_68010
	M68K_CPU_TYPE_68EC020
	M68K_CPU_TYPE_68020
	M68K_CPU_TYPE_68EC030
	M68K_CPU_TYPE_68030
	M68K_CPU_TYPE_68EC040
	M68K_CPU_TYPE_68LC040
	M68K_CPU_TYPE_68040
	M68K_CPU_TYPE_SCC68070

	BusError                = 2
	AdressError             = 3
	IllegalInstruction      = 4
	ZeroDivideError         = 5
	PrivilegeViolationError = 8
	UnintializedInterrupt   = 15
	TrapBase                = 32

	HaltSignal Signal = iota
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
	// Error type for CPU Errors
	Error struct {
		index   int32
		name    string
		address *int32
		c       *M68K
		ir      *uint16
		x       *int32
	}

	// Type of CPU
	Type int32
	// Signal external CPU events
	Signal int32

	// Reader accessor for read accesses
	Reader func(int32, *Size) int32
	// Writer accessor for write accesses
	Writer func(int32, *Size, int32)
	// Reset prototype
	Reset func()

	// AddressArea container for address space area
	AddressArea struct {
		read  Reader
		write Writer
		raw   []byte
		reset Reset
	}

	// AddressBus for accessing address areas
	AddressBus interface {
		read(address int32, s *Size) int32
		write(address int32, s *Size, value int32)
		reset()
	}

	instruction func(*M68K)

	// M68K CPU core
	M68K struct {
		cpuType      Type
		instructions [0x10000]instruction
		icount       int // overall instructions performed

		cycles      [0x10000]int
		eaIdxCycles [64]int
		iclocks     int // number of clocks remaining

		bus   AddressBus
		read  Reader
		write Writer
		eaSrc []ea
		eaDst []ea

		da  [16]int32 // data and address registers
		d   []int32   // slice of data registers
		a   []int32   // slice of address registers
		pc  int32     // Program counter
		ppc int32     // previous program counter
		ir  uint16    // Instruction Register

		ssp, usp int32 // Supervisoro Stackpointer#
		sr       ssr   // Supervisor Status Register (+ CCR)

		stopped bool // stopped state
		hasFPU  bool
	}
)

func (cpu *M68K) String() string {
	result := fmt.Sprintf("SR %s PC %08x USP %08x SSP %08x\n", cpu.sr, cpu.pc, cpu.usp, cpu.ssp)
	for i := range cpu.d {
		result += fmt.Sprintf("D%d %08x ", i, uint32(cpu.d[i]))
	}
	result += "\n"
	for i := range cpu.a {
		result += fmt.Sprintf("A%d %08x ", i, uint32(cpu.a[i]))
	}
	result += "\n"

	return result
}

func (cpu *M68K) catchError() {
	if r := recover(); r != nil {
		if err, ok := r.(Error); ok {
			log.Printf("cpu.catchError: %s\n", err)
			oldSR := cpu.sr
			if !cpu.sr.S {
				cpu.sr.S = true
				cpu.usp = cpu.a[7]
			}
			cpu.a[7] = cpu.ssp
			cpu.push(Long, cpu.pc)
			cpu.push(Word, oldSR.bits())

			exception := err.index
			if err.x != nil {
				exception += *err.x
			}
			if xaddr := cpu.read(exception<<2, Long); xaddr == 0 {
				if xaddr = cpu.read(UnintializedInterrupt<<2, Long); xaddr == 0 {
					panic(fmt.Sprintf("Interrupt vector not set for uninitialised interrupt vector from 0x%08x", cpu.pc))
					// cpu.stopped = true
				}
			} else {
				cpu.pc = xaddr
			}
		} else {
			panic(r)
		}
	}
}

// Run until halted
// TODO: prefetch in goroutine?
func (cpu *M68K) Run(signals <-chan Signal) {
	cpu.stopped = false
	for !cpu.stopped {
		select {
		case signal := <-signals:
			if signal == ResetSignal {
				cpu.Reset()
			} else if signal == HaltSignal {
				cpu.stopped = true
				break
			}
		default:
			cpu.Step()
		}
	}
}

// Step through a single instruction
func (cpu *M68K) Step() *M68K {
	defer cpu.catchError()
	cpu.ir = uint16(cpu.popPc(Word))
	if instruction := cpu.instructions[cpu.ir]; instruction != nil {
		instruction(cpu)
		cpu.icount++
	} else {
		// debug.PrintStack()
		panic(NewError(IllegalInstruction, cpu, cpu.pc, nil))
	}
	// TODO: trace
	return cpu
}

// Reset sets the cpu back to initial state
func (cpu *M68K) Reset() {
	cpu.bus.reset()
	cpu.sr.setbits(0x2700)
	cpu.ssp, cpu.pc = cpu.read(0, Long), cpu.read(4, Long)
	cpu.a[7] = cpu.ssp
	cpu.stopped = false
	cpu.icount = 0
}

var errorString = map[int32]string{
	BusError:                "bus error",
	IllegalInstruction:      "illegal instruction",
	PrivilegeViolationError: "privilege violation",
	TrapBase:                "trap #%d",
}

// NewError creates a new CPU Error object
func NewError(index int32, c *M68K, address int32, extra *int32) Error {
	res := Error{index: index, x: extra, c: c}
	pc := address
	res.address = &pc
	if c != nil {
		ir := c.ir
		res.ir = &ir
	}
	return res
}

func (e Error) Error() string {
	dsc := errorString[e.index]
	if e.x != nil && dsc != "" {
		dsc = fmt.Sprintf(dsc, *e.x)
	} else if dsc == "" {
		dsc = fmt.Sprintf("error %d", e.index)
	}

	if e.ir != nil {
		dasm, _ := Disassemble(e.c.cpuType, *e.address, e.c.bus)
		return fmt.Sprintf("cpu %s (PC %08x, IR %04x)\n%s", dsc, *e.address, *e.ir, dasm)
	}
	return fmt.Sprintf("cpu %s at %08x", dsc, *e.address)
}
