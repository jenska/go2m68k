package core

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignedExtend(t *testing.T) {
	o := []Operand{Byte, Word, Long}

	assert.Equal(t, int32(0), o[0].SignedExtend(0))
	assert.Equal(t, int32(0), o[1].SignedExtend(0))
	assert.Equal(t, int32(0), o[2].SignedExtend(0))

	var b uint32 = 0xff
	var w uint32 = 0xff
	var l uint32 = 0xff

	assert.Equal(t, int32(-1), o[0].SignedExtend(b))
	assert.Equal(t, int32(255), o[1].SignedExtend(w))
	assert.Equal(t, int32(255), o[2].SignedExtend(l))

	w = 0xffff
	l = 0xffff
	assert.Equal(t, int32(-1), o[0].SignedExtend(b))
	assert.Equal(t, int32(-1), o[1].SignedExtend(w))
	assert.Equal(t, int32(0xffff), o[2].SignedExtend(l))

	l = 0xffffffff
	assert.Equal(t, int32(-1), o[0].SignedExtend(b))
	assert.Equal(t, int32(-1), o[1].SignedExtend(w))
	assert.Equal(t, int32(-1), o[2].SignedExtend(l))
}

func TestUnsignedExtend(t *testing.T) {
	o := []Operand{Byte, Word, Long}

	assert.Equal(t, uint32(0), o[0].UnsignedExtend(0))
	assert.Equal(t, uint32(0), o[1].UnsignedExtend(0))
	assert.Equal(t, uint32(0), o[2].UnsignedExtend(0))

	assert.Equal(t, uint32(255), o[0].UnsignedExtend(0xff))
	assert.Equal(t, uint32(255), o[1].UnsignedExtend(0xff))
	assert.Equal(t, uint32(255), o[2].UnsignedExtend(0xff))

	assert.Equal(t, uint32(0xff), o[0].UnsignedExtend(0xffff))
	assert.Equal(t, uint32(0xffff), o[1].UnsignedExtend(0xffff))
	assert.Equal(t, uint32(0xffff), o[2].UnsignedExtend(0xffff))

	assert.Equal(t, uint32(0xff), o[0].UnsignedExtend(0xffffffff))
	assert.Equal(t, uint32(0xffff), o[1].UnsignedExtend(0xffffffff))
	assert.Equal(t, uint32(0xffffffff), o[2].UnsignedExtend(0xffffffff))
}

func TestMSB(t *testing.T) {
	assert.True(t, !Byte.MSB(0))
	assert.True(t, !Word.MSB(0))
	assert.True(t, !Long.MSB(0))

	var b uint32 = 0xff
	var w uint32 = 0xff
	var l uint32 = 0xff

	assert.True(t, Byte.MSB(b))
	assert.True(t, !Word.MSB(w))
	assert.True(t, !Long.MSB(l))

	w = 0xffff
	l = 0xffff
	assert.True(t, Byte.MSB(b))
	assert.True(t, Word.MSB(w))
	assert.True(t, !Long.MSB(l))

	l = 0xffffffff
	assert.True(t, Byte.MSB(b))
	assert.True(t, Word.MSB(w))
	assert.True(t, Long.MSB(l))
}

func TestReadWrite(t *testing.T) {
	b := make([]byte, 8)
	d := make([]byte, 8)

	binary.BigEndian.PutUint32(b, uint32(0x12345678))

	assert.Equal(t, uint32(0x12), Byte.Read(b[0:]))
	assert.Equal(t, uint32(0x34), Byte.Read(b[1:]))
	assert.Equal(t, uint32(0x56), Byte.Read(b[2:]))
	assert.Equal(t, uint32(0x78), Byte.Read(b[3:]))

	Byte.Write(0x12, d[0:])
	Byte.Write(0x34, d[1:])
	Byte.Write(0x56, d[2:])
	Byte.Write(0x78, d[3:])
	assert.Equal(t, uint32(0x12), Byte.Read(d[0:]))
	assert.Equal(t, uint32(0x34), Byte.Read(d[1:]))
	assert.Equal(t, uint32(0x56), Byte.Read(d[2:]))
	assert.Equal(t, uint32(0x78), Byte.Read(d[3:]))

	assert.Equal(t, uint32(0x1234), Word.Read(b[0:]))
	assert.Equal(t, uint32(0x3456), Word.Read(b[1:]))
	assert.Equal(t, uint32(0x5678), Word.Read(b[2:]))
	assert.Equal(t, uint32(0x7800), Word.Read(b[3:]))

	Word.Write(0x1234, d[0:])
	assert.Equal(t, uint32(0x1234), Word.Read(d[0:]))
	Word.Write(0x3456, d[1:])
	assert.Equal(t, uint32(0x3456), Word.Read(d[1:]))
	Word.Write(0x5678, d[2:])
	assert.Equal(t, uint32(0x5678), Word.Read(d[2:]))
	Word.Write(0x7800, d[3:])
	assert.Equal(t, uint32(0x7800), Word.Read(d[3:]))

	assert.Equal(t, uint32(0x12345678), Long.Read(b[0:]))
	assert.Equal(t, uint32(0x34567800), Long.Read(b[1:]))
	assert.Equal(t, uint32(0x56780000), Long.Read(b[2:]))
	assert.Equal(t, uint32(0x78000000), Long.Read(b[3:]))

	Long.Write(0x1234, d[0:])
	assert.Equal(t, uint32(0x1234), Long.Read(d[0:]))
	Long.Write(0x3456, d[1:])
	assert.Equal(t, uint32(0x3456), Long.Read(d[1:]))
	Long.Write(0x5678, d[2:])
	assert.Equal(t, uint32(0x5678), Long.Read(d[2:]))
	Long.Write(0x7800, d[3:])
	assert.Equal(t, uint32(0x7800), Long.Read(d[3:]))
}

func TestWriteToLong(t *testing.T) {
	var a uint32 = 0xffffffff
	Byte.WriteToLong(0x12, &a)
	assert.Equal(t, uint32(0xffffff12), a)
	Word.WriteToLong(0x1234, &a)
	assert.Equal(t, uint32(0xffff1234), a)
	Long.WriteToLong(0x12345678, &a)
	assert.Equal(t, uint32(0x12345678), a)
}

func TestSize(t *testing.T) {
	assert.EqualValues(t, ByteSize, Byte.Size())
	assert.EqualValues(t, WordSize, Word.Size())
	assert.EqualValues(t, LongSize, Long.Size())
}
