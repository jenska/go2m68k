package m68k

import "log"

//asl, asr, lsl, lsr, rol, ror, roxl, roxr								-> shift/rotate instructions
//bchg, bclr, bset, btst												-> bit manipulation instructions
//clr, nbcd, neg, negx, not, scc, tas, tst								-> single operand instructions
//add, adda, and, cmp, cmpa, sub, suba, or, eor, mulu, muls, divu, divs -> standard instructions
//move, movea															-> move instructions
//addi, addq, andi, cmpi, eori, ori, moveq, subi, subq					-> immediate instructions
//bcc, bra, bsr, dbcc													-> specificational instructions
//jmp, jsr, lea, pea, movem												-> misc 1 instructions
//addx, cmpm, subx, abcd, sbcd											-> multiprecision instructions
//andiccr, andisr, chk, eorisr, eoriccr, orisr, oriccr, move from sr
//move to ccr, move to sr, exg, ext, link, moveusp, nop, reset
//rte, rtr, rts, stop, swap, trap, trapv, unlk							-> misc 2 instructions
//movep																	-> peripheral instruction

type instruction func(cpu *M68K, opcode uint16)

func (cpu *M68K) Execute() {
	if cpu.doubleFault {
		cpu.cpuHalted(cpu)
		cpu.sync(4)
	}

	defer func() {
		if r := recover(); r != nil {
			if x, ok := r.(group0exception); !ok {
				log.Fatalf("unable to recover from unexpeced error %d", x)
			}
		}
	}()

	cpu.PC += 2
	cpu.instructions[cpu.IRD](cpu, cpu.IRD)
}

func (cpu *M68K) init68000InstructionSet() {
	cpu.irqMode = AutoVectorInterrut

	cpu.instructions = make([]instruction, 1<<16)
	for i := range cpu.instructions {
		cpu.instructions[i] = illegal
	}

	// LSL, LSR, ASL, ASR reg
	for i := 0; i <= 7; i++ {
		for k := 0; k <= 2; k++ {
			for l := 0; l <= 1; l++ {
				for m := 0; m <= 7; m++ {
					instructions[(i << 9) + (0 << 8) + (k << 6) + (l << 5) + m + 0xe008] = &Core_68k::op_xsx<false, false, false>;
					instructions[(i << 9) + (1 << 8) + (k << 6) + (l << 5) + m + 0xe008] = &Core_68k::op_xsx<true, false, false>;

					instructions[(i << 9) + (0 << 8) + (k << 6) + (l << 5) + m + 0xe000] = &Core_68k::op_xsx<false, false, true>;
					instructions[(i << 9) + (1 << 8) + (k << 6) + (l << 5) + m + 0xe000] = &Core_68k::op_xsx<true, false, true>;
				}
			}
		}
	}

}

func illegal(cpu *M68K, opcode uint16) {
	switch {
	case opcode&0xa000 == 0xa000:
		cpu.raiseException(LineA)
	case opcode&0xf000 == 0xf000:
		cpu.raiseException(LineF)
	default:
		cpu.raiseException(IllegalOpcode)
	}
}

func rte(cpu *M68K, opcode uint16) {

}
