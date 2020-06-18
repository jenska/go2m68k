package cpu

func init() {
	addOpcode("bcc", bcc, 0x6000, 0xf000, 0x0000, "01:10", "7:13", "234fc:6")
	addOpcode("dbcc", dbcc, 0x50c8, 0xf0f8, 0x0000, "0:12", "7:14", "1:10", "234fc:6")
	addOpcode("rts", rts, 0x4e75, 0xffff, 0x0000, "01:16", "7:15", "234fc:10")
	addOpcode("rtr", rtr, 0x4e77, 0xffff, 0x0000, "01:20", "7:22", "234fc:14")
	addOpcode("stop #sr", stop, 0x4e72, 0xffff, 0x0000, "01:4", "7:13", "234fc:8")
}

func bcc(c *M68K) {
	cc := (c.ir >> 8) & 0xf
	dis := int32(c.ir & 0xff)
	if cc == 1 { // bsr
		if dis == 0 {
			dis = int32(int16(c.popPc(Word)))
		} else {
			dis = int32(int8(dis))
		}
		c.push(Long, c.pc)
		c.pc += dis
		// c.cycles +=
	} else if c.sr.testCC(cc) {
		if dis == 0 {
			dis = int32(int16(c.popPc(Word)))
		} else {
			dis = int32(int8(dis))
		}
		c.pc += dis
		// c.cycles +=
	} else {
		// c.cycles +=
		c.pc += Word.size
	}
}

func dbcc(c *M68K) {
	if c.sr.testCC(c.ir >> 8) {
		c.pc += Word.size // skip displacement value
		// c.cycles +=
	} else {
		count := int32(int16(*dy(c))) - 1
		Word.set(count, dy(c))
		dis := int32(int16(c.popPc(Word)))
		if count == -1 {
			// c.cycles +=
		} else {
			// c.cycles +=
			c.pc += dis - Word.size
		}

	}
}

func rts(c *M68K) {
	c.pc = c.pop(Long)
}

func rtr(c *M68K) {
	c.sr.setccr(c.pop(Word))
	c.pc = c.pop(Long)
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
