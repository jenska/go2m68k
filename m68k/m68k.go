package m68k

import (
	"fmt"
	"log"

	"github.com/jenska/atari2go/mem"
)

const (
	UserVectorInterrupt = iota
	AutoVectorInterrut
	Uninitialized
)

type M68K struct {
	SR           StatusRegister
	A, D         [8]uint32
	SSP, USP, PC uint32
	IRC, IRD, IR uint16

	CycleCount int64

	irqMode          int
	irqSamplingLevel uint8

	doubleFault bool
	cpuHalted   func(cpu *M68K)

	memory       mem.AddressHandler
	instructions []instruction
	eaTable      []ea
}

// New M68k CPU instance
func NewM68k(memory mem.AddressHandler) *M68K {
	cpu := &M68K{}
	cpu.memory = memory
	cpu.init68000InstructionSet()
	cpu.initEATable()
	cpu.SR = newStatusRegister(cpu)
	cpu.cpuHalted = func(cpu *M68K) {
		log.Fatal("Interrupt vector not set for uninitialised interrupt vector while trapping uninitialised vector")
	}
	cpu.Reset()
	return cpu
}

func (cpu *M68K) Reset() {
	cpu.doubleFault = false
	cpu.irqSamplingLevel = 0
	cpu.sync(16)
	cpu.SR.Set(0x2700)
	cpu.A[7] = cpu.Read(Long, 0)
	cpu.PC = cpu.Read(Long, 4)
	cpu.fullPrefetch()
}

func (cpu *M68K) sync(cycles uint32) {
	cpu.CycleCount += int64(cycles)
}

func (cpu *M68K) String() string {
	d := cpu.D
	a := cpu.A
	format := "D0: %08x   D4: %08x   A0: %08x   A4: %08x     PC:  %08x\n" +
		"D1: %08x   D5: %08x   A1: %08x   A5: %08x     SR:  %04x\n" +
		"D2: %08x   D6: %08x   A2: %08x   A6: %08x     USP: %08x\n" +
		"D3: %08x   D7: %08x   A3: %08x   A7: %08x     SSP: %08x\n\n" +
		"%s\n"

	return fmt.Sprintf(format, d[0], d[4], a[0], a[4], cpu.PC,
		d[1], d[5], a[1], a[5], cpu.SR.Get(),
		d[2], d[6], a[2], a[6], cpu.USP,
		d[3], d[7], a[3], a[7], cpu.SSP,
		disassembler(cpu.memory, cpu.PC).detailedFormat())
}
