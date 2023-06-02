package instructions

import (
	. "github.com/jenska/m68k/internal/core"
)

func ori2ccr(c *Core) {
	c.SR = NewStatusRegister(c.SR.Word() | uint16(c.PopPc(Byte)))
}

func ori2sr(c *Core) {
	if c.SR.S {
		c.SR = NewStatusRegister(c.SR.Word() | uint16(c.PopPc(Word)))
	} else {
		c.ProcessException(XPrivilegeViolation)
	}
}

func init() {
	Register("ori #imm, ccr", ori2ccr, 0x3c, 0xffff, 0x0000)
	Register("ori #imm, sr", ori2sr, 0x7c, 0xffff, 0x0000)
}
