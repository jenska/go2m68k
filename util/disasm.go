package util

import (
	"fmt"
	"log"

	"github.com/jenska/atari2go/cpu"
)

type disassembledOperand struct {
	operand string
	size    int
	memory  int
}

type disassembledOpcode struct {
	address     cpu.Address
	opcode      int
	instruction string
	op1         *disassembledOperand
	op2         *disassembledOperand
}

type disassemble func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode

var opcodes []disassemble

func Disassemble(bus cpu.AddressBus, start cpu.Address, size int) {
	end := start + cpu.Address(size)
	for start < end {
		if opcode, err := bus.Read(start, cpu.Word); err == nil {
			instruction := opcodes[opcode](bus, start)
			fmt.Println(instruction)
			start += cpu.Address(instruction.Size())
		} else {
			panic(fmt.Sprintf("invalid read %s", err))
		}
	}
}

func (opcode disassembledOpcode) Size() int {
	size := 2
	if opcode.op1 != nil {
		size += opcode.op1.size
	}
	if opcode.op2 != nil {
		size += opcode.op2.size
	}
	return size
}

func (opcode disassembledOpcode) String() string {
	op1hex, op2hex, opstr := "", "", ""
	if opcode.op1 != nil {
		switch opcode.op1.size {
		case 2:
			op1hex = fmt.Sprintf("%04x", opcode.op1.memory)
		case 4:
			op1hex = fmt.Sprintf("%08x", opcode.op1.memory)
		}
		opstr = opcode.op1.operand
	}
	if opcode.op2 != nil {
		switch opcode.op2.size {
		case 2:
			op2hex = fmt.Sprintf("%04x", opcode.op2.memory)
		case 4:
			op2hex = fmt.Sprintf("%08x", opcode.op2.memory)
		}
		opstr += ", " + opcode.op2.operand
	}
	return fmt.Sprintf("%08x %04x %8s %8s %s %s",
		opcode.address, opcode.opcode, op1hex, op2hex, opcode.instruction, opstr)
}

func newDisassembleOpcode(bus cpu.AddressBus, address cpu.Address, ins string) *disassembledOpcode {
	if o, err := bus.Read(address, cpu.Word); err == nil {
		return &disassembledOpcode{
			address:     address,
			opcode:      o,
			instruction: ins,
		}
	} else {
		panic("illegal address")
	}
}

func init() {
	log.Println("disasm init")

	moveq := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode {
		result := newDisassembleOpcode(bus, address, "moveq")
		opcode := result.opcode
		result.op1 = &disassembledOperand{operand: fmt.Sprintf("#%d", int8(opcode&0xff))}
		result.op2 = &disassembledOperand{operand: fmt.Sprintf("d%d", ((opcode >> 9) & 0x07))}
		return result
	}

	illegal := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode {
		return newDisassembleOpcode(bus, address, "????")
	}

	or := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode {
		return newDisassembleOpcode(bus, address, "or")
	}

	and := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	btst := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil } // movep
	bchg := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	bclr := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	bset := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	sub := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	add := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	eor := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	cmp := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	neg := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	chk := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	lea := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	not := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	clr := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	nbcd := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	swap := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	tst := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	tas := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	ext := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil } // movem
	movem := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	move := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	link := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil } // unlk, move
	jsr := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	jmp := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	addq := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	subq := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	dbcc := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil } // scc
	bsr := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }

	bcc := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode {
		return nil
	}

	divu := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	divs := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	mulu := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	muls := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	exg := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }  // and
	abcd := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil } // and
	xsr := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }  // ror, lsr, asr
	xsl := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }  // rol, lsl, asl
	linea := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }
	linef := func(bus cpu.AddressBus, address cpu.Address) *disassembledOpcode { return nil }

	opcodes = []disassemble{
		or, or, or, illegal, btst, bchg, bclr, bset, /* 0000 */
		and, and, and, illegal, btst, bchg, bclr, bset, /* 0200 */
		sub, sub, sub, illegal, btst, bchg, bclr, bset, /* 0400 */
		add, add, add, illegal, btst, bchg, bclr, bset, /* 0600 */
		btst, bchg, bclr, bset, btst, bchg, bclr, bset, /* 0800 */
		eor, eor, eor, illegal, btst, bchg, bclr, bset, /* 0A00 */
		cmp, cmp, cmp, illegal, btst, bchg, bclr, bset, /* 0C00 */
		move, move, move, illegal, btst, bchg, bclr, bset, /* 0E00 */
		move, illegal, move, move, move, move, move, move, /* 1000 */
		move, illegal, move, move, move, move, move, move, /* 1200 */
		move, illegal, move, move, move, move, move, move, /* 1400 */
		move, illegal, move, move, move, move, move, move, /* 1600 */
		move, illegal, move, move, move, move, move, move, /* 1800 */
		move, illegal, move, move, move, move, move, move, /* 1A00 */
		move, illegal, move, move, move, move, move, move, /* 1C00 */
		move, illegal, move, move, move, move, move, move, /* 1E00 */
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
		neg, neg, neg, move, illegal, illegal, chk, lea, /* 4000 */
		clr, clr, clr, move, illegal, illegal, chk, lea, /* 4200 */
		neg, neg, neg, move, illegal, illegal, chk, lea, /* 4400 */
		not, not, not, move, illegal, illegal, chk, lea, /* 4600 */
		nbcd, swap, ext, ext, illegal, illegal, chk, lea, /* 4800 */
		tst, tst, tst, tas, illegal, illegal, chk, lea, /* 4A00 */
		illegal, illegal, movem, movem, illegal, illegal, chk, lea, /* 4C00 */
		illegal, link, jsr, jmp, illegal, illegal, chk, lea, /* 4E00 */
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
		moveq, moveq, moveq, moveq, illegal, illegal, illegal, illegal, /* 7000 */
		moveq, moveq, moveq, moveq, illegal, illegal, illegal, illegal, /* 7200 */
		moveq, moveq, moveq, moveq, illegal, illegal, illegal, illegal, /* 7400 */
		moveq, moveq, moveq, moveq, illegal, illegal, illegal, illegal, /* 7600 */
		moveq, moveq, moveq, moveq, illegal, illegal, illegal, illegal, /* 7800 */
		moveq, moveq, moveq, moveq, illegal, illegal, illegal, illegal, /* 7A00 */
		moveq, moveq, moveq, moveq, illegal, illegal, illegal, illegal, /* 7C00 */
		moveq, moveq, moveq, moveq, illegal, illegal, illegal, illegal, /* 7E00 */
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
		xsr, xsr, xsr, illegal, xsl, xsl, xsl, illegal, /* E800 */
		xsr, xsr, xsr, illegal, xsl, xsl, xsl, illegal, /* EA00 */
		xsr, xsr, xsr, illegal, xsl, xsl, xsl, illegal, /* EC00 */
		xsr, xsr, xsr, illegal, xsl, xsl, xsl, illegal, /* EE00 */
		linef, linef, linef, linef, linef, linef, linef, linef, /* F000 */
		linef, linef, linef, linef, linef, linef, linef, linef,
		linef, linef, linef, linef, linef, linef, linef, linef,
		linef, linef, linef, linef, linef, linef, linef, linef,
		linef, linef, linef, linef, linef, linef, linef, linef,
		linef, linef, linef, linef, linef, linef, linef, linef,
		linef, linef, linef, linef, linef, linef, linef, linef,
		linef, linef, linef, linef, linef, linef, linef, linef,
	}
}
