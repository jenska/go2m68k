package instructions

import (
	. "github.com/jenska/m68k/internal/core"
)

func reset(c *Core) {
	if !c.SR.S {
		// privilege violation error
		c.ProcessException(XPrivilegeViolation)
	} else {
		c.Reset()
	}
}

func init() {
	Register("reset", reset, 0x4e70, 0xffff, 0x0000)
}
