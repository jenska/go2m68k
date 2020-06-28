package cpu

func init() {
	addOpcode("tst", tst, 0x4a00, 0xff00, 0xbf8, "01:4", "7:7", "234fc:2")
	addOpcode("tas ea", tas, 0x4ac0, 0xffc0, 0x0bf8, "01:14", "7:15", "234fc:12")
}

func tst(c *M68K) {
	size := operands[c.ir>>6&0x3]
	c.sr.setLogicalFlags(size, c.resolveDstEA(size).read())
}

func tas(c *M68K) {
	ea := c.resolveSrcEA(Byte)
	res := ea.read()
	c.sr.setLogicalFlags(Byte, res)
	ea.write(res | 0x80)
}
