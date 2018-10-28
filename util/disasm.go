package util

import (
	"fmt"
	"log"

	"github.com/jenska/atari2go/cpu"
)

type disassembledOperand struct {
	operand string
	size    int
	memory  int
}

type disassembledOpcode struct {
	address     cpu.Address
	opcode      int
	instruction string
	op1         *disassembledOperand
	op2         *disassembledOperand
}

type disassemble func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode

var opcodes []disassemble

func init() {
	log.Println("init disasm")

	unknown := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode {
		return &disassembledOpcode{address: address, opcode: readO(bus, address), instruction: "????"}
	}

	opcodes = make([]disassemble, 0x10000)
	for i, _ := range opcodes {
		opcodes[i] = unknown
	}

	moveq := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode {
		opcode := readO(bus, address)
		return &disassembledOpcode{
			address:     address,
			opcode:      opcode,
			instruction: "moveq",
			op1:         &disassembledOperand{operand: fmt.Sprintf("#$%02x", opcode&0xff)},
			op2:         &disassembledOperand{operand: fmt.Sprintf("d%d", ((opcode >> 9) & 0x07))},
		}
	}
	base := 0x7000
	for reg := 0; reg < 8; reg++ {
		for imm := 0; imm < 256; imm++ {
			opcodes[base+(reg<<9)+imm] = moveq
		}
	}
}

func readO(bus cpu.AddressBus, address cpu.Address) int {
	o, err := bus.Read(address, cpu.Word)
	if err != nil {
		panic("invalid address")
	}
	return o
}

func (opcode *disassembledOpcode) Size() int {
	size := 2
	if opcode.op1 != nil {
		size += opcode.op1.size
	}
	if opcode.op2 != nil {
		size += opcode.op2.size
	}
	return size
}

func (opcode *disassembledOpcode) String() string {
	result := fmt.Sprintf("%08x %04x", opcode.address, opcode.opcode)
	op1hex := ""
	op2hex := ""
	opstr := ""

	if opcode.op1 != nil {
		switch opcode.op1.size {
		case 2:
			op1hex = fmt.Sprintf("%04x", opcode.op1.memory)
		case 4:
			op1hex = fmt.Sprintf("%08x", opcode.op1.memory)
		}
		opstr = opcode.op1.operand
	}
	if opcode.op2 != nil {
		switch opcode.op2.size {
		case 2:
			op2hex = fmt.Sprintf("%04x", opcode.op2.memory)
		case 4:
			op2hex = fmt.Sprintf("%08x", opcode.op2.memory)
		}
		opstr += ", " + opcode.op2.operand
	}
	result += fmt.Sprintf("%8s %8s %s %s", op1hex, op2hex, opcode.instruction, opstr)

	return result
}
