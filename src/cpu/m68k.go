package cpu

import (
	"fmt"
	"github.com/golang/glog"
)

type Instruction func() int

type M68k struct {
	A        [8]uint32
	D        [8]uint32
	SR       StatusRegister
	SSP, USP uint32
	PC       uint32

	memory AddressHandler

	ira          uint32
	instructions [0x10000]Instruction
}

func NewM68k(memory AddressHandler) *M68k {
	cpu := &M68k{}
	cpu.memory = memory
	cpu.SR = NewStatusRegister(cpu)
	cpu.SR.Set(0x2700)
	cpu.A[7] = cpu.read(Long,XPT_SPR<<2)
	cpu.PC = cpu.read(Long,XPT_PCR<<2)
	return cpu
}

func (cpu *M68k) Reset() {
	// TODO reset
}

func (cpu *M68k) Execute() int {
	cpu.ira = cpu.PC
	opcode := cpu.popPC(Word)

	if instruction := cpu.instructions[opcode]; instruction != nil {
		return instruction()
	} else {
		if (opcode & 0xa000) == 0xa000 {
			return cpu.RaiseException(XPT_LNA)
		} else if (opcode & 0xf000) == 0xf000 {
			return cpu.RaiseException(XPT_LNF)
		}
		glog.Errorf("Illegal instrcution $%04x ar $%x\n", opcode, cpu.ira)
		return cpu.RaiseException(XPT_ILL)
	}
}

func (cpu *M68k) RaiseException(vector uint16) int {
	address := uint32(vector) << 2
	cpu.pushSP(Long, cpu.PC)
	cpu.pushSP(Word, uint32(cpu.SR.Get()))

	if !cpu.SR.S() {
		cpu.SR.SetS(true)
	}
	cpu.SR.T = false

	if address = cpu.read(Long, address); address == 0 {
		if address = cpu.read(Long, UNINITIALIZED_INTERRUPT_VECTOR<<2); address == 0 {
			panic(fmt.Sprintf("Interrupt vector not set for uninitialised interrupt vector while trapping uninitialised vector %d\n", vector))
		}
	}

	cpu.PC = address
	return 34
}

// TODO Disassemble current instruction
func (cpu *M68k) String() string {
	d := cpu.D
	a := cpu.A
	format := "D0: %08x   D4: %08x   A0: %08x   A4: %08x     PC:  %08x\n" +
		"D1: %08x   D5: %08x   A1: %08x   A5: %08x     SR:  %04x %s\n" +
		"D2: %08x   D6: %08x   A2: %08x   A6: %08x     USP: %08x\n" +
		"D3: %08x   D7: %08x   A3: %08x   A7: %08x     SSP: %08x\n\n"

	return fmt.Sprintf(format, d[0], d[4], a[0], a[4], cpu.PC,
		d[1], d[5], a[1], a[5], cpu.SR.Get(), cpu.SR,
		d[2], d[6], a[2], a[6], cpu.USP,
		d[3], d[7], a[3], a[7], cpu.SSP)
}

func (cpu *M68k) read(o *Operand, address uint32) uint32 {
	address &= 0x00ffffff
	if v, ok := cpu.memory.Mem(o, address); ok {
		return v
	}
	// TODO raise exception
	return 0
}

func (cpu *M68k) write(o *Operand, address uint32, value uint32) {
	address &= 0x00ffffff
	if cpu.memory.setMem(o, address, value) {
		return
	}
	// TODO raise exception
}

func (cpu *M68k) popPC(o *Operand) uint32 {
	result := cpu.read(o, cpu.PC)
	cpu.PC += o.Size
	return result
}

func (cpu *M68k) pushPC(o *Operand, v uint32) {
	cpu.PC -= o.Size
	cpu.write(o, cpu.PC, v)
}

func (cpu *M68k) popSP(o *Operand) uint32 {
	result := cpu.read(o, cpu.A[7])
	cpu.A[7] += o.Size
	return result
}

func (cpu *M68k) pushSP(o *Operand, v uint32) {
	cpu.A[7] -= o.Size
	cpu.write(o, cpu.A[7], v)
}
