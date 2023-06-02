package core

import "fmt"

const (
	ByteSize = 1
	WordSize = 2
	LongSize = 4
)

type (
	// Represents the operands of M68K CPU
	Operand interface {
		Write(v uint32, dst []byte)
		WriteToLong(v uint32, dst *uint32) // set v into dst
		Read(src []byte) uint32

		UnsignedExtend(v uint32) uint32
		SignedExtend(v uint32) int32
		MSB(v uint32) bool // most significant bit

		Size() uint32
		Align() uint32

		HexString(value uint32) string
	}

	byteOperand struct{}
	wordOperand struct{}
	longOperand struct{}
)

var (
	Byte byteOperand
	Word wordOperand
	Long longOperand
)

func (byteOperand) Write(v uint32, dst []byte) {
	dst[0] = byte(v)
}

func (byteOperand) WriteToLong(v uint32, dest *uint32) {
	*dest = (*dest)&0xffffff00 | (v & 0xff)
}

func (byteOperand) Read(src []byte) uint32 {
	return uint32(src[0])
}

func (byteOperand) UnsignedExtend(v uint32) uint32 {
	return v & 0xff
}

func (byteOperand) SignedExtend(v uint32) int32 {
	return int32(int8(v))
}

func (byteOperand) MSB(v uint32) bool {
	return v&0x80 != 0
}

func (byteOperand) Size() uint32 {
	return ByteSize
}

func (byteOperand) Align() uint32 {
	return WordSize
}

func (byteOperand) HexString(value uint32) string {
	return fmt.Sprintf("%02x", value)
}

func (byteOperand) String() string {
	return "b"
}

func (wordOperand) Write(v uint32, dst []byte) {
	dst[0] = byte(v >> 8)
	dst[1] = byte(v)
}

func (wordOperand) WriteToLong(v uint32, dest *uint32) {
	*dest = (*dest)&0xffff0000 | (v & 0xffff)
}

func (wordOperand) Read(src []byte) uint32 {
	return uint32(src[1]) | uint32(src[0])<<8
}

func (wordOperand) SignedExtend(v uint32) int32 {
	return int32(int16(v))
}

func (wordOperand) UnsignedExtend(v uint32) uint32 {
	return v & 0xffff
}

func (wordOperand) MSB(v uint32) bool {
	return v&0x8000 != 0
}

func (wordOperand) Size() uint32 {
	return WordSize
}

func (wordOperand) Align() uint32 {
	return WordSize
}

func (wordOperand) HexString(value uint32) string {
	return fmt.Sprintf("%04x", value)
}

func (wordOperand) String() string {
	return "w"
}

func (longOperand) Write(v uint32, dst []byte) {
	dst[0] = byte(v >> 24)
	dst[1] = byte(v >> 16)
	dst[2] = byte(v >> 8)
	dst[3] = byte(v)
}

func (longOperand) WriteToLong(v uint32, dest *uint32) {
	*dest = v
}

func (longOperand) Read(src []byte) uint32 {
	return uint32(src[3]) | uint32(src[2])<<8 | uint32(src[1])<<16 | uint32(src[0])<<24
}

func (longOperand) SignedExtend(v uint32) int32 {
	return int32(v)
}

func (longOperand) UnsignedExtend(v uint32) uint32 {
	return v
}

func (longOperand) MSB(v uint32) bool {
	return v&0x80000000 != 0
}

func (longOperand) Size() uint32 {
	return LongSize
}

func (longOperand) Align() uint32 {
	return LongSize
}

func (longOperand) HexString(value uint32) string {
	return fmt.Sprintf("%08x", value)
}

func (longOperand) String() string {
	return "l"
}
