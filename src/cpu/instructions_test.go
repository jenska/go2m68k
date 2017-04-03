package cpu

import "testing"

func TestSwap(t *testing.T) {
	c := NewM68k(NewMemoryHandler(0x10000))
	c.write(Long, 0, 0x0100)
	c.write(Long, 4, 0xa000)
	c.Reset()
	// TODO
}
