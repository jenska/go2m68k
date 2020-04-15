package cpu

func (cpu *M68K) InitISA68000() CPUBuilder {
	/*
		or := func() {}
		and := func() {}
		btst := func() {} // movep
		bchg := func() {}
		bclr := func() {}
		bset := func() {}
		sub := func() {}
		add := func() {}
		eor := func() {}
		cmp := func() {}
		neg := func() {}
		chk := func() {}
		lea := func() {}
		not := func() {}
		clr := func() {}
		nbcd := func() {}
		swap := func() {}
		tst := func() {}
		tas := func() {}
		ext := func() {} // movem
		movem := func() {}
		move := func() {}

		moveq := func() {
			cpu.D[(cpu.IR>>9)&0x7] = int32(cpu.IR & 0xff)
		}

		link := func() {} // unlk, move
		jsr := func() {}
		jmp := func() {}
		addq := func() {}
		subq := func() {}
		dbcc := func() {} // scc
		bsr := func() {}
		bcc := func() {}
		divu := func() {}
		divs := func() {}
		mulu := func() {}
		muls := func() {}
		exg := func() {}  // and
		abcd := func() {} // and
		xsr := func() {}  // ror, lsr, asr
		xsl := func() {}  // rol, lsl, asl
		linea := func() {}
		linef := func() {}

		isaTable := []func(){
			or, or, or, nil, btst, bchg, bclr, bset, // 0000
			and, and, and, nil, btst, bchg, bclr, bset, // 0200
			sub, sub, sub, nil, btst, bchg, bclr, bset, // 0400
			add, add, add, nil, btst, bchg, bclr, bset, // 0600
			btst, bchg, bclr, bset, btst, bchg, bclr, bset, // 0800
			eor, eor, eor, nil, btst, bchg, bclr, bset, // 0A00
			cmp, cmp, cmp, nil, btst, bchg, bclr, bset, // 0C00
			move, move, move, nil, btst, bchg, bclr, bset, // 0E00
			move, nil, move, move, move, move, move, move, // 1000
			move, nil, move, move, move, move, move, move, // 1200
			move, nil, move, move, move, move, move, move, // 1400
			move, nil, move, move, move, move, move, move, // 1600
			move, nil, move, move, move, move, move, move, // 1800
			move, nil, move, move, move, move, move, move, // 1A00
			move, nil, move, move, move, move, move, move, // 1C00
			move, nil, move, move, move, move, move, move, // 1E00
			move, move, move, move, move, move, move, move, // 2000
			move, move, move, move, move, move, move, move, // 2200
			move, move, move, move, move, move, move, move, // 2400
			move, move, move, move, move, move, move, move, // 2600
			move, move, move, move, move, move, move, move, // 2800
			move, move, move, move, move, move, move, move, // 2A00
			move, move, move, move, move, move, move, move, // 2C00
			move, move, move, move, move, move, move, move, // 2E00
			move, move, move, move, move, move, move, move, // 3000
			move, move, move, move, move, move, move, move, // 3200
			move, move, move, move, move, move, move, move, // 3400
			move, move, move, move, move, move, move, move, // 3600
			move, move, move, move, move, move, move, move, // 3800
			move, move, move, move, move, move, move, move, // 3A00
			move, move, move, move, move, move, move, move, // 3C00
			move, move, move, move, move, move, move, move, // 3E00
			neg, neg, neg, move, nil, nil, chk, lea, // 4000
			clr, clr, clr, move, nil, nil, chk, lea, // 4200
			neg, neg, neg, move, nil, nil, chk, lea, // 4400
			not, not, not, move, nil, nil, chk, lea, // 4600
			nbcd, swap, ext, ext, nil, nil, chk, lea, // 4800
			tst, tst, tst, tas, nil, nil, chk, lea, // 4A00
			nil, nil, movem, movem, nil, nil, chk, lea, // 4C00
			nil, link, jsr, jmp, nil, nil, chk, lea, // 4E00
			addq, addq, addq, dbcc, subq, subq, subq, dbcc, // 5000
			addq, addq, addq, dbcc, subq, subq, subq, dbcc, // 5200
			addq, addq, addq, dbcc, subq, subq, subq, dbcc, // 5400
			addq, addq, addq, dbcc, subq, subq, subq, dbcc, // 5600
			addq, addq, addq, dbcc, subq, subq, subq, dbcc, // 5800
			addq, addq, addq, dbcc, subq, subq, subq, dbcc, // 5A00
			addq, addq, addq, dbcc, subq, subq, subq, dbcc, // 5C00
			addq, addq, addq, dbcc, subq, subq, subq, dbcc, // 5E00
			bcc, bcc, bcc, bcc, bsr, bsr, bsr, bsr, // 6000
			bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, // 6200
			bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, // 6400
			bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, // 6600
			bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, // 6800
			bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, // 6A00
			bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, // 6C00
			bcc, bcc, bcc, bcc, bcc, bcc, bcc, bcc, // 6E00
			moveq, moveq, moveq, moveq, nil, nil, nil, nil, // 7000
			moveq, moveq, moveq, moveq, nil, nil, nil, nil, // 7200
			moveq, moveq, moveq, moveq, nil, nil, nil, nil, // 7400
			moveq, moveq, moveq, moveq, nil, nil, nil, nil, // 7600
			moveq, moveq, moveq, moveq, nil, nil, nil, nil, // 7800
			moveq, moveq, moveq, moveq, nil, nil, nil, nil, // 7A00
			moveq, moveq, moveq, moveq, nil, nil, nil, nil, // 7C00
			moveq, moveq, moveq, moveq, nil, nil, nil, nil, // 7E00
			or, or, or, divu, or, or, or, divs, // 8000
			or, or, or, divu, or, or, or, divs, // 8200
			or, or, or, divu, or, or, or, divs, // 8400
			or, or, or, divu, or, or, or, divs, // 8600
			or, or, or, divu, or, or, or, divs, // 8800
			or, or, or, divu, or, or, or, divs, // 8A00
			or, or, or, divu, or, or, or, divs, // 8C00
			or, or, or, divu, or, or, or, divs, // 8E00
			sub, sub, sub, sub, sub, sub, sub, sub, // 9000
			sub, sub, sub, sub, sub, sub, sub, sub, // 9200
			sub, sub, sub, sub, sub, sub, sub, sub, // 9400
			sub, sub, sub, sub, sub, sub, sub, sub, // 9600
			sub, sub, sub, sub, sub, sub, sub, sub, // 9800
			sub, sub, sub, sub, sub, sub, sub, sub, // 9A00
			sub, sub, sub, sub, sub, sub, sub, sub, // 9C00
			sub, sub, sub, sub, sub, sub, sub, sub, // 9E00
			linea, linea, linea, linea, linea, linea, linea, linea, // A000
			linea, linea, linea, linea, linea, linea, linea, linea,
			linea, linea, linea, linea, linea, linea, linea, linea,
			linea, linea, linea, linea, linea, linea, linea, linea,
			linea, linea, linea, linea, linea, linea, linea, linea,
			linea, linea, linea, linea, linea, linea, linea, linea,
			linea, linea, linea, linea, linea, linea, linea, linea,
			linea, linea, linea, linea, linea, linea, linea, linea,
			cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, // B000
			cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, // B200
			cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, // B400
			cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, // B600
			cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, // B800
			cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, // BA00
			cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, // BC00
			cmp, cmp, cmp, cmp, cmp, cmp, cmp, cmp, // BE00
			and, and, and, mulu, abcd, exg, exg, muls, // C000
			and, and, and, mulu, abcd, exg, exg, muls, // C200
			and, and, and, mulu, abcd, exg, exg, muls, // C400
			and, and, and, mulu, abcd, exg, exg, muls, // C600
			and, and, and, mulu, abcd, exg, exg, muls, // C800
			and, and, and, mulu, abcd, exg, exg, muls, // CA00
			and, and, and, mulu, abcd, exg, exg, muls, // CC00
			and, and, and, mulu, abcd, exg, exg, muls, // CE00
			add, add, add, add, add, add, add, add, // D000
			add, add, add, add, add, add, add, add, // D200
			add, add, add, add, add, add, add, add, // D400
			add, add, add, add, add, add, add, add, // D600
			add, add, add, add, add, add, add, add, // D800
			add, add, add, add, add, add, add, add, // DA00
			add, add, add, add, add, add, add, add, // DC00
			add, add, add, add, add, add, add, add, // DE00
			xsr, xsr, xsr, xsr, xsl, xsl, xsl, xsl, // E000
			xsr, xsr, xsr, xsr, xsl, xsl, xsl, xsl, // E200
			xsr, xsr, xsr, xsr, xsl, xsl, xsl, xsl, // E400
			xsr, xsr, xsr, xsr, xsl, xsl, xsl, xsl, // E600
			xsr, xsr, xsr, nil, xsl, xsl, xsl, nil, // E800
			xsr, xsr, xsr, nil, xsl, xsl, xsl, nil, // EA00
			xsr, xsr, xsr, nil, xsl, xsl, xsl, nil, // EC00
			xsr, xsr, xsr, nil, xsl, xsl, xsl, nil, // EE00
			linef, linef, linef, linef, linef, linef, linef, linef, // F000
			linef, linef, linef, linef, linef, linef, linef, linef,
			linef, linef, linef, linef, linef, linef, linef, linef,
			linef, linef, linef, linef, linef, linef, linef, linef,
			linef, linef, linef, linef, linef, linef, linef, linef,
			linef, linef, linef, linef, linef, linef, linef, linef,
			linef, linef, linef, linef, linef, linef, linef, linef,
			linef, linef, linef, linef, linef, linef, linef, linef}
	*/
	return cpu
}
