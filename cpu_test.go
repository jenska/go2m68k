package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	assert.Equal(t, initialSSP, tcpu.read(0, Long))
	assert.Equal(t, initialSSP&0xff, tcpu.read(3, Byte))
	assert.Equal(t, initialSSP&0xffff, tcpu.read(2, Word))
	assert.Equal(t, initialSSP, tcpu.read(0, Long))
	// bounds overflow
	assert.Equal(t, initialSSP, tcpu.read(0x12000000, Long))
}

func TestReadWrite(t *testing.T) {
	assert.Panics(t, func() {
		tcpu.write(0, Long, 400)
	})
	assert.Panics(t, func() {
		tcpu.write(4, Long, 400)
	})
	assert.Equal(t, initialSSP, tcpu.read(0, Long))
	assert.Equal(t, initialPC, tcpu.read(4, Long))

	tcpu.write(100, Long, 3)
	assert.Equal(t, int32(3), tcpu.read(100, Long))
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
	assert.Panics(t, func() {
		tcpu.write(100, Long, 0)
	})

	tcpu.sr.S = true
	assert.NotPanics(t, func() {
		tcpu.write(100, Long, oldV)
	})
}

func TestAddressError(t *testing.T) {
	oldV := tcpu.read(100, Long)

	assert.Panics(t, func() {
		tcpu.read(101, Long)
	})
	assert.Panics(t, func() {
		tcpu.read(101, Word)
	})
	assert.Panics(t, func() {
		tcpu.write(101, Long, 0)
	})
	assert.Panics(t, func() {
		tcpu.write(101, Word, 0)
	})

	tcpu.write(100, Long, oldV)
}

func TestPop(t *testing.T) {
	assert.NotEqual(t, 0, tcpu.a[7])
	tcpu.push(Long, 1001)
	assert.Equal(t, int32(1001), tcpu.pop(Long))
}
