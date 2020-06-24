package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbcd(t *testing.T) {
	tcpu.pc = 0x4000

	twrite(0x7011, 0x7222, 0x7433, 0x7644, 0x7801) // moveq d0=0x11, d1=0x22, d2=0x33, d3=0x44, d4=1
	twrite(0xc101)                                 // abcd d1, d0 = 33
	twrite(0xc102)                                 // abcd d2, d0 = 66
	twrite(0xc102)                                 // abcd d2, d0 = 99
	twrite(0xc104)                                 // abcd d4, d0 = 00 + x
	twrite(0xc104)                                 // abcd d4, d0 + x = 2
	twrite(0xc104)                                 // abcd d4, d0  = 3

	tcpu.write(0x1234, Long, 0x11223344)
	tcpu.write(0x1234+Long.size, Long, 0x11223344)
	twrite(0x3440+eaModeImmidiate, 0x1234+4) // movea.w #$1234, a2
	twrite(0x3640+eaModeImmidiate, 0x1238+4) // movea.w #$1234, a3
	twrite(0xc50b)                           // abcd -(a2), -(a3)
	twrite(0xc50b)                           // abcd -(a2), -(a3)
	twrite(0xc50b)                           // abcd -(a2), -(a3)
	twrite(0xc50b)                           // abcd -(a2), -(a3)
	twrite(0x2040+eaModeIndirect+3, 0x1238)  // movea.l (a3), a0
	trun(0x4000)
	// fmt.Println(tcpu)
	assert.Equal(t, int32(3), tcpu.d[0])
	assert.Equal(t, int32(0x22446688), tcpu.a[0])
}

func TestSbcd(t *testing.T) {
	tcpu.pc = 0x4000

	twrite(0x7060, 0x7259) // moveq d0=60, d1=0x59
	twrite(0x8101)         // sbcd d1, d0 = 0

	tcpu.write(0x1234, Long, 0x11223344)
	tcpu.write(0x1234+Long.size, Long, 0x11223344)
	twrite(0x3440+eaModeImmidiate, 0x1234+4) // movea.w #$1234, a2
	twrite(0x3640+eaModeImmidiate, 0x1238+4) // movea.w #$1234, a3
	twrite(0x850b)                           // sbcd -(a2), -(a3)
	twrite(0x850b)                           // sbcd -(a2), -(a3)
	twrite(0x850b)                           // sbcd -(a2), -(a3)
	twrite(0x850b)                           // sbcd -(a2), -(a3)
	twrite(0x2040+eaModeIndirect+3, 0x1238)  // movea.l (a3), a0
	trun(0x4000)
	//	fmt.Println(tcpu)
	assert.Equal(t, int32(1), tcpu.d[0])
	assert.Equal(t, int32(0), tcpu.a[0])
}

func TestNbcd(t *testing.T) {
	tcpu.pc = 0x4000
	twrite(0x7011, 0x7222, 0x7433, 0x7644, 0x7801, 0x7900, 0x7bff) // moveq d0=0x11, d1=0x22, d2=0x33, d3=0x44, d4=1, d5=0, d6=-1

	twrite(0x4800) // nbcd d0
	twrite(0x4801) // nbcd d1
	twrite(0x4802) // nbcd d2
	twrite(0x4803) // nbcd d3
	twrite(0x4804) // nbcd d4
	twrite(0x4805) // nbcd d5
	twrite(0x4806) // nbcd d6
	trun(0x4000)
	//fmt.Println(tcpu)
	assert.Equal(t, int32(0x89), tcpu.d[0])
	assert.Equal(t, int32(0x77), tcpu.d[1])
	assert.Equal(t, int32(0x66), tcpu.d[2])
	assert.Equal(t, int32(0x55), tcpu.d[3])
	assert.Equal(t, int32(0x98), tcpu.d[4])
	assert.Equal(t, int32(0x99), tcpu.d[5])
	assert.Equal(t, int32(0x99), tcpu.d[6])

}
