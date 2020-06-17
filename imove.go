package cpu

func init() {
	addOpcode("moveq #i, dy", moveq, 0x7000, 0xf100, 0x0000, "01:4", "7:7", "234fc:2")
}

func moveq(c *M68K) {
	res := int32(int8(c.ir & 0xff))
	c.sr.setLogicalFlags(Long, res)
	*dx(c) = res
}
