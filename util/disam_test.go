package util

import (
	"fmt"
	"testing"

	"github.com/jenska/atari2go/cpu"
	"github.com/jenska/atari2go/mem"
	"github.com/stretchr/testify/assert"
)

var bus = mem.NewAddressBus(
	mem.NewRAM(1014, 1024),
)

func TestDisasmMoveq(t *testing.T) {
	bus.Write(1024, cpu.Word, 0x7000)
	bus.Write(1026, cpu.Word, 0x7001)
	bus.Write(1028, cpu.Word, 0x7002)
	bus.Write(1030, cpu.Word, 0x7201)
	bus.Write(1032, cpu.Word, 0x72FF)

	start := cpu.Address(1024)
	end := start + cpu.Address(10)
	for start < end {
		opcode, err := bus.Read(start, cpu.Word)
		if err != nil {
			panic(fmt.Sprintf("invalid read %s", err))
		}

		instruction := opcodes[opcode>>6](bus, start)
		assert.Equal(t, "moveq", instruction.instruction)
		assert.Equal(t, 2, instruction.Size())
		start += cpu.Address(instruction.Size())
	}
}
