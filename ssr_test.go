package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBits(t *testing.T) {
	sr := ssr{}
	assert.Equal(t, 0, sr.bits())
	sr.S = true
	assert.Equal(t, 0x2000, sr.bits())
	sr.T0 = true
	assert.Equal(t, 0x6000, sr.bits())
}

func TestBits(t *testing.T) {
	sr := &ssr{}
	assert.Equal(t, 0, sr.bits())
	sr.setbits(0x2000)
	assert.True(t, sr.S)
	assert.Equal(t, 0x2000, sr.bits())
}

func TestAllBits(t *testing.T) {
	ssr := ssr{C: true, V: true, Z: true, N: true, X: true, S: true, T1: true, T0: true, M: true,
		Interrupts: 7,
	}
	bits := ssr.bits()
	ssr.setbits(bits)
	assert.True(t, ssr.C)
	assert.True(t, ssr.V)
	assert.True(t, ssr.Z)
	assert.True(t, ssr.N)
	assert.True(t, ssr.X)
	assert.True(t, ssr.S)
	assert.True(t, ssr.T1)
	assert.True(t, ssr.T0)
	assert.True(t, ssr.M)
	assert.Equal(t, 7, ssr.Interrupts)
}

func TestCCR(t *testing.T) {
	ssr := ssr{C: true, V: true, Z: true, N: true, X: true, S: true, T1: true, T0: true, M: true,
		Interrupts: 7,
	}
	bits := ssr.bits()
	ccr := ssr.ccr()
	assert.Equal(t, bits&0xff, ccr)

	ssr.setccr(bits)
	bits2 := ssr.bits()
	assert.Equal(t, bits, bits2)

	ssr.S = false
	bits = ssr.bits()
	ccr = ssr.ccr()
	assert.Equal(t, bits&0xff, ccr)
}
