package instructions

import (
	. "github.com/jenska/m68k/internal/core"
)

func eori2ccr(c *Core) {
	c.SR = NewStatusRegister(c.SR.Word() ^ uint16(c.PopPc(Byte)))
}
func eori2sr(c *Core) {
	if c.SR.S {
		c.SR = NewStatusRegister(c.SR.Word() ^ uint16(c.PopPc(Word)))
	} else {
		c.ProcessException(XPrivilegeViolation)
	}
}

func init() {
	Register("eori #imm, ccr", eori2ccr, 0xa3c, 0xffff, 0x0000)
	Register("eori #imm, sr", eori2sr, 0xa7c, 0xffff, 0x0000)
}
