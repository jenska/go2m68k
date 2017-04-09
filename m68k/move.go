package m68k

import (
	"fmt"
)

var shortExecutionTime = [][]int{{4, 4, 8, 8, 8, 12, 14, 12, 16},
	{4, 4, 8, 8, 8, 12, 14, 12, 16}, {8, 8, 12, 12, 12, 16, 18, 16, 20},
	{8, 8, 12, 12, 12, 16, 18, 16, 20}, {10, 10, 14, 14, 14, 18, 20, 18, 22},
	{12, 12, 16, 16, 16, 20, 22, 20, 24}, {14, 14, 18, 18, 18, 22, 24, 22, 26},
	{12, 12, 16, 16, 16, 20, 22, 20, 24}, {16, 16, 20, 20, 20, 24, 26, 24, 28},
	{12, 12, 16, 16, 16, 20, 22, 20, 24}, {14, 14, 18, 18, 18, 22, 24, 22, 26},
	{8, 8, 12, 12, 12, 16, 18, 16, 20}}

var longExecutionTime = [][]int{{4, 4, 12, 12, 12, 16, 18, 16, 20},
	{4, 4, 12, 12, 12, 16, 18, 16, 20}, {12, 12, 20, 20, 20, 24, 26, 24, 28},
	{12, 12, 20, 20, 20, 24, 26, 24, 28}, {14, 14, 22, 22, 22, 26, 28, 26, 30},
	{16, 16, 24, 24, 24, 28, 30, 28, 32}, {18, 18, 26, 26, 26, 30, 32, 30, 34},
	{16, 16, 24, 24, 24, 28, 30, 28, 32}, {20, 20, 28, 28, 28, 32, 34, 32, 36},
	{16, 16, 24, 24, 24, 28, 30, 28, 32}, {18, 18, 26, 26, 26, 30, 32, 30, 34},
	{12, 12, 20, 20, 20, 24, 26, 24, 28}}

var m2RTiming = []int{0, 0, 12, 12, 0, 16, 18, 16, 20, 16, 18}
var r2MTiming = []int{0, 0, 8, 0, 8, 12, 14, 12, 16}

func registerMoveInstructions(cpu *M68K) {

	// moveq
	for reg := 0; reg < 8; reg++ {
		for v := -128; v < 128; v++ {
			opcode := 0x7000 + (reg << 9) + int(uint8(v))
			value := v
			target := &cpu.D[reg]
			n := value < 0
			z := value == 0
			cpu.registerInstruction(opcode, func(cpu *M68K) int {
				return moveq(uint32(value), target, &cpu.SR, n, z)
			})
		}
	}
}

func moveq(value uint32, target *uint32, sr *StatusRegister, n, z bool) int {
	*target = value
	sr.N, sr.Z = n, z
	sr.C, sr.V = false, false
	return 4
}

func dMoveq(handler AddressHandler, address uint32) *disassembledInstruction {
	opcode := dOpcode(handler, address)
	dEA := []disassembledEA{
		disassembledEA{fmt.Sprintf("#%02d", int8(opcode&0xff)), nil, 0},
		disassembledEA{fmt.Sprintf("d%d", (opcode>>9)&7), nil, 0}}
	return &disassembledInstruction{"moveq", opcode, address, dEA}
}
