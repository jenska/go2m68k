package cpu

func init() {
	addOpcode("clr", clr, 0x4200, 0xff00, 0xbf8, "01:4", "7:7", "234fc:2") // Word cycles
	addOpcode("swap", swap, 0x4840, 0xfff8, 0x0000, "01234fc:4", "7:7")
	addOpcode("reset", reset, 0x4e70, 0xffff, 0x0000, "071234fc:0")
	addOpcode("nop", func(c *M68K) {}, 0x4e71, 0xffff, 0x0000, "01:4", "7:7", "234fc:2:0")
	addOpcode("illegal", illegal, 0x4afc, 0xffff, 0x0000, "071234fc:4")
	addOpcode("stop #sr", stop, 0x4e72, 0xffff, 0x0000, "01:4", "7:13", "234fc:8")
	addOpcode("lea ea, ax", lea, 0x41c0, 0xf1c0, 0x027b, "01:0", "7:7", "234fc:2")
}

func lea(c *M68K) {
	*ax(c) = c.resolveDstEA(Long).computedAddress()
}

func clr(c *M68K) {
	size := operands[(c.ir>>6)&0x3]
	ea := c.resolveDstEA(size)
	c.sr.setLogicalFlags(size, 0)
	ea.write(0)
}

func swap(c *M68K) {
	d := *dy(c)
	d = (d>>16)&0xffff | (d << 16)
	c.sr.setLogicalFlags(Long, d)
	*dy(c) = d
}

func reset(c *M68K) {
	if !c.sr.S {
		panic(NewError(PrivilegeViolationError, c, c.pc, nil))
	}
	c.Reset()
	// c.cycles -=
}

func illegal(c *M68K) {
	// if c.ir != 0x4afc {
	// 	debug.PrintStack()
	// }
	panic(NewError(IllegalInstruction, c, c.pc, nil))
}

func stop(c *M68K) {
	if c.sr.S {
		newSR := c.popPc(Word)
		c.stopped = true
		c.sr.setbits(newSR)
	} else {
		panic(NewError(PrivilegeViolationError, c, c.pc, nil))
	}
}
