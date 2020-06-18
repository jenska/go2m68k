package cpu

func init() {
	addOpcode("clr", clr, 0x4200, 0xff00, 0xbf8, "01:4", "7:7", "234fc:2") // Word cycles
	addOpcode("swap", swap, 0x4840, 0xfff8, 0x0000, "01234fc:4", "7:7")
	addOpcode("reset", reset, 0x4e70, 0xffff, 0x0000, "071234fc:0")
	addOpcode("nop", func(c *M68K) {}, 0x4e71, 0xffff, 0x0000, "01:4", "7:7", "234fc:2:0")
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
		panic(PrivilegeViolationError)
	}
	c.Reset()
	// c.cycles -=
}
