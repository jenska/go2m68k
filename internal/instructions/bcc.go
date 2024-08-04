package instructions

import (
	. "github.com/jenska/m68k/internal/core"
)

func bcc(c *Core) {
	cc := uint8(c.IRC>>8) & 0xf      // condition code
	dis := int32(int8(c.IRC & 0xff)) // signed displacement

	if cc == 1 { // bsr
		if dis == 0 {
			dis = Word.SignedExtend(c.PopPc(Word))
		}
		c.Push(Long, c.PC) // ?
		c.PC = uint32(int32(c.PC0) + dis + WordSize)
	} else if c.SR.TestCC(cc) {
		if dis == 0 {
			dis = Word.SignedExtend(c.PopPc(Word))
		}
		c.PC = uint32(int32(c.PC0) + dis + WordSize)
	} else if dis == 0 {
		c.PC += Word.Size()
	}
}

func init() {
	Register("bcc(.s) dis", bcc, 0x6000, 0xf000, 0x0000)
}
