package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jenska/atari2go/cpu"
)

func TestRAMRead(t *testing.T) {
	a := NewRAM(1024, 1024)
	read := a.read
	for _, o := range []*cpu.Operand{cpu.Byte, cpu.Word, cpu.Long} {
		v, err := read(1024, o)
		assert.Equal(t, 0, v)
		assert.Equal(t, nil, err)
		v, err = read(1023, o)
		assert.Equal(t, 0, v)
		assert.NotEqual(t, nil, err)
	}
}

func TestRAMWrite(t *testing.T) {
	a := NewRAM(1024, 1024)
	read := a.read
	write := a.write

	write(1024, cpu.Long, 1)
	v, _ := read(1024, cpu.Long)
	assert.Equal(t, 1, v)
}

func BenchmarkRAMLongWrite(b *testing.B) {
	a := NewRAM(1024, 1024)
	read := a.read
	write := a.write

	for i := 0; i < 256; i++ {
		offset := 1024 + cpu.Address(i<<2)
		write(offset, cpu.Long, i)
		read(offset, cpu.Long)
	}
}
