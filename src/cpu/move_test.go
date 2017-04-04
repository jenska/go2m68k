package cpu

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMoveq(t *testing.T) {
	cpu := NewM68k(NewMemoryHandler(1024*1024))
	var address uint32 = 0x1000

	for reg := 0; reg<8; reg++ {
		for imm := 0;imm<256; imm++ {
			opcode := 0x7000 + (reg<<9)  + imm
			cpu.write(Word, address, uint32(opcode) )
			address += Word.Size
		}
	}

	cpu.PC = 0x1000
	cpu.D[0] = 1
	cpu.Execute() // moveq #0, D0
	assert.Equal(t, uint32(0), cpu.D[0])
}