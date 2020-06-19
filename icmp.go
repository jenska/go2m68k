package cpu

func init() {
	addOpcode("tst", tst, 0x4a00, 0xff00, 0xbf8, "01:4", "7:7", "234fc:2")
}

func tst(c *M68K) {
	size := operands[c.ir>>6&0x3]
	c.sr.setLogicalFlags(size, c.resolveDstEA(size).read())
}
