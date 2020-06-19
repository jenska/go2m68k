package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClr(t *testing.T) {
	tcpu.write(0x4000, Word, 0x7000+100) // moveq #100, d0
	tcpu.write(0x4002, Word, 0x7200+100) // moveq #100, d1
	tcpu.write(0x4004, Word, 0x4200)     // clr.b d0
	tcpu.write(0x4006, Word, 0x4201)     // clr.b d1
	tcpu.write(0x4008, Word, 0x4e72)     // stop
	tcpu.write(0x400a, Word, 0x2700)     // #$27000
	tcpu.pc = 0x4000
	signals := make(chan Signal)
	tcpu.Run(signals)
	if tcpu.d[0] != 0 {
		t.Error("must be 0")
	}
	if tcpu.d[1] != 0 {
		t.Error("must be 0")
	}

	tcpu.pc = 0x4000
	twrite(0x7000 + 100) // moveq #100, d0
	twrite(0x4840)       // swap d0
	twrite(0x4200)       // clr.b d0
	twrite(0x4840)       // swap d0
	twrite(0x4a00)       // tst.b d0
	// twrite(0x66)
	twrite(0x4240) // clr.w d0
	twrite(0x4280) // clr.l d0
	trun(0x4000)
	assert.Equal(t, int32(0), tcpu.d[0])
}

func TestIllegal(t *testing.T) {
	tcpu.pc = 0x4000
	twrite(0x4afc) // illegal
	trun(0x4000)
}

func TestStop(t *testing.T) {
	signals := make(chan Signal)
	tcpu.write(0x400c, Word, 0x4e72) // stop
	tcpu.write(0x400e, Word, 0x2000) // #$2000
	tcpu.pc = 0x400c
	tcpu.Run(signals)
	assert.True(t, tcpu.stopped)
	assert.Equal(t, int32(0x2000), tcpu.sr.bits())

	tcpu.pc = 0x400c
	tcpu.write(int32(PrivilegeViolationError)<<2, Long, 0x400c)
	tcpu.sr.S = false
	tcpu.Run(signals)
	assert.Equal(t, int32(0x400c), tcpu.pc)
	assert.True(t, tcpu.sr.S)
	assert.Equal(t, int32(0x2000), tcpu.sr.bits())
}

func TestSwap(t *testing.T) {
	tcpu.write(0x4000, Word, 0x7000+100) // moveq #100, d0
	tcpu.write(0x4002, Word, 0x7200+100) // moveq #100, d1
	tcpu.write(0x4004, Word, 0x4840)     // swap d0
	tcpu.write(0x4006, Word, 0x4841)     // swap d1
	tcpu.write(0x4008, Word, 0x4e72)     // stop
	tcpu.write(0x400a, Word, 0x2700)     // #$27000
	tcpu.pc = 0x4000
	signals := make(chan Signal)
	tcpu.Run(signals)
	assert.Equal(t, tcpu.d[0], int32(100<<16))
	assert.Equal(t, tcpu.d[1], int32(100<<16))
}
