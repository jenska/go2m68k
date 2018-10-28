package cpu

import (
	"log"
)

func set(instructions []Instruction, index int, f Instruction) {
	if instructions[index] != nil {
		panic("opcode already in use")
	}
	instructions[index] = f
}

func registerM68KInstructions(c *M68K) {
	log.Println("register M68K instruction set")

	log.Println("\tmoveq")
	moveq := func(opcode int) int {
		c.D[(opcode>>9)&0x7] = Data(int8(opcode & 0xff))
		return 4
	}
	base := 0x7000
	for reg := 0; reg < 8; reg++ {
		for imm := 0; imm < 256; imm++ {
			set(c.instructions, base+(reg<<9)+imm, moveq)
		}
	}

	log.Println("\tunknown")
	unknown := func(opcode int) int {
		if opcode&0xf000 == 0xf000 {
			c.RaiseException(LineF)
		} else if opcode&0xa000 == 0xa000 {
			c.RaiseException(LineA)
		} else {
			c.RaiseException(IllegalInstruction)
		}
		return 34
	}
	for i := range c.instructions {
		if c.instructions[i] == nil {
			c.instructions[i] = unknown
		}
	}

}
