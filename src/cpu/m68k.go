package cpu

import (
	"log"
	"mem"
	"fmt"
)

type Address uint32

type Instruction interface {
	 execute() uint32
}

type Disassembler interface {
	Disassemble(address Address) DisassembledInstruction
}

type DisassembledEA struct {
	operand Operand
	size int
	memory Address
	ea string
}

type DisassembledInstruction struct {
	address Address
	opcode uint16
	numOperands int
	instruction string
	op1, op2 DisassembledEA
}

type M68k struct {
	A [8]Address
	D [8]int32
	SR StatusRegister
	PC, SP Address

	instructions [0x10000]Instruction
	disassembler [0x10000]Disassembler
	loadedInstructions int

	log log.Logger
}

func NewM68k(memory mem.PhysicalAddressSpace) M68k {
	cpu := M68k{}
	cpu.SR = NewStatusRegister( &cpu )

	return cpu
}

func (cpu *M68k) SetSupervisorMode(mode bool) {
}

func (cpu *M68k) SetInstruction(opcode uint16, i Instruction, d Disassembler) {
	if current := cpu.instructions[opcode]; current != nil {
		cpu.instructions[opcode] = i
		cpu.disassembler[opcode] = d
		cpu.loadedInstructions++
	} else {
		panic(fmt.Sprintf("Attempted to overwrite existing instruction [%s] at 0x%04x with [%s]",
			cpu.disassembler[opcode],opcode,d))
	}
}