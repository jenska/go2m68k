package mem

import (
	"encoding/binary"
	"testing"

	"github.com/jenska/go2m68k/pkg/cpu"
	"github.com/stretchr/testify/assert"
)

func TestROMRead(t *testing.T) {
	var b [1024]byte
	for i := 0; i < 256; i++ {
		binary.BigEndian.PutUint32(b[i*4:], uint32(i))
	}

	a := NewROM(1024, b[:])
	read := a.read
	for _, o := range []*cpu.Operand{cpu.Byte, cpu.Word, cpu.Long} {
		v, err := read(1024, o)
		assert.Equal(t, 0, v)
		assert.Equal(t, nil, err)
		v, err = read(1023, o)
		assert.Equal(t, 0, v)
		assert.NotEqual(t, nil, err)
	}

	v, _ := read(1024, cpu.Long)
	assert.Equal(t, 0, v)
	v, _ = read(1028, cpu.Long)
	assert.Equal(t, 1, v)
	v, _ = read(1032, cpu.Long)
	assert.Equal(t, 2, v)

	v, _ = read(1030, cpu.Word)
	assert.Equal(t, 1, v)
	v, _ = read(1034, cpu.Word)
	assert.Equal(t, 2, v)
}
