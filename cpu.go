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

type (
	instruction func(*M68K)

	// M68K CPU core
	M68K struct {
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
			oldSR := cpu.sr
			if !cpu.sr.S {
				cpu.sr.S = true
				cpu.usp = cpu.a[7]
			}
			cpu.a[7] = cpu.ssp
			cpu.push(Long, cpu.pc)
			cpu.push(Word, int32(oldSR.bits()))

			xaddr := cpu.read(int32(err)<<2, Long)
			if xaddr == 0 {
				if xaddr = cpu.read(int32(UnintializedInterrupt)<<2, Long); xaddr == 0 {
					log.Print("Interrupt vector not set for uninitialised interrupt vector")
					cpu.stopped = true
				}
			}
			cpu.pc = xaddr
		}
	}
}

// Run until halted
func (cpu *M68K) Run(signals <-chan Signal) {
	defer cpu.catchError()
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
	cpu.instructions[cpu.ir](cpu)
	cpu.icount++
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
