package instructions

import (
	. "github.com/jenska/m68k/internal/core"
)

func rts(c *Core) {
	c.PC = c.Pop(Long)
}

func init() {
	Register("rts", rts, 0x4e75, 0xffff, 0x0000)
}
