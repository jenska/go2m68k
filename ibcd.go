package cpu

import (
	"runtime/debug"
)

func init() {
	addOpcode("abcd dx, dy", abcdRR, 0xc100, 0xf1f8, 0x0000, "01:6", "7:10", "234fc:4")
	addOpcode("abcd -(ax), -(ay)", abcdMM, 0xc108, 0xf1f8, 0x0000, "01:18", "7:31", "234fc:16")
	addOpcode("sbcd dx, dy", sbcdRR, 0x8100, 0xf1f8, 0x0000, "01:6", "7:10", "234fc:4")
	addOpcode("sbcd -(ax), -(ay)", sbcdMM, 0x8108, 0xf1f8, 0x0000, "01:18", "7:31", "234fc:16")
	addOpcode("nbcd", nbcd, 0x4800, 0xffc0, 0x0bf8, "01:8", "7:14", "234fc:6")
}

var bcdeaSrcPd = &eaPreDecrement{eaRegister{reg: ax}, 0}
var bcdeaDstPd = &eaPreDecrement{eaRegister{reg: ay}, 0}

func abcd(c *M68K, a, b int32) int32 {
	res := (a & 0x0f) + (b & 0x0f) + c.sr.x1()

	corf := int32(0)
	if res > 9 {
		corf = 6
	}

	res += (a & 0xf0) + (b & 0xf0)
	res += corf
	c.sr.X = res > 0x9f
	c.sr.C = c.sr.X
	if c.sr.X {
		res -= 0xa0
	}
	c.sr.setLogicalFlags(Byte, res)
	return res
}

func abcdRR(c *M68K) {
	dst := dx(c)
	res := abcd(c, *dy(c), *dst)
	Byte.set(res, dst)
}

//	c108 f1f8 abcd     b .          01:18 7:31 234fc:16
func abcdMM(c *M68K) {
	ea := bcdeaDstPd.init(c, Byte)
	src := bcdeaSrcPd.init(c, Byte).read()
	dst := ea.read()
	res := abcd(c, src, dst)
	ea.write(res)
}

func sbcd(c *M68K, s, d int32) int32 {
	lo := (d & 0xf) - (s & 0xf) - c.sr.x1()
	carry := int32(0)
	if lo < 0 {
		lo += 10
		carry = 1
	}
	hi := (((d & 0xf0) - (s & 0xf0)) >> 4) - carry
	carry = 0
	if hi < 0 {
		hi += 10
		carry = 1
	}
	res := (hi << 4) + lo
	c.sr.X = carry != 0
	c.sr.C = c.sr.X
	c.sr.setLogicalFlags(Byte, res)
	return res
}

func sbcdRR(c *M68K) {
	dst := dx(c)
	res := sbcd(c, *dy(c), *dst)
	Byte.set(res, dst)
}

//	c108 f1f8 abcd     b .          01:18 7:31 234fc:16
func sbcdMM(c *M68K) {
	ea := bcdeaDstPd.init(c, Byte)
	src := bcdeaSrcPd.init(c, Byte).read()
	dst := ea.read()
	res := sbcd(c, src, dst)
	ea.write(res)
}

func nbcd(c *M68K) {
	ea := c.resolveDstEA(Byte)
	dst := ea.read()
	if dst == 0xff {
		debug.PrintStack()
	}
	res := -dst - c.sr.x1()
	c.sr.X = res != 0
	if (res|dst)&0xf == 0 {
		res = (res & 0xf0) | 6
	}
	res += 0x9a
	ea.write(res)
	c.sr.setLogicalFlags(Byte, res)
}
