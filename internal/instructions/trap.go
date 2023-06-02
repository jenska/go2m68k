package instructions

import (
	. "github.com/jenska/m68k/internal/core"
)

func trap(c *Core) {
	c.ProcessException(XTrapBase + (c.IRC & 0xf))
}

func init() {
	Register("trap #imm", trap, 0x4e40, 0xfff0, 0x000)
}
