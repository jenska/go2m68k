package instructions

import (
	. "github.com/jenska/m68k/internal/core"
)

func moveq(c *Core) {
	res := uint32(Byte.SignedExtend(uint32(c.IRC)))
	c.SR.SetFlags(Long, res)
	c.D[(c.IRC>>9)&7] = res
}

func init() {
	Register("moveq #i, dy", moveq, 0x7000, 0xf100, 0x0000)
}
