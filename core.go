package cpu

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type (
	// Error type for CPU Errors
	Error int32
	// Type of CPU
	Type int32
	// Signal external CPU events
	Signal int32

	// floatX80 not supported yet
	floatX80 float64

	// Reader accessor for read accesses
	Reader func(int32, *Size) int32
	// Writer accessor for write accesses
	Writer func(int32, *Size, int32)
	// Reset prototype
	Reset func()

	// AddressArea container for address space area
	AddressArea struct {
		name  string
		read  Reader
		write Writer
		reset Reset
	}

	// AddressBus for accessing address areas
	AddressBus interface {
		read(address int32, s *Size) int32
		write(address int32, s *Size, value int32)
		reset()
	}

	opcode struct {
		instruction         instruction
		match, mask, eaMask uint16
		cycles              map[rune]*int
	}
)

// Exceptions handled by emulation
const (
	MMUAtcEntries = 22  // 68851 has 64, 030 has 22
	M68KICSize    = 128 // instruction cache size

	BusError                Error = 2
	AdressError             Error = 3
	IllegalInstruction      Error = 4
	ZeroDivideError         Error = 5
	PrivilegeViolationError Error = 8
	UnintializedInterrupt   Error = 15

	HaltSignal Signal = iota
	ResetSignal
	Int1Signal
	Int2Signal
	Int3Signal
	Int4Signal
	Int5Signal
	Int6Singal
	Int7Signal
)

var opcodeTable = []*opcode{}

func (e Error) Error() string {
	return fmt.Sprintf("CPU error %v", int32(e))
}

func validEA(opcode, mask uint16) bool {
	if mask == 0 {
		return true
	}
	switch (opcode & 0x3f) >> 3 {
	case 0x00:
		return (mask & 0x800) != 0
	case 0x01:
		return (mask & 0x400) != 0
	case 0x02:
		return (mask & 0x200) != 0
	case 0x03:
		return (mask & 0x100) != 0
	case 0x04:
		return (mask & 0x080) != 0
	case 0x05:
		return (mask & 0x040) != 0
	case 0x06:
		switch opcode & 0x3f {
		case 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37:
			return (mask & 0x020) != 0
		case 0x38:
			return (mask & 0x010) != 0
		case 0x39:
			return (mask & 0x008) != 0
		case 0x3a:
			return (mask & 0x002) != 0
		case 0x3b:
			return (mask & 0x001) != 0
		case 0x3c:
			return (mask & 0x004) != 0
		}
	}
	return false
}

func illegal(c *M68K) {
	panic(IllegalInstruction)
}

func buildInstructionTable(c *M68K, r rune) {
	var counter int
	for _, opcode := range opcodeTable {
		match := opcode.mask
		if opcode.cycles[r] != nil {
			mask := opcode.mask
			for value := uint16(0); ; {
				index := match | value

				if validEA(index, opcode.eaMask) {
					// if c.instructions[index] != nil {
					// 	fmt.Printf("instruction 0x%04x already set\n", index)
					// }
					counter++
					c.instructions[index] = opcode.instruction
					c.cycles[index] = *opcode.cycles[r]
				}

				value = ((value | mask) + 1) & ^mask
				if value == 0 {
					break
				}
			}
		}
	}
	log.Printf("added %d cpu instructions", counter)

	for i := range c.instructions {
		if c.instructions[i] == nil {
			c.instructions[i] = illegal
			c.cycles[i] = 4
		}
	}

}

func addOpcode(ins instruction, match, mask uint16, eaMask uint16, cycles ...string) {
	cycleMap := map[rune]*int{}
	for _, entry := range cycles {
		parts := strings.Split(entry, ":")
		c := parts[1]
		cnt := toInt(c, 10)
		for _, r := range parts[0] {
			cycleMap[r] = &cnt
		}
	}
	opcodeTable = append(opcodeTable, &opcode{ins, match, mask, eaMask, cycleMap})
}

func toInt(s string, base int) int {
	if v, x := strconv.ParseInt(s, base, 64); x == nil {
		return int(v)
	}
	panic(x)
}
