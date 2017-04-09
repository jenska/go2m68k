package m68k

import (
	"fmt"
)

type disassemble func(handler AddressHandler, address uint32) *disassembledInstruction

var dInstructionTable = make([]disassemble, 0x10000)

func disassembler(handler AddressHandler, address uint32) *disassembledInstruction {
	opcode := dOpcode(handler, address)
	if disassemble := dInstructionTable[opcode]; disassemble != nil {
		return disassemble(handler, address)
	}
	return &disassembledInstruction{"data.w", opcode, address, []disassembledEA{
		disassembledEA{fmt.Sprintf("$%04x  ; (%d) unknown instruction", opcode, opcode), nil, opcode}}}
}

type disassembledInstruction struct {
	name    string
	opcode  uint32
	address uint32
	ea      []disassembledEA
}

func (i *disassembledInstruction) size() int {
	size := Word.Size
	for _, a := range i.ea {
		if o := a.o; o != nil {
			size += o.AlignedSize
		}
	}
	return int(size)
}

func (i *disassembledInstruction) shortFormat() string {
	result := fmt.Sprintf("%08x %-9s ", i.address, i.name)
	if len(i.ea) > 0 {
		result += i.ea[0].name
	}
	if len(i.ea) > 1 {
		result += ", " + i.ea[1].name
	}
	return result
}

func (i *disassembledInstruction) detailedFormat() string {
	result := ""
	if len(i.ea) > 0 {
		result = i.ea[0].toHex()
	}
	if len(i.ea) > 1 {
		if len(result) > 0 {
			result += " "
		}
		result += i.ea[1].toHex()
	}

	result = fmt.Sprintf("%08x %04x %-17s %-9s ", i.address, i.opcode, result, i.name)
	if len(i.ea) > 0 {
		result += i.ea[0].name
	}
	if len(i.ea) > 1 {
		result += ", " + i.ea[1].name
	}
	return result
}

func (i *disassembledInstruction) String() string {
	return i.detailedFormat()
}

type disassembledEA struct {
	name   string
	o      *operand
	memory uint32
}

type eaDisassembler interface {
	disassemble(handler *AddressHandler, address uint32) disassembledEA
}

func (ea *disassembledEA) toHex() string {
	if ea.o != nil {
		return fmt.Sprintf(ea.o.formatter, ea.memory)
	}
	return ""
}

func (ea *eaDataRegister) disassemble(_ AddressHandler, _ uint32) disassembledEA {
	return disassembledEA{fmt.Sprintf("d%d", ea.register), ea.o, 0}
}

func (ea *eaAddressRegister) disassemble(_ AddressHandler, _ uint32) disassembledEA {
	return disassembledEA{fmt.Sprintf("a%d", ea.register), ea.o, 0}
}

func (ea *eaAddressRegisterIndirect) disassemble(_ AddressHandler, _ uint32) disassembledEA {
	return disassembledEA{fmt.Sprintf("(a%d)", ea.register), ea.o, 0}
}

func (ea *eaAddressRegisterPreDec) disassemble(_ AddressHandler, _ uint32) disassembledEA {
	return disassembledEA{fmt.Sprintf("-(a%d)", ea.register), ea.o, 0}
}

func (ea *eaAddressRegisterPostInc) disassemble(_ AddressHandler, _ uint32) disassembledEA {
	return disassembledEA{fmt.Sprintf("(a%d)+", ea.register), ea.o, 0}
}

func (ea *eaAddressRegisterWithDisplacement) disassemble(handler AddressHandler, address uint32) disassembledEA {
	mem := dOpcode(handler, address)
	displacement := int(uint16(mem))
	return disassembledEA{fmt.Sprintf("$%04x(a%d)", displacement, ea.register), ea.o, mem}
}

func (ea *eaPCWithDisplacement) disassemble(handler AddressHandler, address uint32) disassembledEA {
	mem := dOpcode(handler, address)
	displacement := int(uint16(mem))
	return disassembledEA{fmt.Sprintf("$%04x(pc)", displacement), ea.o, mem}
}

func (ea *eaAddressRegisterWithIndex) disassemble(handler AddressHandler, address uint32) disassembledEA {
	mem := dOpcode(handler, address)
	displacement := int(uint8(mem))
	reg := "%s%d."
	if (mem & 0x8000) == 0x8000 {
		reg = fmt.Sprintf(reg, "a", (mem>>12)&7)
	} else {
		reg = fmt.Sprintf(reg, "d", (mem>>12)&7)
	}
	if (mem & 0x0800) == 0x0800 {
		reg += "l"
	} else {
		reg += "w"
	}
	return disassembledEA{fmt.Sprintf("%d(a%d,%s)", displacement, ea.register, reg), ea.o, mem}
}

func (ea *eaPCWithIndex) disassemble(handler AddressHandler, address uint32) disassembledEA {
	mem := dOpcode(handler, address)
	displacement := int(uint8(mem))
	reg := "%s%d."
	if (mem & 0x8000) == 0x8000 {
		reg = fmt.Sprintf(reg, "a", (mem>>12)&7)
	} else {
		reg = fmt.Sprintf(reg, "d", (mem>>12)&7)
	}
	if (mem & 0x0800) == 0x0800 {
		reg += "l"
	} else {
		reg += "w"
	}
	return disassembledEA{fmt.Sprintf("%d(pc,%s)", displacement, reg), ea.o, mem}
}

func (ea *eaAbsoluteWord) disassemble(handler AddressHandler, address uint32) disassembledEA {
	mem := dOpcode(handler, address)
	return disassembledEA{fmt.Sprintf("$"+Word.formatter, mem), ea.o, mem}
}

func (ea *eaAbsoluteLong) disassemble(handler AddressHandler, address uint32) disassembledEA {
	mem := dOpcode(handler, address)
	return disassembledEA{fmt.Sprintf("$"+Long.formatter, mem), ea.o, mem}
}

func (ea *eaImmediate) disassemble(handler AddressHandler, address uint32) disassembledEA {
	mem := dOpcode(handler, address)
	return disassembledEA{fmt.Sprintf("#$"+ea.o.formatter, mem), ea.o, mem}
}

func dOpcode(handler AddressHandler, address uint32) uint32 {
	opcode, err := handler.Read(Word, address)
	if err != nil {
		panic(err)
	}
	return opcode
}
