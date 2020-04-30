package cpu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisassemble(t *testing.T) {
	assert.NotNil(t, tcpu)
	assert.NotNil(t, trom)

	bus := tcpu.bus
	bra := bus.read(romTop, Word)
	assert.Equal(t, int32(0x602e), bra)
	d := Disassembler(romTop, bus)
	assert.Equal(t, "00fc0000 bra.s      $00fc0030", d.Next().String())

	d = Disassembler(0xfc0030, bus)
	assert.Equal(t, "00fc0030 move       #$2700, sr", d.Next().String())
	for i := 0; i < 10; i++ {
		fmt.Println(d.Next())
	}
}
