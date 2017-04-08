package m68k

import (
	"fmt"
)

type disassembledInstruction struct {
	address uint32
	opcode  uint16
	name    string
	ea      []disassembledEA
}

func (i *disassembledInstruction) size() int {
	size := Word.Size
	for _, a := range i.ea {
		size += a.o.AlignedSize
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
	disassemble(cpu *M68K, address uint32) disassembledEA
}

func (ea *disassembledEA) toHex() string {
	if ea.o != nil {
		return fmt.Sprintf(ea.o.formatter, ea.memory)
	}
	return ""
}

func (ea *eaDataRegister) disassemble(_ *M68K, _ uint32) disassembledEA {
	return disassembledEA{fmt.Sprintf("d%d", ea.register), ea.o, 0}
}

func (ea *eaAddressRegister) disassemble(_ *M68K, _ uint32) disassembledEA {
	return disassembledEA{fmt.Sprintf("a%d", ea.register), ea.o, 0}
}

func (ea *eaAddressRegisterIndirect) disassemble(_ *M68K, _ uint32) disassembledEA {
	return disassembledEA{fmt.Sprintf("(a%d)", ea.register), ea.o, 0}
}

func (ea *eaAddressRegisterPreDec) disassemble(_ *M68K, _ uint32) disassembledEA {
	return disassembledEA{fmt.Sprintf("-(a%d)", ea.register), ea.o, 0}
}

func (ea *eaAddressRegisterPostInc) disassemble(_ *M68K, _ uint32) disassembledEA {
	return disassembledEA{fmt.Sprintf("(a%d)+", ea.register), ea.o, 0}
}

func (ea *eaAddressRegisterWithDisplacement) disassemble(cpu *M68K, address uint32) disassembledEA {
	mem := cpu.Read(Word, address)
	displacement := int(uint16(mem))
	return disassembledEA{fmt.Sprintf("$%04x(a%d)", displacement, ea.register), ea.o, mem}
}

func (ea *eaPCWithDisplacement) disassemble(cpu *M68K, address uint32) disassembledEA {
	mem := cpu.Read(Word, address)
	displacement := int(uint16(mem))
	return disassembledEA{fmt.Sprintf("$%04x(pc)", displacement), ea.o, mem}
}

func (ea *eaAddressRegisterWithIndex) disassemble(cpu *M68K, address uint32) disassembledEA {
	mem := cpu.Read(Word, address)
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

func (ea *eaPCWithIndex) disassemble(cpu *M68K, address uint32) disassembledEA {
	mem := cpu.Read(Word, address)
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

func (ea *eaAbsoluteWord) disassemble(cpu *M68K, address uint32) disassembledEA {
	mem := cpu.Read(Word, address)
	return disassembledEA{fmt.Sprintf("$"+Word.formatter, mem), ea.o, mem}
}

func (ea *eaAbsoluteLong) disassemble(cpu *M68K, address uint32) disassembledEA {
	mem := cpu.Read(Long, address)
	return disassembledEA{fmt.Sprintf("$"+Long.formatter, mem), ea.o, mem}
}

func (ea *eaImmediate) disassemble(cpu *M68K, address uint32) disassembledEA {
	mem := cpu.Read(ea.o, address)
	return disassembledEA{fmt.Sprintf("#$"+ea.o.formatter, mem), ea.o, mem}
}
