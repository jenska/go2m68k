package cpu

import (
	"encoding/binary"
	"fmt"
)

type Size struct {
	size  uint32
	align uint32
	mask  uint
	bits  uint
	msb   uint
	fmt   string
	Write func([]byte, int)
	Read  func([]byte) int
}

var (
	Byte = &Size{
		size:  1,
		align: 2,
		mask:  0x000000ff,
		bits:  8,
		msb:   0x00000080,
		fmt:   "$%02x",
		Write: func(sclice []byte, value int) {
			sclice[0] = byte(value)
		},
		Read: func(slice []byte) int {
			return int(slice[0])
		},
	}

	Word = &Size{
		size:  2,
		align: 2,
		mask:  0x0000ffff,
		bits:  16,
		msb:   0x00008000,
		fmt:   "$%04x",
		Write: func(slice []byte, value int) {
			binary.BigEndian.PutUint16(slice, uint16(value))
		},
		Read: func(slice []byte) int {
			return int(binary.BigEndian.Uint16(slice))
		},
	}

	Long = &Size{
		size:  4,
		align: 4,
		mask:  0xffffffff,
		bits:  32,
		msb:   0x80000000,
		fmt:   "$%08x",
		Write: func(sclice []byte, value int) {
			binary.BigEndian.PutUint32(sclice, uint32(value))
		},
		Read: func(slice []byte) int {
			return int(binary.BigEndian.Uint32(slice))
		},
	}
)

func (s *Size) IsNegative(value int) bool {
	return s.msb&uint(value) != 0
}

func (s *Size) HexString(value int) string {
	v := uint(value) & s.mask
	return fmt.Sprintf(s.fmt, v)
}

// SignedHexString returns a signed hex with leading zeroes (-$0001)
func (s *Size) SignedHexString(value int) string {
	if uint(value) == s.msb {
		return s.HexString(value)
	} else if s.IsNegative(value) {
		return "-" + s.HexString(-value)
	}
	return s.HexString(value)
}
