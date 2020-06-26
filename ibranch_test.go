package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBra16(t *testing.T) {
	tcpu.pc = 0x4000
	twrite(0x6000) // bra
	twrite(0x0002) // #+2
	twrite(0x7001) // moveq #1, d0
	twrite(0x7002) // moveq #2, d0
	twrite(0x4e72) // stop
	twrite(0x2700) // #$27000
	trun(0x4000)
	assert.Equal(t, int32(2), tcpu.d[0])
}

func TestBsr8(t *testing.T) {
	tcpu.pc = 0x4000
	twrite(0x7001) // moveq #1, d0
	twrite(0x6104) // bsr
	twrite(0x4e72) // stop
	twrite(0x2700) // #$27000
	twrite(0x4e71) // nop
	twrite(0x7002) // moveq #2, d0
	twrite(0x4e75) // rts
	trun(0x4000)
	assert.Equal(t, int32(2), tcpu.d[0])
}

func TestBsr16(t *testing.T) {
	tcpu.pc = 0x4000
	twrite(0x7001) // moveq #1, d0
	twrite(0x6100) // bsr
	twrite(0x0004) // #+4
	twrite(0x4e72) // stop
	twrite(0x2700) // #$27000
	twrite(0x4e71) // nop
	twrite(0x7002) // moveq #2, d0
	twrite(0x4e75) // rts
	trun(0x4000)
	assert.Equal(t, int32(2), tcpu.d[0])
}

func TestDbra(t *testing.T) {
	tcpu.write(0x4000, Word, 0x7005) // moveq #5, d0
	tcpu.write(0x4002, Word, 0x51c8) // dbra d0,
	tcpu.write(0x4004, Word, 0xfffe) // -2
	tcpu.write(0x4006, Word, 0x7201) // moveq #1, d1
	tcpu.write(0x4008, Word, 0x70FF) // moveq #-1, d0
	tcpu.write(0x400A, Word, 0x72FF) // moveq #-1, d1

	tcpu.pc = 0x4000
	tcpu.Step()
	for i := 5; i > 0; i-- {
		assert.Equal(t, int32(i), tcpu.d[0])
		assert.Equal(t, int32(0x4002), tcpu.pc)
		tcpu.Step()
		assert.Equal(t, int32(0x4002), tcpu.pc)
	}
	tcpu.Step()
	assert.Equal(t, int32(0x4006), tcpu.pc)

	tcpu.write(0x4000, Word, 0x7000+5) // moveq #5, d0
	tcpu.write(0x4002, Word, 0x7200+5) // moveq #5, d1
	tcpu.write(0x4004, Word, 0x51c9)   // dbra d1,
	tcpu.write(0x4006, Word, 0xfffe)   // -2
	tcpu.write(0x4008, Word, 0x51c8)   // dbra d0,
	tcpu.write(0x400a, Word, 0xfff8)   // -8
	tcpu.write(0x400c, Word, 0x4e72)   // stop
	tcpu.write(0x400e, Word, 0x2300)   // #$23000
	tcpu.pc = 0x4000
	signals := make(chan Signal)
	tcpu.Run(signals)
	assert.Equal(t, int32(0x2300), tcpu.sr.bits())
}

func BenchmarkDbra(b *testing.B) {
	b.StopTimer()
	tcpu.pc = 0x4000
	twrite(0x7000 + 100)   // moveq #100, d0
	twrite(0x7200 + 100)   // moveq #100, d1
	twrite(0x4e71)         // nop
	twrite(0x51c9, 0xfffc) // dbra d1, #-4
	twrite(0x51c8, 0xfff6) // dbra d0, #-10
	twrite(0x4e72, 0x2700)
	signals := make(chan Signal)
	b.StartTimer()
	for j := 0; j < b.N; j++ {
		tcpu.pc = 0x4000
		tcpu.Run(signals)
	}
	// fmt.Println(tcpu.icount)
}
