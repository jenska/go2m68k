package cpu

// TODO: throw bus error on odd memory accesses

import (
	"log"
	"os"
)

var Trace = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
var Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
var Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
var Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

type Instruction func(cpu *M68k)

type M68k struct {
	A            [8]uint32
	D            [8]uint32
	SR           StatusRegister
	SSP, USP uint32
	memory       AddressHandler

	PC uint32
	instructions []Instruction
	opcode       uint16

	exception, halt      bool
}

func NewM68k(memory AddressHandler) *M68k {
	cpu := &M68k{}
	cpu.memory = memory
	cpu.SR = NewStatusRegister(cpu)
	cpu.instructions = []Instruction{move}
	cpu.Reset()

	return cpu
}

func (cpu *M68k) Reset() {
	// TODO reset
}

func (cpu *M68k) Execute() {
/*	defer func() {
		//r := recover()

	}n*/
	if !cpu.halt {
		// TODO execute
	}
}
