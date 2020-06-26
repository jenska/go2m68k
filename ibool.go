package cpu

func init() {
	addOpcode("andi #imm, ccr", andi2ccr, 0x23c, 0xffff, 0x0000, "0:20", "7:14", "1:16", "234fc:12")
	addOpcode("andi #imm, sr", andi2sr, 0x27c, 0xffff, 0x0000, "0:20", "7:14", "1:16", "234fc:12")
	addOpcode("eori #imm, ccr", eori2ccr, 0xa3c, 0xffff, 0x0000, "0:20", "7:14", "1:16", "234fc:12")
	addOpcode("eori #imm, sr", eori2sr, 0xa7c, 0xffff, 0x0000, "0:20", "7:14", "1:16", "234fc:12")
	addOpcode("ori #imm, ccr", ori2ccr, 0x3c, 0xffff, 0x0000, "0:20", "7:14", "1:16", "234fc:12")
	addOpcode("ori #imm, sr", ori2sr, 0x7c, 0xffff, 0x0000, "0:20", "7:14", "1:16", "234fc:12")
}

func andi2ccr(c *M68K) {
	res := c.sr.ccr() & c.popPc(Byte)
	c.sr.setccr(res)
}

func andi2sr(c *M68K) {
	if c.sr.S {
		res := c.sr.bits() & c.popPc(Word)
		c.sr.setbits(res)
	} else {
		panic(NewError(PrivilegeViolationError, c, c.pc, nil))
	}
}

func eori2ccr(c *M68K) {
	res := c.sr.ccr() ^ c.popPc(Byte)
	c.sr.setccr(res)
}

func eori2sr(c *M68K) {
	if c.sr.S {
		res := c.sr.bits() ^ c.popPc(Word)
		c.sr.setbits(res)
	} else {
		panic(NewError(PrivilegeViolationError, c, c.pc, nil))
	}
}

func ori2ccr(c *M68K) {
	res := c.sr.ccr() | c.popPc(Byte)
	c.sr.setccr(res)
}

func ori2sr(c *M68K) {
	if c.sr.S {
		res := c.sr.bits() | c.popPc(Word)
		c.sr.setbits(res)
	} else {
		panic(NewError(PrivilegeViolationError, c, c.pc, nil))
	}
}
