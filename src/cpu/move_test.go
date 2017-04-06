package cpu

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var cpu = NewM68k(NewMemoryHandler(1024 * 1024))

func BenchmarkMoveq(b *testing.B) {
	b.StopTimer()
	var address uint32 = 0x1000
	for reg := 0; reg < 8; reg++ {
		for imm := 0; imm < 256; imm++ {
			opcode := 0x7000 + (reg << 9) + imm
			cpu.write(Word, address, uint32(opcode))
			address += Word.Size
		}
	}
	b.StartTimer()

	n := 0
	for {
		cpu.PC = 0x1000
		for reg := 0; reg < 8; reg++ {
			for imm := 0; imm < 256; imm++ {
				cpu.Execute() // moveq #imm, Dreg
				if n++; n >= b.N {
					return
				}
			}
		}
	}
}

func TestMoveq(t *testing.T) {
	cpu := NewM68k(NewMemoryHandler(1024 * 1024))
	var address uint32 = 0x1000

	for reg := 0; reg < 8; reg++ {
		for imm := 0; imm < 256; imm++ {
			opcode := 0x7000 + (reg << 9) + imm
			cpu.write(Word, address, uint32(opcode))
			address += Word.Size
		}
	}
	cpu.PC = 0x1000
	for reg := 0; reg < 8; reg++ {
		cpu.D[reg] = 1
		for imm := 0; imm < 256; imm++ {
			cpu.Execute() // moveq #0, D0
			assert.Equal(t, int32(int8(imm)), int32(cpu.D[reg]), fmt.Sprintf("reg = D%02d, imm = #%d", reg, imm))
		}
	}
}
