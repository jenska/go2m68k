package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMovea(t *testing.T) {
	tcpu.pc = 0x4000
	tcpu.write(0x1234, Word, 0x4321)

	twrite(0x3440+eaModeImmidiate, 0x1234)         // movea.w #$1234, a2
	twrite(0x2640+eaModeImmidiate, 0x1234, 0x5678) // movea.l #$12345678, a3

	twrite(0x3840 + eaModePostIncrement + 2) // movea.w (a2)+, a4
	twrite(0x3a40 + eaModePreDecrement + 2)  // movea.w -(a2), a5

	trun(0x4000)
	assert.Equal(t, int32(0x1234), tcpu.a[2])
	assert.Equal(t, int32(0x12345678), tcpu.a[3])
	assert.Equal(t, int32(0x4321), tcpu.a[4])
	assert.Equal(t, int32(0x4321), tcpu.a[5])
}

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
