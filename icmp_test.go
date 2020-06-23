package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTst(t *testing.T) {
	tcpu.pc = 0x4000
	twrite(0x7000, 0x7200, 0x7400, 0x7600) // moveq #0, d0, d1, d2, d3
	twrite(0x4A00)                         // tst d0
	twrite(0x6000 + (0b0110 << 8) + 2)     // bne.s +2
	twrite(0x7201)                         // moveq #1, d1
	twrite(0x4A01)                         // tst d1
	twrite(0x6000 + (0b1101 << 8) + 2)     // blt.s +2
	twrite(0x74FF)                         // moveq #-1, d2
	twrite(0x6000 + (0b1100 << 8) + 2)     // bge.s +2
	twrite(0x76FF)                         // moveq #-1, d3
	trun(0x4000)

	assert.Equal(t, int32(0), tcpu.d[0])
	assert.Equal(t, int32(1), tcpu.d[1])
	assert.Equal(t, int32(-1), tcpu.d[2])
	assert.Equal(t, int32(-1), tcpu.d[3])
}
