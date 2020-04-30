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

	DasmIterator struct {
		bus     AddressBus
		address int32
	}

	dasm func(ir uint16, pc int32, bus AddressBus) DasmInstruction
)

var dasmTable = make([]dasm, 0x10000)

// Disassemble an M68K instruction
func Disassembler(pc int32, bus AddressBus) *DasmIterator {
	return &DasmIterator{bus, pc}
}

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
	} else {
		return ""
	}
}

// Operand2 as string if exists, otherwise empty string
func (di DasmInstruction) Operand2() string {
	if len(di.operands) == 2 {
		return di.operands[1].operand
	} else {
		return ""
	}
}

// Short formatted string
func (di DasmInstruction) String() string {
	result := fmt.Sprintf("%08x %-10s ", di.address, di.Instruction())
	if op1 := di.Operand1(); op1 != "" {
		result += op1
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
					dasmTable[i] = info.dasm
					counter++
					break
				}
			}
		}
	}
	log.Printf("added %d disassembler instructions", counter)
}

/*
func dasmEAMode(s *Size, d *core) string {
	switch ea := d.ir & 0x3f; ea {
	case 0, 1, 2, 3, 4, 5, 6, 7:
		return fmt.Sprintf("D%d", d.ir&7)
	case 8, 9, 10, 11, 12, 13, 14, 15:
		return fmt.Sprintf("A%d", d.ir&7)
	case 16, 17, 18, 19, 20, 21, 22, 23:
		return fmt.Sprintf("(A%d)", d.ir&7)
	case 24, 25, 26, 27, 28, 29, 30, 31:
		return fmt.Sprintf("(A%d)+", d.ir&7)
	case 32, 33, 34, 35, 36, 37, 38, 39:
		return fmt.Sprintf("-(A%d)", d.ir&7)

	case 0x3c:
		return "#" + s.SignedHexString(d.pop(s))
	}

	return "xxx not implemented xxx"
}

func d68000_illegal(d *core) (string, string) {
	return "dc.w", fmt.Sprintf("$%04x ; illegal", d.ir)
}
*/

//------------------------------------------------------------------------------

func dasmIllegal(ir uint16, pc int32, bus AddressBus) DasmInstruction {
	return dasmInstruction("dc.w", pc, dasmOperand{Word.HexString(int32(ir)), nil})
}

func dasmBra8(ir uint16, pc int32, bus AddressBus) DasmInstruction {
	target := pc + int32(int8(ir)) + 2
	return dasmInstruction("bra.s", pc, dasmOperand{Long.HexString(target), nil})
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
