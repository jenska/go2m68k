package cpu

import (
	"fmt"
	"log"
)

type (
	dasmOperand struct {
		operand string
		size    *Size
	}

	// DasmInstruction is a container for a disassembled M68K instruction
	DasmInstruction struct {
		instruction string
		address     int32
		operands    []dasmOperand
	}

	// DasmIterator iterates over DasmInstruction
	DasmIterator struct {
		bus     AddressBus
		address int32
	}

	dasm func(ir uint16, pc int32, bus AddressBus) DasmInstruction
)

var dasmTable = make([]dasm, 0x10000)

// Disassembler returns an M68K DasmIterator
func Disassembler(pc int32, bus AddressBus) *DasmIterator {
	return &DasmIterator{bus, pc}
}

// Next disassembled instruction
func (iter *DasmIterator) Next() DasmInstruction {
	ir := uint16(iter.bus.read(iter.address, Word))
	result := dasmTable[ir](ir, iter.address, iter.bus)

	offset := int32(2)
	for _, op := range result.operands {
		if operand := op.size; operand != nil {
			offset += operand.size
		}
	}
	iter.address += offset

	return result
}

func dasmInstruction(name string, address int32, ops ...dasmOperand) DasmInstruction {
	return DasmInstruction{
		instruction: name,
		address:     address,
		operands:    ops,
	}
}

// Instruction as string
func (di DasmInstruction) Instruction() string {
	return di.instruction
}

// Operand1 as string if exists, otherwise empty string
func (di DasmInstruction) Operand1() string {
	if len(di.operands) >= 1 {
		return di.operands[0].operand
	}
	return ""
}

// Operand2 as string if exists, otherwise empty string
func (di DasmInstruction) Operand2() string {
	if len(di.operands) == 2 {
		return di.operands[1].operand
	}
	return ""
}

// Short formatted string
func (di DasmInstruction) String() string {
	result := fmt.Sprintf("%08x %-10s ", di.address, di.Instruction())
	if op := di.Operand1(); op != "" {
		result += op
	}
	if op := di.Operand2(); op != "" {
		result += ", " + op
	}

	return result
}

//------------------------------------------------------------------------------

func init() {
	counter := 0
	for i := range dasmTable {
		opcode := uint16(i)
		dasmTable[i] = dasmIllegal
		for _, info := range opcodeTable {
			if (opcode & info.mask) == info.match {
				if info.move && !validEA(eam(uint16(opcode)), 0xbf8) {
					continue
				}
				if validEA(opcode, info.eaMask) {
					// log.Printf("%04x %+v\n", opcode, info)
					dasmTable[i] = info.dasm
					counter++
					break
				}
			}
		}
	}
	log.Printf("added %d disassembler instructions", counter)
}

//------------------------------------------------------------------------------

func dasmIllegal(ir uint16, pc int32, bus AddressBus) DasmInstruction {
	return dasmInstruction("dc.w", pc, dasmOperand{Word.HexString(int32(ir)), nil})
}

func dasmBra8(ir uint16, pc int32, bus AddressBus) DasmInstruction {
	target := pc + int32(int8(ir)) + 2
	return dasmInstruction("bra.s", pc, dasmOperand{Long.HexString(target), nil})
}

func dasmDbra(ir uint16, pc int32, bus AddressBus) DasmInstruction {
	target := pc + int32(int8(ir)) + 2
	return dasmInstruction("bra.s", pc, dasmOperand{Long.HexString(target), nil})
}
func dasmStop(ir uint16, pc int32, bus AddressBus) DasmInstruction {
	target := pc + int32(int8(ir)) + 2
	return dasmInstruction("bra.s", pc, dasmOperand{Long.HexString(target), nil})
}

func dasmMoveq(ir uint16, pc int32, bus AddressBus) DasmInstruction {
	op1 := dasmOperand{fmt.Sprintf("d[%v]", (ir>>9)&7), nil}
	op2 := dasmOperand{fmt.Sprintf("#%s", Long.SignedHexString(int32(int8(ir)))), nil}
	return dasmInstruction("moveq", pc, op1, op2)
}

/*
func dasmBra16(ir uint16, pc uint32, bus AddressBus) DasmInstruction {
	return "bra", fmt.Sprintf("$%08x", int(d.pc)+Word.signed(d.pop(Word)))
}

func d68000_move_8(d *dasm) (string, string) {
	return "move.b", ""
}
func d68000_move_16(d *dasm) (string, string) {
	return "move", ""
}
func d68000_move_32(d *dasm) (string, string) {
	return "move.l", ""
}

func d68000_movea_16(d *dasm) (string, string) {
	return "movea", ""
}

func d68000_movea_32(d *dasm) (string, string) {
	return "movea.l", ""
}

func d68000_move_to_sr(d *dasm) (string, string) {
	return "move", dasmEAMode(Word, d) + ", sr"
}

/* ======================================================================== */
/* ============================= BitOp Helpers ============================ */
/* ======================================================================== */
