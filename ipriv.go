package cpu

func init() {
	addOpcode("stop #sr", stop, 0x4e72, 0xffff, 0x0000, "01:4", "7:13", "234fc:8")
}

func stop(c *M68K) {
	if c.sr.S {
		newSR := c.popPc(Word)
		c.stopped = true
		c.sr.setbits(newSR)
	} else {
		panic(PrivilegeViolationError)
	}
}
