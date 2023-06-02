package instructions

import (
	. "github.com/jenska/m68k/internal/core"
)

func nop(c *Core) {
}

func init() {
	Register("nop", nop, 0x4e71, 0xffff, 0x0000)
}
