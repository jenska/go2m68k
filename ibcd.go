package cpu

func init() {
	addOpcode("abcd dx, dy", abcdRR, 0xc100, 0xf1f8, 0x0000, "01:6", "7:10", "234fc:4")
	addOpcode("abcd -(ax), -(ay)", abcdMM, 0xc108, 0xf1f8, 0x0000, "01:18", "7:31", "234fc:16")
	addOpcode("nbcd", nbcd, 0x4800, 0xffc0, 0x0, "01:8", "7:14", "234fc:6")
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
	if c.sr.X {
		res -= 0xa0
	}
	c.sr.setLogicalFlags(Byte, res)
	return res
}

func abcdRR(c *M68K) {
	dst := dx(c)
	src := dy(c)
	res := abcd(c, *src, *dst)
	Byte.set(res, dst)
}

//	c108 f1f8 abcd     b .          01:18 7:31 234fc:16
func abcdMM(c *M68K) {
	src := bcdeaDstPd.init(c, Byte).read()
	ea := bcdeaSrcPd.init(c, Byte)
	dst := ea.read()
	res := abcd(c, src, dst)
	ea.write(res)
}

func nbcd(c *M68K) {
	dst := dy(c)
	res := -*dst - c.sr.x1()
	if res != 0 {
		c.sr.V = true // undefined
		if (res|*dst)&0xf == 0 {
			res = (res & 0xf0) | 6
		}
		res += 0x9a
		c.sr.X = true
	} else {
		c.sr.X = false
	}
	c.sr.setLogicalFlags(Byte, res)
}
