package cpu

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDisassemble(t *testing.T) {
	assert.NotNil(t, tcpu)
	assert.NotNil(t, trom)
	assert.NotNil(t, trom.Raw())

	bus := tcpu.bus
	bra := bus.Read(romTop, Word)
	assert.Equal(t, 0x602e, bra)
	d, next := Disassemble(romTop, bus)
	assert.Equal(t, "bra.s $00fc0030", d )
	assert.Equal(t, uint32(romTop+2), next)
}

func ExampleDisassemble() {

}
