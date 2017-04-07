package m68

import (
	"fmt"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEAVectors(t *testing.T) {
	cpu := NewM68k(NewMemoryHandler(1024))
	assert.NotNil(t, cpu)

	eaVec := NewEAVectors(cpu)
	assert.NotNil(t, eaVec)
}

func TestEADataRegister(t *testing.T) {
	cpu := NewM68k(NewMemoryHandler(1024))
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
	cpu := NewM68k(NewMemoryHandler(1024))
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
	cpu := NewM68k(NewMemoryHandler(1024))
	eaVec := NewEAVectors(cpu)

	peek := func(a uint32) uint32 { return cpu.read(Byte, a) }

	// (A0)
	const A0 = (2 << 3) | 0

	for i := uint32(0x102); i < 0x110; i += 2 {
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
		assert.Equal(t, uint32(0), peek(i))
		assert.Equal(t, uint32(0), peek(i+1))
		assert.Equal(t, uint32(0), peek(i+2))
		assert.Equal(t, uint32(0), peek(i+3))
		mb.write(0xff)
		assert.Equal(t, uint32(0xff), peek(i))
		assert.Equal(t, uint32(0), peek(i+1))
		assert.Equal(t, uint32(0), peek(i+2))
		assert.Equal(t, uint32(0), peek(i+3))

		assert.Equal(t, uint32(0xff), mb.read())
		assert.Equal(t, uint32(0xff), mw.read(), fmt.Sprintf("error at address 0x%08x %02x %02x",
			i, peek(i), peek(i+1)))
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

func TestEAPostInc(t *testing.T) {
	cpu := NewM68k(NewMemoryHandler(1024))
	eaVec := NewEAVectors(cpu)

	// (A0)+
	const A0 = (3 << 3) | 0
	cpu.A[0] = uint32(0x100)
	eb := eaVec[A0+Byte.eaVecOffset]
	ew := eaVec[A0+Word.eaVecOffset]
	el := eaVec[A0+Long.eaVecOffset]

	mb := eb.compute()
	assert.Equal(t, uint32(0x101), cpu.A[0])
	mb.write(0xff)
	assert.Equal(t, uint32(0xff), cpu.read(Byte, 0x100))

	cpu.A[0] = uint32(0x100)
	mw := ew.compute()
	assert.Equal(t, uint32(0x102), cpu.A[0])
	mw.write(0xffff)
	assert.Equal(t, uint32(0xffff), cpu.read(Word, 0x100))

	cpu.A[0] = uint32(0x100)
	ml := el.compute()
	assert.Equal(t, uint32(0x104), cpu.A[0])
	ml.write(0xffffffff)
}

func TestEAPreDec(t *testing.T) {
	cpu := NewM68k(NewMemoryHandler(1024))
	eaVec := NewEAVectors(cpu)

	// -(A0)
	const A0 = (4 << 3) | 0

	eb := eaVec[A0+Byte.eaVecOffset]
	ew := eaVec[A0+Word.eaVecOffset]
	el := eaVec[A0+Long.eaVecOffset]
	assert.IsType(t, &EAAddressRegisterPreDec{}, eb)
	assert.IsType(t, &EAAddressRegisterPreDec{}, ew)
	assert.IsType(t, &EAAddressRegisterPreDec{}, el)

	cpu.A[0] = 0x100
	mb := eb.compute()
	assert.Equal(t, uint32(0xff), cpu.A[0])
	mb.write(0xff)
	assert.Equal(t, uint32(0xff), cpu.read(Byte, 0xFF))

	mw := ew.compute()
	ml := el.compute()

	mb.write(0xff)
	assert.Equal(t, uint32(0xff), cpu.read(Byte, 0xff))
	mw.write(0xffff)
	ml.write(0)
}

func TestEAAddressRegisterWithDisplacement(t *testing.T) {
	cpu := NewM68k(NewMemoryHandler(1024))
	eaVec := NewEAVectors(cpu)

	// xxxx(A0)
	const A0 = (5 << 3) | 0

	eb := eaVec[A0+Byte.eaVecOffset]
	ew := eaVec[A0+Word.eaVecOffset]
	el := eaVec[A0+Long.eaVecOffset]
	assert.IsType(t, &EAAddressRegisterWithDisplacement{}, eb)
	assert.IsType(t, &EAAddressRegisterWithDisplacement{}, ew)
	assert.IsType(t, &EAAddressRegisterWithDisplacement{}, el)

	cpu.A[0] = 0x100
	cpu.PC = 0x200
	cpu.pushPC(Word, 0x80)
	cpu.write(Byte, 0x180, 0x23)
	mb := eb.compute()
	assert.Equal(t, uint32(0x23), mb.read())

	cpu.A[0] = 0x100
	cpu.PC = 0x200
	cpu.pushPC(Word, 0x80)
	cpu.write(Word, 0x180, 0x1234)
	mw := ew.compute()
	assert.Equal(t, uint32(0x1234), mw.read())

}

func TestEAAddressRegisterWithIndex(t *testing.T) {

}
