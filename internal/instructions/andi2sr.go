package instructions

import (
	. "github.com/jenska/m68k/internal/core"
)

func andi2ccr(c *Core) {
	c.SR = NewStatusRegister(c.SR.Word() & uint16(c.PopPc(Byte)))
}

func andi2sr(c *Core) {
	if c.SR.S {
		c.SR = NewStatusRegister(c.SR.Word() & uint16(c.PopPc(Word)))
	} else {
		c.ProcessException(XPrivilegeViolation)
	}
}

func init() {
	Register("andi #imm, ccr", andi2ccr, 0x23c, 0xffff, 0x0000)
	Register("andi #imm, sr", andi2sr, 0x27c, 0xffff, 0x0000)
}
