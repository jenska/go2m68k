package m68k

import "testing"

func TestSwap(t *testing.T) {
	c := NewM68k(NewMemoryHandler(0x10000))
	c.Write(Long, 0, 0x0100)
	c.Write(Long, 4, 0xa000)
	c.Reset()
	// TODO
}
