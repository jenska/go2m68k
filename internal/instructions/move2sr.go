package instructions

import (
	. "github.com/jenska/m68k/internal/core"
)

func move2ccr(c *Core) {
	ea := c.FetchEA((c.IRC>>3)&7, c.IRC&7, Byte)
	c.SR = NewStatusRegister(c.SR.Word() | uint16(ea.Read()))
}

func move2sr(c *Core) {
	if c.SR.S {
		ea := c.FetchEA((c.IRC>>3)&7, c.IRC&7, Word)
		c.SR = NewStatusRegister(uint16(ea.Read()))
	} else {
		c.ProcessException(XPrivilegeViolation)
	}
}

func init() {
	var eaMask uint16 = MaskDataRegister |
		MaskAbsoluteLong |
		MaskAbsoluteShort |
		MaskDisplacement |
		MaskIndex |
		MaskIndirect |
		MaskPostIncrement |
		MaskPreDecrement |
		MaskImmediate

	Register("move <ea>, ccr", move2ccr, 0x44c0, 0xffc0, eaMask)
	Register("move <ea>, sr", move2sr, 0x46c0, 0xffc0, eaMask)
}
