package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoveq(t *testing.T) {
	tcpu.Reset()

	tcpu.write(0x4000, Word, 0x7000) // moveq #0, d0
	tcpu.write(0x4002, Word, 0x7001) // moveq #1, d0
	tcpu.write(0x4004, Word, 0x7200) // moveq #0, d1
	tcpu.write(0x4006, Word, 0x7201) // moveq #1, d1
	tcpu.write(0x4008, Word, 0x70FF) // moveq #-1, d0
	tcpu.write(0x400A, Word, 0x72FF) // moveq #-1, d1

	tcpu.pc = 0x4000
	tcpu.Step()
	assert.Equal(t, int32(0), tcpu.d[0])
	assert.True(t, tcpu.sr.Z)

	assert.Equal(t, int32(0x4002), tcpu.pc)
	tcpu.Step()
	assert.Equal(t, int32(1), tcpu.d[0])
	assert.False(t, tcpu.sr.Z)

	assert.Equal(t, int32(0x4004), tcpu.pc)
	tcpu.Step()
	assert.Equal(t, int32(1), tcpu.d[0])
	assert.Equal(t, int32(0), tcpu.d[1])
	assert.True(t, tcpu.sr.Z)

	assert.Equal(t, int32(0x4006), tcpu.pc)
	tcpu.Step()
	assert.Equal(t, int32(1), tcpu.d[0])
	assert.Equal(t, int32(1), tcpu.d[1])
	assert.False(t, tcpu.sr.Z)

	assert.Equal(t, int32(0x4008), tcpu.pc)
	tcpu.Step()
	assert.Equal(t, int32(-1), tcpu.d[0])
	assert.Equal(t, int32(1), tcpu.d[1])
	assert.False(t, tcpu.sr.Z)
	assert.True(t, tcpu.sr.N)

	assert.Equal(t, int32(0x400A), tcpu.pc)
	tcpu.Step()
	assert.Equal(t, int32(-1), tcpu.d[0])
	assert.Equal(t, int32(-1), tcpu.d[1])
	assert.False(t, tcpu.sr.Z)
	assert.True(t, tcpu.sr.N)
}
