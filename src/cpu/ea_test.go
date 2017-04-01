package cpu

import (
	"fmt"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEAVectors(t *testing.T) {
	cpu := NewM68k(NewMemoryHandler(1024, nil))
	assert.NotNil(t, cpu)

	eaVec := NewEAVectors(cpu)
	assert.NotNil(t, eaVec)
}

func TestEADataRegister(t *testing.T) {
	cpu := NewM68k(NewMemoryHandler(1024, nil))
	assert.NotNil(t, cpu)
	eaVec := NewEAVectors(cpu)
	assert.NotNil(t, eaVec)
	// reg = D0
	eb := eaVec[0+Byte.eaVecOffset]
	ew := eaVec[0+Word.eaVecOffset]
	el := eaVec[0+Long.eaVecOffset]
	assert.NotNil(t, eb)

	mb := eb.compute()
	mw := ew.compute()
	ml := el.compute()
	assert.NotNil(t, mb)

	ml.write(0)
	mb.write(0xff)
	assert.Equal(t, uint32(0xff), mb.read())
	assert.Equal(t, uint32(0xff), mw.read())
	assert.Equal(t, uint32(0xff), ml.read())

	mw.write(0xffff)
	assert.Equal(t, uint32(0x00ff), mb.read())
	assert.Equal(t, uint32(0xffff), mw.read())
	assert.Equal(t, uint32(0xffff), ml.read())

	mb.write(0xff)
	assert.Equal(t, uint32(0x00ff), mb.read())
	assert.Equal(t, uint32(0xffff), mw.read())
	assert.Equal(t, uint32(0xffff), ml.read())
}

func TestEAAddressRegister(t *testing.T) {
	cpu := NewM68k(NewMemoryHandler(1024, nil))
	eaVec := NewEAVectors(cpu)

	// reg = A0
	const A0 = (1 << 3) | 0
	eb := eaVec[A0+Byte.eaVecOffset]
	ew := eaVec[A0+Word.eaVecOffset]
	el := eaVec[A0+Long.eaVecOffset]
	assert.IsType(t, &EAAddressRegister{}, eb)
	assert.IsType(t, &EAAddressRegister{}, ew)
	assert.IsType(t, &EAAddressRegister{}, el)

	mb := eb.compute()
	mw := ew.compute()
	ml := el.compute()
	assert.NotNil(t, mb)

	ml.write(0)
	mb.write(0xff)
	assert.Equal(t, uint32(0xff), mb.read())
	assert.Equal(t, uint32(0xff), mw.read())
	assert.Equal(t, uint32(0xff), ml.read())

	mw.write(0xffff)
	assert.Equal(t, uint32(0x00ff), mb.read())
	assert.Equal(t, uint32(0xffff), mw.read())
	assert.Equal(t, uint32(0xffff), ml.read())

	mb.write(0xff)
	assert.Equal(t, uint32(0x00ff), mb.read())
	assert.Equal(t, uint32(0xffff), mw.read())
	assert.Equal(t, uint32(0xffff), ml.read())
}

func TestEAIndirect(t *testing.T) {
	cpu := NewM68k(NewMemoryHandler(1024*1024, nil))
	eaVec := NewEAVectors(cpu)

	// (A0)
	const A0 = (2 << 3) | 0

	for i := 0x100; i < 0x120; i++ {
		cpu.A[0] = uint32(i)
		eb := eaVec[A0+Byte.eaVecOffset]
		ew := eaVec[A0+Word.eaVecOffset]
		el := eaVec[A0+Long.eaVecOffset]
		assert.IsType(t, &EAAddressRegisterIndirect{}, eb)
		assert.IsType(t, &EAAddressRegisterIndirect{}, ew)
		assert.IsType(t, &EAAddressRegisterIndirect{}, el)

		mb := eb.compute()
		mw := ew.compute()
		ml := el.compute()
		assert.NotNil(t, mb)

		ml.write(0)
		mb.write(0xff)
		assert.Equal(t, uint32(0xff), mb.read())
		assert.Equal(t, uint32(0xff), mw.read(), fmt.Sprintf("error at address 0x%08x", i))
		assert.Equal(t, uint32(0xff), ml.read(), fmt.Sprintf("error at address 0x%08x", i))

		mw.write(0xffff)
		assert.Equal(t, uint32(0x00ff), mb.read())
		assert.Equal(t, uint32(0xffff), mw.read(), fmt.Sprintf("error at address 0x%08x", i))
		assert.Equal(t, uint32(0xffff), ml.read(), fmt.Sprintf("error at address 0x%08x", i))

		mb.write(0xff)
		assert.Equal(t, uint32(0x00ff), mb.read())
		assert.Equal(t, uint32(0xffff), mw.read())
		assert.Equal(t, uint32(0xffff), ml.read())

		ml.write(0xff01ff01)
		assert.Equal(t, uint32(0x01), mb.read())
		assert.Equal(t, uint32(0xff01), mw.read())
		assert.Equal(t, uint32(0xff01ff01), ml.read())
	}
}
