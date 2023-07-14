package instructions

import (
	. "github.com/jenska/m68k/internal/core"
)

func jmp(c *Core) {
	c.PC = c.FetchEA((c.IRC>>3)&7, c.IRC&7, Long).Address()
}

func init() {
	Register("jmp <ea>", jmp, 0x4ec0, 0xffc0,
		MaskIndirect+
			MaskDisplacement+
			MaskIndex+
			MaskAbsoluteShort+
			MaskAbsoluteLong+
			MaskPCDisplacement+
			MaskPCIndex)

}
