package cpu

func init() {
	addOpcode("moveq #i, dy", moveq, 0x7000, 0xf100, 0x0000, "01:4", "7:7", "234fc:2")
	addOpcode("movea.w ea, ax", moveaw, 0x3040, 0xf040, 0x0fff, "01:4", "7:7", "234fc:2")
	addOpcode("movea.w ea, ax", moveal, 0x2040, 0xf040, 0x0fff, "01:4", "7:7", "234fc:2")
}

func moveaw(c *M68K) {
	*ax(c) = int32(int16(c.resolveSrcEA((Word)).read()))
}
func moveal(c *M68K) {
	*ax(c) = c.resolveSrcEA(Long).read()
}

func moveq(c *M68K) {
	res := int32(int8(c.ir & 0xff))
	c.sr.setLogicalFlags(Long, res)
	*dx(c) = res
}
