package cpu

func init() {
	addOpcode(a000, 0xa000, 0xf000, 0x0000, "071234fc:4")
	addOpcode(f000, 0xf000, 0xf000, 0x0000, "071234fc:4")
	addOpcode(fpu0, 0xf200, 0xff00, 0x0000, "234f:0")
	addOpcode(fpu1, 0xf300, 0xff00, 0x0000, "234f:0")
}

func a000(c *M68K) {
	c.raiseException1010()
}

func f000(c *M68K) {
	c.raiseException1111()
}

func fpu0(c *M68K) {
	if c.hasFPU {
		//	c.m68040_fpu_op0()
	} else {
		c.raiseException1111()
	}
}

func fpu1(c *M68K) {
	if c.hasFPU {
		//	c.m68040_fpu_op1()
	} else {
		c.raiseException1111()
	}
}

func (c *M68K) raiseException1111() {
	panic("not implemented")
}

func (c *M68K) raiseException1010() {
	panic("not implemented")
}
