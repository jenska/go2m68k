package mem

import (
	"reflect"
	"testing"

	"github.com/jenska/atari2go/cpu"
	"github.com/stretchr/testify/assert"
)

func TestNewRAM(t *testing.T) {
	a := NewRAM(0, 1024)
	assert.Equal(t, a.start, cpu.Address(0))
	assert.Equal(t, a.end, cpu.Address(1024))
	assert.NotNil(t, a.read)
	assert.NotNil(t, a.write)
}

func TestNewROM(t *testing.T) {
	a := NewROM(0, make([]byte, 1024))
	assert.Equal(t, a.start, cpu.Address(0))
	assert.Equal(t, a.end, cpu.Address(1024))
	assert.NotNil(t, a.read)
	assert.Nil(t, a.write)
}

func TestNewAddressBus(t *testing.T) {
	ab := NewAddressBus(
		NewRAM(0, 1024),
		NewROM(2048, make([]byte, 1024)),
	)

	v, err := ab.Read(0, cpu.Byte)
	assert.Equal(t, v, 0)
	assert.Equal(t, err, nil)

	v, err = ab.Read(1500, cpu.Byte)
	print(reflect.TypeOf(err))
	assert.IsType(t, err, BusError(0))
	assert.Equal(t, err, cpu.Address(1500))

}
