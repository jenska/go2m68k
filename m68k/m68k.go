package m68k

import (
	"fmt"

	"github.com/golang/glog"
)

type M68K struct {
	A        [8]uint32
	D        [8]uint32
	SR       StatusRegister
	SSP, USP uint32
	PC       uint32

	memory       AddressHandler
	instructions []instruction
}

// New M68k CPU instance
func NewM68k(memory AddressHandler) *M68K {
	cpu := &M68K{}
	cpu.memory = memory
	cpu.init68000InstructionSet()
	cpu.SR = newStatusRegister(cpu)
	cpu.SR.Set(0x2700)
	cpu.A[7] = cpu.Read(Long, XptSpr<<2)
	cpu.PC = cpu.Read(Long, XptPcr<<2)
	return cpu
}

func (cpu *M68K) Reset() {
	// TODO reset
}

func (cpu *M68K) Execute() int {
	ira := cpu.PC
	opcode := uint16(cpu.popPC(Word))
	if instruction := cpu.instructions[opcode]; instruction != nil {
		return instruction(cpu)
	} else {
		glog.Errorf("Illegal instruction #$%04x at $%08x\n", opcode, ira)
		return cpu.RaiseException(XptIll)
	}
}

func (cpu *M68K) RaiseException(vector uint16) int {
	address := uint32(vector) << 2
	cpu.pushSP(Long, cpu.PC)
	cpu.pushSP(Word, uint32(cpu.SR.Get()))

	if !cpu.SR.S() {
		cpu.SR.SetS(true)
	}
	cpu.SR.T = false

	if address = cpu.Read(Long, address); address == 0 {
		if address = cpu.Read(Long, UNINITIALIZED_INTERRUPT_VECTOR<<2); address == 0 {
			panic(fmt.Errorf("Interrupt vector not set for uninitialised interrupt vector while trapping uninitialised vector %d\n", vector))
		}
	}

	cpu.PC = address
	return 34
}

// TODO Disassemble current instruction
func (cpu *M68K) String() string {
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

func (cpu *M68K) Read(o *operand, address uint32) uint32 {
	address &= 0x00ffffff
	if v, err := cpu.memory.Read(o, address); err == nil {
		return v
	} else {
		// TODO raise exception
		return 0
	}
}

func (cpu *M68K) Write(o *operand, address uint32, value uint32) {
	address &= 0x00ffffff
	if err := cpu.memory.Write(o, address, value); err == nil {
		return
	} else {
		// TODO raise exception
	}
}

func (cpu *M68K) popPC(o *operand) uint32 {
	result := cpu.Read(o, cpu.PC)
	cpu.PC += o.Size
	return result
}

func (cpu *M68K) pushPC(o *operand, v uint32) {
	cpu.PC -= o.Size
	cpu.Write(o, cpu.PC, v)
}

func (cpu *M68K) popSP(o *operand) uint32 {
	result := cpu.Read(o, cpu.A[7])
	cpu.A[7] += o.Size
	return result
}

func (cpu *M68K) pushSP(o *operand, v uint32) {
	cpu.A[7] -= o.Size
	cpu.Write(o, cpu.A[7], v)
}
