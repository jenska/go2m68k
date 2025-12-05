package instructions

import (
	"github.com/jenska/m68k/internal/core"
	. "github.com/jenska/m68k/internal/core"
)

func init() {
	Register("pea <ea>", pea, 0x41c0, 0xf1c0,
		MaskAbsoluteLong|
			MaskAbsoluteShort|
			MaskDisplacement|
			MaskIndex|
			MaskIndirect|
			MaskPCDisplacement|
			MaskPCIndex|
			MaskPostIncrement|
			MaskPreDecrement)
}

func pea(c *core.Core) {
	ax := &c.A[(c.IRC>>9)&7]
	*ax = c.FetchEA((c.IRC>>3)&7, c.IRC&7, Long).Address()
}
