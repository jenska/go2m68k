package cpu

func init() {
	addOpcode("bcc", bcc, 0x6000, 0xf000, 0x0000, "01:10", "7:13", "234fc:6")
	addOpcode("dbcc", dbcc, 0x50c8, 0xf0f8, 0x0000, "0:12", "7:14", "1:10", "234fc:6")
}

func bcc(c *M68K) {
	dis := int32(int8(c.ir & 0xff))
	if dis == 0 {
		dis = int32(int16(c.popPc(Word)))
		// c.cycles +=
	}
	cc := (c.ir >> 8) & 0xf
	if cc == 1 { // bsr
		c.push(Long, c.pc)
		c.pc += dis
		// c.cycles +=
	} else if c.sr.testCC(cc) {
		c.pc += dis
		// c.cycles +=
	} else {
		// c.cycles +=
	}
}

func dbcc(c *M68K) {
	if c.sr.testCC((c.ir >> 8) & 0xf) {
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
