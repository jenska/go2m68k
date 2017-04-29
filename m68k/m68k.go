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

	// required for bus/address error stackframe
	statusCode struct {
		program, read, instruction bool
	}

	irqMode          int
	irqSamplingLevel uint8

	rmwCycle    bool
	stop        bool
	doubleFault bool
	CycleCount  int64

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
	cpu.Reset()
	return cpu
}

func (cpu *M68K) Reset() {
	cpu.statusCode.program = false
	cpu.statusCode.read = true
	cpu.statusCode.instruction = true
	cpu.irqSamplingLevel = 0
	cpu.rmwCycle = false
	cpu.stop = false
	cpu.doubleFault = false
	cpu.sync(16)
	cpu.SR.Set(0x2700)
	cpu.A[7] = cpu.Read(Long, 0)
	cpu.PC = cpu.Read(Long, 4)
	cpu.fullPrefetch()
}

func (cpu *M68K) sync(cycles int64) {
	cpu.CycleCount += cycles
}

func (cpu *M68K) Execute() {
	if cpu.doubleFault {
		cpu.sync(4)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(exception); ok {
				cpu.statusCode.program = false
				cpu.statusCode.read = true
				cpu.statusCode.instruction = true
			} else {
				log.Fatalf("unable to recover from unexpeced error %s", err)
			}
		}
	}()

	cpu.PC += 2
	cpu.instructions[cpu.IRD](cpu, cpu.IRD)
}

func (cpu *M68K) conditionalTest(code uint32) bool {
	switch code & 0xF {
	case 0:
		return true
	case 1:
		return false
	case 2:
		return !cpu.SR.C && !cpu.SR.Z
	case 3:
		return cpu.SR.C || cpu.SR.Z
	case 4:
		return !cpu.SR.C
	case 5:
		return cpu.SR.C
	case 6:
		return !cpu.SR.Z
	case 7:
		return cpu.SR.Z
	case 8:
		return !cpu.SR.V
	case 9:
		return cpu.SR.V
	case 10:
		return !cpu.SR.N
	case 11:
		return cpu.SR.N
	case 12:
		return !(cpu.SR.N != cpu.SR.V)
	case 13:
		return cpu.SR.N != cpu.SR.V
	case 14:
		return (cpu.SR.N && cpu.SR.V && !cpu.SR.Z) || (!cpu.SR.N && !cpu.SR.V && !cpu.SR.Z)
	default:
		return cpu.SR.Z || (cpu.SR.N && !cpu.SR.V) || (!cpu.SR.N && cpu.SR.V)
	}
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
