package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBits(t *testing.T) {
	ssr := SSR{}
	assert.Equal(t, 0, ssr.ToBits())
	ssr.S = true
	assert.Equal(t, 0x2000, ssr.ToBits())
	ssr.T0 = true
	assert.Equal(t, 0x6000, ssr.ToBits())
}

func TestBits(t *testing.T) {
	ssr := SSR{}
	assert.Equal(t, 0, ssr.ToBits())
	ssr.Bits(0x2000)
	assert.True(t, ssr.S)
	assert.Equal(t, 0x2000, ssr.ToBits())
}

func TestAllBits(t *testing.T) {
	ssr := SSR{C: true, V: true, Z: true, N: true, X: true, S: true, T1: true, T0: true, M: true,
		Interrupts: 7,
	}
	bits := ssr.ToBits()
	ssr.Bits(bits)
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
	ssr := SSR{C: true, V: true, Z: true, N: true, X: true, S: true, T1: true, T0: true, M: true,
		Interrupts: 7,
	}
	bits := ssr.ToBits()
	ccr := ssr.GetCCR()
	assert.Equal(t, bits&0xff, ccr)

	ssr.SetCCR(bits)
	bits2 := ssr.ToBits()
	assert.Equal(t, bits, bits2)

	ssr.S = false
	bits = ssr.ToBits()
	ccr = ssr.GetCCR()
	assert.Equal(t, bits&0xff, ccr)
}
