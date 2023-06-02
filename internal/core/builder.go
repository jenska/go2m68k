package core

import (
	"fmt"
	"log"
)

type (
	opcode struct {
		name                string
		instruction         Instruction
		match, mask, eaMask uint16
	}

	buildInstructionTable func(c *Core)
)

var opcodeTable = []*opcode{}

var M68000InstructionSet buildInstructionTable = func(c *Core) {
	var counter int
	for _, opcode := range opcodeTable {
		match := opcode.match
		mask := opcode.mask
		for value := uint16(0); ; {
			index := match | value
			if validEA(index, opcode.eaMask) {
				if c.instructions[index] != nil {
					panic(fmt.Errorf("instruction 0x%04x (%s) already set", index, opcode.name))
				} else {
					counter++
				}
				c.instructions[index] = opcode.instruction
			}

			value = ((value | mask) + 1) & ^mask
			if value == 0 {
				break
			}
		}
	}

	log.Printf("%d cpu instructions available", counter)
}

func Build(r Reader, w Writer, reset func()) func(buildInstructionTable) *Core {
	return func(it buildInstructionTable) *Core {
		c := &Core{}
		c.SR = NewStatusRegister(0x2700)
		c.Regs = make([]uint32, 16)
		c.D = c.Regs[0:8]
		c.A = c.Regs[8:16]
		c.SP = &c.A[7]

		c.readRaw = r
		c.writeRaw = w
		c.resetHandler = reset
		it(c)
		c.Reset()
		return c
	}
}

func Register(name string, ins Instruction, match, mask uint16, eaMask uint16) {
	log.Printf("registering instruction '%s (base $%04x. mask $%04x, ea $%04x)'\n", name, match, mask, eaMask)
	opcodeTable = append(opcodeTable, &opcode{name, ins, match, mask, eaMask})
}

// EA Masks
const (
	MaskDataRegister    = 0x0800
	MaskAddressRegister = 0x0400
	MaskIndirect        = 0x0200
	MaskPostIncrement   = 0x0100
	MaskPreDecrement    = 0x0080
	MaskDisplacement    = 0x0040
	MaskIndex           = 0x0020
	MaskAbsoluteShort   = 0x0010
	MaskAbsoluteLong    = 0x0008
	MaskImmediate       = 0x0004
	MaskPCDisplacement  = 0x0002
	MaskPCIndex         = 0x0001
)

func validEA(opcode, mask uint16) bool {
	if mask == 0 {
		return true
	}

	switch opcode & 0x3f {
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07:
		return (mask & MaskDataRegister) != 0
	case 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f:
		return (mask & MaskAddressRegister) != 0
	case 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17:
		return (mask & MaskIndirect) != 0
	case 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f:
		return (mask & MaskPostIncrement) != 0
	case 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27:
		return (mask & MaskPreDecrement) != 0
	case 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f:
		return (mask & MaskDisplacement) != 0
	case 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37:
		return (mask & MaskIndex) != 0
	case 0x38:
		return (mask & MaskAbsoluteShort) != 0
	case 0x39:
		return (mask & MaskAbsoluteLong) != 0
	case 0x3a:
		return (mask & MaskPCDisplacement) != 0
	case 0x3b:
		return (mask & MaskPCIndex) != 0
	case 0x3c:
		return (mask & MaskImmediate) != 0
	}
	return false
}
