package instructions

import (
	"github.com/jenska/m68k/internal/core"
	. "github.com/jenska/m68k/internal/core"
)

func init() {
	Register("cmpi #imm, <ea>", cmpi, 0x0c00, 0xff00,
		MaskDataRegister|
			MaskAbsoluteLong|
			MaskAbsoluteShort|
			MaskDisplacement|
			MaskIndex|
			MaskIndirect|
			MaskPostIncrement|
			MaskPreDecrement)
}

func cmpi(c *core.Core) {
	if op := OperandFromValue(c.IRC >> 6); op != nil {
		val1 := c.PopPc(op)
		val2 := c.FetchEA((c.IRC>>3)&7, c.IRC&7, op).Read()
		c.SR.SetFlags(op, val2-val1)
	} else {
		c.ProcessException(XIllegalInstruction)
	}

}
