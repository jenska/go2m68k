package cpu

import (
	"encoding/binary"
)

type Operand struct {
	name  string
	short string
	size  uint
	mask  uint
	bits  uint
	msb   uint
	Write func([]byte, int)
	Read  func([]byte) int
}

var (
	Byte = &Operand{
		name:  "Byte",
		short: ".b",
		size:  1,
		mask:  0x000000ff,
		bits:  8,
		msb:   0x00000080,
		Write: func(sclice []byte, value int) {
			sclice[0] = byte(value)
		},
		Read: func(slice []byte) int {
			return int(slice[0])
		},
	}

	Word = &Operand{
		name:  "Word",
		short: ".w",
		size:  2,
		mask:  0x0000ffff,
		bits:  16,
		msb:   0x00008000,
		Write: func(slice []byte, value int) {
			binary.BigEndian.PutUint16(slice, uint16(value))
		},
		Read: func(slice []byte) int {
			return int(binary.BigEndian.Uint16(slice))
		},
	}

	Long = &Operand{
		name:  "Long",
		short: ".l",
		size:  4,
		mask:  0xffffffff,
		bits:  32,
		msb:   0x80000000,
		Write: func(sclice []byte, value int) {
			binary.BigEndian.PutUint32(sclice, uint32(value))
		},
		Read: func(slice []byte) int {
			return int(binary.BigEndian.Uint32(slice))
		},
	}
)

func (o *Operand) IsNegative(value int) bool {
	return o.msb&uint(value) != 0
}
