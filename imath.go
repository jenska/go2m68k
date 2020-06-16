package cpu

func init() {
	addOpcode("add ea,reg", add2reg, 0xd000, 0xf100, 0x0fff, "071234fc:4")
	addOpcode("add reg,ea", add2ea, 0xd100, 0xf100, 0x0fff, "071234fc:4")
	addOpcode("adda ea, ax", adda, 0xd0c0, 0xf1c0, 0x0fff, "071234fc:4")
	addOpcode("sub ea,reg", sub2reg, 0x9000, 0xf100, 0x0fff, "071234fc:4")
	addOpcode("sub reg,ea", sub2ea, 0x9100, 0xf100, 0x0fff, "071234fc:4")
	addOpcode("suba ea, ax", suba, 0x90c0, 0xf1c0, 0x0fff, "071234fc:4")
}

func add2ea(c *M68K) {
	size := operands[(c.ir>>6)&0x3]
	src := dx(c)
	ea := c.resolveDstEA(size)
	res := *src + ea.read()
	c.sr.setAddSubFlags(size, *src, ea.read(), res)
	ea.write(res)
}

func add2reg(c *M68K) {
	size := operands[(c.ir>>6)&0x3]
	dst := dx(c)
	ea := c.resolveDstEA(size)
	res := *dst + ea.read()
	c.sr.setAddSubFlags(size, ea.read(), *dst, res)
	size.set(res, dst)
}

func adda(c *M68K) {
	size := operands[((c.ir>>9)&0x1)+1]
	dst := ax(c)
	ea := c.resolveDstEA(size)
	res := *dst + ea.read()
	size.set(res, dst)
}

func sub2ea(c *M68K) {
	size := operands[(c.ir>>6)&0x3]
	src := ax(c)
	ea := c.resolveDstEA(size)
	res := ea.read() - *src
	c.sr.setAddSubFlags(size, *src, ea.read(), res)
	ea.write(res)
}

func sub2reg(c *M68K) {
	size := operands[(c.ir>>6)&0x3]
	dst := ax(c)
	ea := c.resolveDstEA(size)
	res := *dst - ea.read()
	c.sr.setAddSubFlags(size, ea.read(), *dst, res)
	size.set(res, dst)
}

func suba(c *M68K) {
	size := operands[((c.ir>>9)&0x1)+1]
	dst := ax(c)
	ea := c.resolveDstEA(size)
	res := *dst - ea.read()
	size.set(res, dst)
}
