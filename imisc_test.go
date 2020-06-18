package cpu

import (
	"testing"

	"github.com/magiconair/properties/assert"
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
