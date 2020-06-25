package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImmidiate(t *testing.T) {
	tcpu.pc = 0x4000
	twrite(0x41c0+eaModeImmidiate, 0x001) // lea #1, a0 ==> illegal instruction error
	tcpu.write(int32(IllegalInstruction)<<2, Long, 0x5000)
	tcpu.pc = 0x4000
	tcpu.Step()
	assert.Equal(t, int32(0x5000), tcpu.pc)
}

func TestIndirect(t *testing.T) {
	tcpu.pc = 0x4000
	twrite(0x45c0+eaModeAbsoluteLong, 0x0000, 0x4000) // lea $4000.l, a2
	twrite(0x45c0 + eaModeIndirect + 2)               // lea (a2), a2
	trun(0x4000)
	assert.Equal(t, int32(0x4000), tcpu.a[2])
}

func TestDisplacement(t *testing.T) {
	tcpu.pc = 0x4000
	twrite(0x45c0+eaModeAbsoluteLong, 0x0000, 0x4000) // lea $4000.l, a2
	twrite(0x45c0+eaModeDisplacement+2, 0x0100)       // lea $100(a2), a2
	twrite(0x45c0+eaModeDisplacement+2, 0xFF00)       // lea -$100(a2), a2
	trun(0x4000)
	assert.Equal(t, int32(0x4000), tcpu.a[2])
}

func TestAbsolute(t *testing.T) {
	tcpu.pc = 0x4000
	twrite(0x41c0+eaModeAbsoluteShort, 0x0001)        // lea $1.w, a0
	twrite(0x43c0+eaModeAbsoluteLong, 0x0000, 0x0001) // lea $1.l, a1
	twrite(0x45c0+eaModeAbsoluteLong, 0x0000, 0x5000) // lea $1.l, a2
	trun(0x4000)
	assert.Equal(t, int32(1), tcpu.a[0])
	assert.Equal(t, int32(1), tcpu.a[1])
	assert.Equal(t, int32(0x5000), tcpu.a[2])
}
