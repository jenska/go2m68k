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

/*
42c0 fff8 move     w .          1234fc:4
	DY() = MASK_OUT_BELOW_16(DY()) | m68ki_get_ccr()


42c0 ffc0 move     w A+-DXWL    1:8 234fc:4
	m68ki_write_16(M68KMAKE_GET_EA_AY_16, m68ki_get_ccr())


44c0 fff8 move     w .          01:12 7:10 234fc:4
	m68ki_set_ccr(DY())


44c0 ffc0 move     w A+-DXWLdxI 01:12 7:10 234fc:4
	m68ki_set_ccr(M68KMAKE_GET_OPER_AY_16)

	40c0 fff8 move     w .          0:6 7:7
	DY() = MASK_OUT_BELOW_16(DY()) | m68ki_get_sr()


40c0 fff8 move     w .          1:4p 234fc:8p
	if(m_s_flag)
		DY() = MASK_OUT_BELOW_16(DY()) | m68ki_get_sr()
	} else {
		m68ki_exception_privilege_violation()
	}
*/
