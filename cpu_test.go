package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	assert.Equal(t, initialSSP, tcpu.readAddress(0))
	assert.Equal(t, initialSSP&0xff, tcpu.read(3, Byte))
	assert.Equal(t, initialSSP&0xffff, tcpu.read(2, Word))
	assert.Equal(t, initialSSP, tcpu.read(0, Long))
	// bounds overflow
	assert.Equal(t, initialSSP, tcpu.readAddress(0x12000000))
}

func TestReadWrite(t *testing.T) {
	assert.Panics(t, func() {
		tcpu.write(0, Long, 400)
	})
	assert.Panics(t, func() {
		tcpu.write(4, Long, 400)
	})
	assert.Equal(t, initialSSP, tcpu.readAddress(0))
	assert.Equal(t, initialPC, tcpu.readAddress(4))

	tcpu.write(100, Long, 3)
	assert.Equal(t, int32(3), tcpu.readAddress(100))
	assert.Equal(t, int32(3), tcpu.read(103, Byte))
	assert.Equal(t, int32(3), tcpu.read(102, Word))
	assert.Equal(t, int32(3), tcpu.read(100, Long))
}

func TestReset(t *testing.T) {
	tcpu.Reset()
	assert.True(t, tcpu.sr.S)
	assert.Equal(t, initialPC, tcpu.pc)
	assert.Equal(t, initialSSP, tcpu.ssp)
	assert.Equal(t, int32(0x2700), tcpu.sr.bits())
}

func TestRaiseException(t *testing.T) {
	tcpu.Reset()
	tcpu.write(int32(PrivilegeViolationError)<<2, Long, 500)
	tcpu.raiseException(PrivilegeViolationError)
	assert.Equal(t, int32(500), tcpu.pc)

	assert.Panics(t, func() {
		tcpu.raiseException(ZeroDivideError)
	})

	tcpu.write(int32(UnintializedInterrupt)<<2, Long, 600)
	tcpu.raiseException(ZeroDivideError)
	assert.Equal(t, int32(600), tcpu.pc)

	tcpu.sr.S = false
	tcpu.raiseException(ZeroDivideError)
	assert.Equal(t, int32(600), tcpu.pc)
	assert.True(t, tcpu.sr.S)
}

func TestPrivileViolationException(t *testing.T) {
	oldS := tcpu.sr.S
	tcpu.sr.S = false
	defer func() { tcpu.sr.S = oldS }()

	tcpu.read(0, Long)
	tcpu.read(4, Long)

	assert.Panics(t, func() {
		tcpu.write(0, Long, 0)
	})

	oldV := tcpu.read(100, Long)
	assert.PanicsWithValue(t, PrivilegeViolationError, func() {
		tcpu.write(100, Long, 0)
	})

	tcpu.sr.S = true
	assert.NotPanics(t, func() {
		tcpu.write(100, Long, oldV)
	})
}

func TestAddressError(t *testing.T) {
	oldV := tcpu.read(100, Long)

	assert.PanicsWithValue(t, AdressError, func() {
		tcpu.read(101, Long)
	})
	assert.PanicsWithValue(t, AdressError, func() {
		tcpu.read(101, Word)
	})
	assert.PanicsWithValue(t, AdressError, func() {
		tcpu.write(101, Long, 0)
	})
	assert.PanicsWithValue(t, AdressError, func() {
		tcpu.write(101, Word, 0)
	})

	tcpu.write(100, Long, oldV)
}

func TestPop(t *testing.T) {
	assert.NotEqual(t, 0, tcpu.a[7])
	tcpu.push(Long, 1001)
	assert.Equal(t, int32(1001), tcpu.pop(Long))
}

func TestBra8(t *testing.T) {
	tcpu.Reset()
	tcpu.pc = romTop
	tcpu.Step()
	assert.Equal(t, romTop+0x30, tcpu.pc)
	assert.Equal(t, 1, tcpu.icount)
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
	tcpu.write(0x400e, Word, 0x2300)   // #$27000
	tcpu.pc = 0x4000
	signals := make(chan Signal)
	tcpu.Run(signals)
	assert.Equal(t, int32(0x2300), tcpu.sr.bits())
}

func TestStop(t *testing.T) {
	signals := make(chan Signal)
	tcpu.write(0x400c, Word, 0x4e72) // stop
	tcpu.write(0x400e, Word, 0x2000) // #$27000
	tcpu.pc = 0x400c
	tcpu.Run(signals)
	assert.True(t, tcpu.stopped)
	assert.Equal(t, int32(0x2000), tcpu.sr.bits())

	tcpu.pc = 0x400c
	tcpu.sr.S = false
	assert.Panics(t, func() {
		tcpu.Run(signals)
	})
}

func BenchmarkDbra(b *testing.B) {
	b.StopTimer()
	tcpu.write(0x4000, Word, 0x7000+100) // moveq #100, d0
	tcpu.write(0x4002, Word, 0x7200+100) // moveq #100, d1
	tcpu.write(0x4004, Word, 0x51c9)     // dbra d1,
	tcpu.write(0x4006, Word, 0xfffe)     // -2
	tcpu.write(0x4008, Word, 0x51c8)     // dbra d0,
	tcpu.write(0x400a, Word, 0xfff8)     // -8
	tcpu.write(0x400c, Word, 0x4e72)     // stop
	tcpu.write(0x400e, Word, 0x2700)     // #$27000
	signals := make(chan Signal)
	b.StartTimer()
	for j := 0; j < b.N; j++ {
		tcpu.pc = 0x4000
		tcpu.Run(signals)
	}
}
