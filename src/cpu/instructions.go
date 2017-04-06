package cpu

type Instruction func(cpu *M68k) int

func (cpu *M68k) init68000InstructionSet() {
	cpu.instructions = []Instruction{
		or, or, or, nil, btst, bchg, bclr, bset, /* 0000 */
		and, and, and, nil, btst, bchg, bclr, bset, /* 0200 */
		sub, sub, sub, nil, btst, bchg, bclr, bset, /* 0400 */
		add, add, add, nil, btst, bchg, bclr, bset, /* 0600 */
		btst, bchg, bclr, bset, btst, bchg, bclr, bset, /* 0800 */
		eor, eor, eor, nil, btst, bchg, bclr, bset, /* 0A00 */
		cmp, cmp, cmp, nil, btst, bchg, bclr, bset, /* 0C00 */
		move, move, move, nil, btst, bchg, bclr, bset, /* 0E00 */
		move, nil, move, move, move, move, move, move, /* 1000 */
		move, nil, move, move, move, move, move, move, /* 1200 */
		move, nil, move, move, move, move, move, move, /* 1400 */
		move, nil, move, move, move, move, move, move, /* 1600 */
		move, nil, move, move, move, move, move, move, /* 1800 */
		move, nil, move, move, move, move, move, move, /* 1A00 */
		move, nil, move, move, move, move, move, move, /* 1C00 */
		move, nil, move, move, move, move, move, move, /* 1E00 */
		move, move, move, move, move, move, move, move, /* 2000 */
		move, move, move, move, move, move, move, move, /* 2200 */
		move, move, move, move, move, move, move, move, /* 2400 */
		move, move, move, move, move, move, move, move, /* 2600 */
		move, move, move, move, move, move, move, move, /* 2800 */
		move, move, move, move, move, move, move, move, /* 2A00 */
		move, move, move, move, move, move, move, move, /* 2C00 */
		move, move, move, move, move, move, move, move, /* 2E00 */
		move, move, move, move, move, move, move, move, /* 3000 */
		move, move, move, move, move, move, move, move, /* 3200 */
		move, move, move, move, move, move, move, move, /* 3400 */
		move, move, move, move, move, move, move, move, /* 3600 */
		move, move, move, move, move, move, move, move, /* 3800 */
		move, move, move, move, move, move, move, move, /* 3A00 */
		move, move, move, move, move, move, move, move, /* 3C00 */
		move, move, move, move, move, move, move, move, /* 3E00 */
		neg, neg, neg, move, nil, nil, chk, lea, /* 4000 */
		clr, clr, clr, move, nil, nil, chk, lea, /* 4200 */
		neg, neg, neg, move, nil, nil, chk, lea, /* 4400 */
		not, not, not, move, nil, nil, chk, lea, /* 4600 */
		nbcd, swap, ext, ext, nil, nil, chk, lea, /* 4800 */
		tst, tst, tst, tas, nil, nil, chk, lea, /* 4A00 */
		nil, nil, movem, movem, nil, nil, chk, lea, /* 4C00 */
		nil, link, jsr, jmp, nil, nil, chk, lea, /* 4E00 */
		addq, addq, addq, dbcc, subq, subq, subq, dbcc, /* 5000 */
		addq, addq, addq, dbcc, subq, subq, subq, dbcc, /* 5200 */
		addq, addq, addq, dbcc, subq, subq, subq, dbcc, /* 5400 */
		addq, addq, addq, dbcc, subq, subq, subq, dbcc, /* 5600 */
		addq, addq, addq, dbcc, subq, subq, subq, dbcc, /* 5800 */
		addq, addq, addq, dbcc, subq, subq, subq, dbcc, /* 5A00 */
		addq, addq, addq, dbcc, subq, subq, subq, dbcc, /* 5C00 */
		addq, addq, addq, dbcc, subq, subq, subq, dbcc, /* 5E00 */
		bcc, bcc, bcc, bcc, bsr, bsr, bsr, bsr, /* 6000 */
		bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, /* 6200 */
		bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, /* 6400 */
		bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, /* 6600 */
		bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, /* 6800 */
		bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, /* 6A00 */
		bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, /* 6C00 */
		bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, /* 6E00 */
		moveq, moveq, moveq, moveq, nil, nil, nil, nil, /* 7000 */
		moveq, moveq, moveq, moveq, nil, nil, nil, nil, /* 7200 */
		moveq, moveq, moveq, moveq, nil, nil, nil, nil, /* 7400 */
		moveq, moveq, moveq, moveq, nil, nil, nil, nil, /* 7600 */
		moveq, moveq, moveq, moveq, nil, nil, nil, nil, /* 7800 */
		moveq, moveq, moveq, moveq, nil, nil, nil, nil, /* 7A00 */
		moveq, moveq, moveq, moveq, nil, nil, nil, nil, /* 7C00 */
		moveq, moveq, moveq, moveq, nil, nil, nil, nil, /* 7E00 */
		or, or, or, divu, or, or, or, divs, /* 8000 */
		or, or, or, divu, or, or, or, divs, /* 8200 */
		or, or, or, divu, or, or, or, divs, /* 8400 */
		or, or, or, divu, or, or, or, divs, /* 8600 */
		or, or, or, divu, or, or, or, divs, /* 8800 */
		or, or, or, divu, or, or, or, divs, /* 8A00 */
		or, or, or, divu, or, or, or, divs, /* 8C00 */
		or, or, or, divu, or, or, or, divs, /* 8E00 */
		sub, sub, sub, sub, sub, sub, sub, sub, /* 9000 */
		sub, sub, sub, sub, sub, sub, sub, sub, /* 9200 */
		sub, sub, sub, sub, sub, sub, sub, sub, /* 9400 */
		sub, sub, sub, sub, sub, sub, sub, sub, /* 9600 */
		sub, sub, sub, sub, sub, sub, sub, sub, /* 9800 */
		sub, sub, sub, sub, sub, sub, sub, sub, /* 9A00 */
		sub, sub, sub, sub, sub, sub, sub, sub, /* 9C00 */
		sub, sub, sub, sub, sub, sub, sub, sub, /* 9E00 */
		linea, linea, linea, linea, linea, linea, linea, linea, /* A000 */
		linea, linea, linea, linea, linea, linea, linea, linea,
		linea, linea, linea, linea, linea, linea, linea, linea,
		linea, linea, linea, linea, linea, linea, linea, linea,
		linea, linea, linea, linea, linea, linea, linea, linea,
		linea, linea, linea, linea, linea, linea, linea, linea,
		linea, linea, linea, linea, linea, linea, linea, linea,
		linea, linea, linea, linea, linea, linea, linea, linea,
		cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, /* B000 */
		cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, /* B200 */
		cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, /* B400 */
		cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, /* B600 */
		cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, /* B800 */
		cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, /* BA00 */
		cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, /* BC00 */
		cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, /* BE00 */
		and, and, and, mulu, abcd, exg, exg, muls, /* C000 */
		and, and, and, mulu, abcd, exg, exg, muls, /* C200 */
		and, and, and, mulu, abcd, exg, exg, muls, /* C400 */
		and, and, and, mulu, abcd, exg, exg, muls, /* C600 */
		and, and, and, mulu, abcd, exg, exg, muls, /* C800 */
		and, and, and, mulu, abcd, exg, exg, muls, /* CA00 */
		and, and, and, mulu, abcd, exg, exg, muls, /* CC00 */
		and, and, and, mulu, abcd, exg, exg, muls, /* CE00 */
		add, add, add, add, add, add, add, add, /* D000 */
		add, add, add, add, add, add, add, add, /* D200 */
		add, add, add, add, add, add, add, add, /* D400 */
		add, add, add, add, add, add, add, add, /* D600 */
		add, add, add, add, add, add, add, add, /* D800 */
		add, add, add, add, add, add, add, add, /* DA00 */
		add, add, add, add, add, add, add, add, /* DC00 */
		add, add, add, add, add, add, add, add, /* DE00 */
		xsr, xsr, xsr, xsr, xsl, xsl, xsl, xsl, /* E000 */
		xsr, xsr, xsr, xsr, xsl, xsl, xsl, xsl, /* E200 */
		xsr, xsr, xsr, xsr, xsl, xsl, xsl, xsl, /* E400 */
		xsr, xsr, xsr, xsr, xsl, xsl, xsl, xsl, /* E600 */
		xsr, xsr, xsr, nil, xsl, xsl, xsl, nil, /* E800 */
		xsr, xsr, xsr, nil, xsl, xsl, xsl, nil, /* EA00 */
		xsr, xsr, xsr, nil, xsl, xsl, xsl, nil, /* EC00 */
		xsr, xsr, xsr, nil, xsl, xsl, xsl, nil, /* EE00 */
		linef, linef, linef, linef, linef, linef, linef, linef, /* F000 */
		linef, linef, linef, linef, linef, linef, linef, linef,
		linef, linef, linef, linef, linef, linef, linef, linef,
		linef, linef, linef, linef, linef, linef, linef, linef,
		linef, linef, linef, linef, linef, linef, linef, linef,
		linef, linef, linef, linef, linef, linef, linef, linef,
		linef, linef, linef, linef, linef, linef, linef, linef,
		linef, linef, linef, linef, linef, linef, linef, linef}
}

func or(cpu *M68k) int    { return 0 }
func and(cpu *M68k) int   { return 0 }
func btst(cpu *M68k) int  { return 0 } // movep
func bchg(cpu *M68k) int  { return 0 }
func bclr(cpu *M68k) int  { return 0 }
func bset(cpu *M68k) int  { return 0 }
func sub(cpu *M68k) int   { return 0 }
func add(cpu *M68k) int   { return 0 }
func eor(cpu *M68k) int   { return 0 }
func cmp(cpu *M68k) int   { return 0 }
func neg(cpu *M68k) int   { return 0 }
func chk(cpu *M68k) int   { return 0 }
func lea(cpu *M68k) int   { return 0 }
func not(cpu *M68k) int   { return 0 }
func clr(cpu *M68k) int   { return 0 }
func nbcd(cpu *M68k) int  { return 0 }
func swap(cpu *M68k) int  { return 0 }
func tst(cpu *M68k) int   { return 0 }
func tas(cpu *M68k) int   { return 0 }
func ext(cpu *M68k) int   { return 0 } // movem
func movem(cpu *M68k) int { return 0 }
func link(cpu *M68k) int  { return 0 } // unlk, move
func jsr(cpu *M68k) int   { return 0 }
func jmp(cpu *M68k) int   { return 0 }
func addq(cpu *M68k) int  { return 0 }
func subq(cpu *M68k) int  { return 0 }
func dbcc(cpu *M68k) int  { return 0 } // scc
func bsr(cpu *M68k) int   { return 0 }
func bcc(cpu *M68k) int   { return 0 }
func divu(cpu *M68k) int  { return 0 }
func divs(cpu *M68k) int  { return 0 }
func mulu(cpu *M68k) int  { return 0 }
func muls(cpu *M68k) int  { return 0 }
func exg(cpu *M68k) int   { return 0 } // and
func abcd(cpu *M68k) int  { return 0 } // and
func xsr(cpu *M68k) int  { return 0 } // ror, lsr, asr
func xsl(cpu *M68k) int  { return 0 } // rol, lsl, asl

func linea(cpu *M68k) int {
	return cpu.RaiseException(XPT_LNA)
}

func linef(cpu *M68k) int {
	return cpu.RaiseException(XPT_LNF)
}
