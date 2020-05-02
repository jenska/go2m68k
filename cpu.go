package cpu

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
		icount       int // overall instructions performend

		bus   AddressBus
		read  Reader
		write Writer

		d [8]int32 // Data register
		a [8]int32 // Address register (A[7] == USP)

		pc int32  // Program counter
		ir uint16 // Instruction Register

		ssp, usp int32 // Supervisoro Stackpointer#
		sr       ssr   // Supervisor Status Register (+ CCR)

		stopped bool // stopped state
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
			cpu.raiseException(err)
		} else {
			log.Print(err)
			cpu.stopped = true
		}
	}
}

func (cpu *M68K) step() *M68K {
	cpu.ir = uint16(cpu.popPC(Word))
	cpu.instructions[cpu.ir](cpu)
	cpu.icount++
	return cpu
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
			cpu.step()
		}
	}
}

// Step through a single instruction
func (cpu *M68K) Step() *M68K {
	defer cpu.catchError()
	return cpu.step()
}

// Reset sets the cpu back to initial state
func (cpu *M68K) Reset() {
	cpu.bus.reset()
	cpu.sr.setbits(0x2700)
	cpu.ssp, cpu.pc = cpu.readAddress(0), cpu.readAddress(4)
	cpu.a[7] = cpu.ssp
	cpu.stopped = false
	cpu.icount = 0
}

func (cpu *M68K) raiseException(err Error) {
	oldSR := cpu.sr
	if !cpu.sr.S {
		cpu.sr.S = true
		cpu.usp = cpu.a[7]
	}
	cpu.a[7] = cpu.ssp
	cpu.push(Long, cpu.pc)
	cpu.push(Word, int32(oldSR.bits()))

	xaddr := cpu.readAddress(int32(err) << 2)
	if xaddr == 0 {
		if xaddr = cpu.readAddress(int32(UnintializedInterrupt) << 2); xaddr == 0 {
			panic("Interrupt vector not set for uninitialised interrupt vector")
		}
	}
	cpu.pc = xaddr
}
