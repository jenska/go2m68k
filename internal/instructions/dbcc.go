package instructions

import (
	. "github.com/jenska/m68k/internal/core"
)

func dbcc(c *Core) {
	if c.SR.TestCC(c.IRC >> 8) {
		c.PC += Word.Size() // skip displacement value
	} else {
		regPtr := &c.D[c.IRC&0x07]
		counter := *regPtr - 1
		dis := Word.SignedExtend(c.PopPc(Word))
		Word.WriteToLong(counter, *&regPtr)

		if counter != 0xffff {
			c.PC = uint32(int32(c.PC0) + dis)
		}
	}
}

func init() {
	Register("dbcc", dbcc, 0x50c8, 0xf0f8, 0x0000)
}
