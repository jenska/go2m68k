package cpu

import (
	"encoding/binary"
	"fmt"
)

// Size contains properties of the M68K operad types
type Size struct {
	size  int32
	align int32
	mask  uint32
	bits  uint32
	msb   uint32
	fmt   string
	ext   string
	write func([]byte, int32)
	read  func([]byte) int32
}

var (
	operands = []*Size{Byte, Word, Long}

	// Byte is an M68K operand
	Byte = &Size{
		size:  1,
		align: 2,
		mask:  0x000000ff,
		bits:  8,
		msb:   0x00000080,
		fmt:   "$%02x",
		ext:   ".b",
		write: func(sclice []byte, value int32) {
			sclice[0] = byte(value)
		},
		read: func(slice []byte) int32 {
			return int32(slice[0])
		},
	}

	// Word is an M68K operand type
	Word = &Size{
		size:  2,
		align: 2,
		mask:  0x0000ffff,
		bits:  16,
		msb:   0x00008000,
		fmt:   "$%04x",
		ext:   ".w",
		write: func(slice []byte, value int32) {
			binary.BigEndian.PutUint16(slice, uint16(value))
		},
		read: func(slice []byte) int32 {
			return int32(int16(binary.BigEndian.Uint16(slice)))
		},
	}

	// Long is an M68K operand type
	Long = &Size{
		size:  4,
		align: 4,
		mask:  0xffffffff,
		bits:  32,
		msb:   0x80000000,
		fmt:   "$%08x",
		ext:   ".b",

		write: func(sclice []byte, value int32) {
			binary.BigEndian.PutUint32(sclice, uint32(value))
		},
		read: func(slice []byte) int32 {
			return int32(binary.BigEndian.Uint32(slice))
		},
	}
)

func (s *Size) IsNegative(value int32) bool {
	return s.msb&uint32(value) != 0
}

func (s *Size) HexString(value int32) string {
	v := uint32(value) & s.mask
	return fmt.Sprintf(s.fmt, v)
}

// SignedHexString returns a signed hex with leading zeroes (-$0001)
func (s *Size) SignedHexString(value int32) string {
	if uint32(value) == s.msb {
		return s.HexString(value)
	} else if s.IsNegative(value) {
		return "-" + s.HexString(-value)
	}
	return s.HexString(value)
}

func (s *Size) set(value int32, target *int32) {
	result := (uint32(*target) & ^s.mask) | (uint32(value) & s.mask)
	*target = int32(result)
}
