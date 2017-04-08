package m68k

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOri2ccr(t *testing.T) {
	sr := StatusRegister{}
	sr.C = true
	sr.N = true

	cpu := NewM68k(NewMemoryHandler(0x10000))

	cpu.SR.C = false
	cpu.SR.N = false

	cpu.PC = uint32(0x01000)
	opcode := uint32(0x0074)
	cpu.Write(Word, cpu.PC, opcode)
	cpu.Write(Word, cpu.PC+Word.Size, uint32(sr.GetCCR()))
	cpu.Execute()

	assert.Equal(t, true, cpu.SR.C)
	assert.Equal(t, true, cpu.SR.N)
}

func TestOri2sr(t *testing.T) {
	sr := StatusRegister{}
	sr.C = true
	sr.N = true

	cpu := NewM68k(NewMemoryHandler(0x10000))
	cpu.SR.C = false
	cpu.SR.N = false

	cpu.PC = uint32(0x01000)
	opcode := uint32(0x00f4)
	cpu.Write(Word, cpu.PC, opcode)
	cpu.Write(Word, cpu.PC+Word.Size, uint32(sr.Get()))
	cpu.Execute()

	assert.Equal(t, true, cpu.SR.C)
	assert.Equal(t, true, cpu.SR.N)

	cpu.SR.SetS(false)

	// privileged exception handler
	cpu.Write(Long, XptPrv*Long.Size, 0x02000)
	cpu.Write(Word, 0x02000, 0x7001) // moveq #1, D0

	cpu.PC = uint32(0x01000)
	assert.Equal(t, 34, cpu.Execute()) // bang
	assert.Equal(t, uint32(0x02000), cpu.PC)
	assert.Equal(t, true, cpu.SR.S())

	cpu.Execute() // moveq
	assert.Equal(t, uint32(1), cpu.D[0])
	assert.Equal(t, uint32(0x02002), cpu.PC)
}
