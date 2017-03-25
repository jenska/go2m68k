package cpu

// TODO: throw bus error on odd memory accesses

import (
	"log"
	"os"
	"util"
)

var Trace = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
var Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
var Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
var Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

type Instruction func(cpu *M68k)

type M68k struct {
	A            [8]Address
	D            [8]int32
	SR           StatusRegister
	PC, SSP, USP Address
	memory       AddressHandler

	instructions []Instruction
	opcode       uint16

	delay, clkcnt   uint64
	lastPC          *util.Stack
	exception, halt bool
}

func NewM68k(ramSize int, memory AddressHandler) *M68k {
	cpu := &M68k{}
	cpu.memory = memory
	cpu.SR = NewStatusRegister(cpu)

	cpu.lastPC = util.NewStack()
	cpu.instructions = []Instruction{move}

	return cpu
}

func (cpu *M68k) Mem8(a Address) uint8 {
	a &= 0xffffff
	if v, ok := cpu.memory.Mem8(a); ok {
		return v
	}
	// TODO bus error
	return 0
}

func (cpu *M68k) setMem8(a Address, v uint8) {
	a &= 0xffffff
	if !cpu.memory.setMem8(a, v) {
		// TODO bus error
	}
}

func (cpu *M68k) Mem16(a Address) uint16 {
	a &= 0xffffff
	if v, ok := cpu.memory.Mem16(a); ok {
		return v
	}
	// TODO bus error
	return 0
}

func (cpu *M68k) setMem16(a Address, v uint16) {
	a &= 0xffffff
	if !cpu.memory.setMem16(a, v) {
		// TODO bus error
	}
}

func (cpu *M68k) Mem32(a Address) uint32 {
	a &= 0xffffff
	if v, ok := cpu.memory.Mem32(a); ok {
		return v
	}
	// TODO bus error
	return 0
}

func (cpu *M68k) setMem32(a Address, v uint32) {
	a &= 0xffffff
	if !cpu.memory.setMem32(a, v) {
		// TODO bus error
	}
}

func (cpu *M68k) execute() {
	if !cpu.halt {
		cpu.lastPC.Push(int32(cpu.PC))
		// TODO execute
	}
}

func (cpu *M68k) clock(n uint64) {
	for n > cpu.delay {
		n -= cpu.delay
		cpu.clkcnt = cpu.delay
		cpu.delay = 0
		cpu.execute()
		if cpu.delay == 0 {
			Warning.Printf("warning: delay == 0 at %08lx\n", cpu.PC)
		}
	}
	cpu.clkcnt += n
	cpu.delay -= n
}
