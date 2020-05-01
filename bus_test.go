package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRAM(t *testing.T) {
	page1 := NewRAMArea(1024 * 1024)
	page2 := NewRAMArea(1024 * 1024)
	assert.NotNil(t, page1.read)
	assert.NotNil(t, page1.write)
	assert.NotNil(t, page1.reset)

	io := NewAddressBusBuilder().AddArea(0, 1024*1024, page1).AddArea(3*1024*1024, 1024*1024, page2).Build()

	io.write(0, Long, 0)
	assert.Equal(t, int32(0), io.read(0, Long))

	io.write(0, Long, 1)
	assert.Equal(t, int32(1), io.read(0, Long))

	io.write(4, Long, 1)
	assert.Equal(t, int32(1), io.read(4, Long))

	io.write(0, Long, 300)
	assert.Equal(t, int32(300), io.read(0, Long))

	io.write(4, Long, 400)
	assert.Equal(t, int32(400), io.read(4, Long))

	assert.Panics(t, func() {
		io.write(1024*1024, Byte, 2)
	})
	assert.Panics(t, func() {
		io.write(1024*1024-1, Word, 2)
	})
	assert.Panics(t, func() {
		io.write(1024*1024-3, Long, 2)
	})

	assert.Panics(t, func() {
		io.read(1024*1024, Byte)
	})

	assert.Panics(t, func() {
		io.read(1024*1024-1, Word)
	})

	assert.Panics(t, func() {
		io.read(1024*1024-3, Long)
	})

	assert.Panics(t, func() {
		io.read(1024*1024+1024, Byte)
	})
	assert.Panics(t, func() {
		io.read(1024*1024+1024, Word)
	})
	assert.Panics(t, func() {
		io.read(1024*1024+1024, Long)
	})
	assert.Panics(t, func() {
		io.write(1024*1024+10124, Byte, 2)
	})
	assert.Panics(t, func() {
		io.write(1024*1024+1024, Word, 2)
	})
	assert.Panics(t, func() {
		io.write(1024*1024+1024, Long, 2)
	})

	io.reset()
	assert.Equal(t, int32(0), io.read(0, Long))
}

func TestROM(t *testing.T) {
	rom := make([]byte, 1024)
	for i := range rom {
		rom[i] = 1
	}
	page1 := NewROMArea(rom)
	io := NewAddressBusBuilder().AddArea(0, 1024, page1).Build()
	assert.Panics(t, func() {
		io.write(0, Long, 0)
	})
	assert.Equal(t, int32(0x01010101), io.read(0, Long))

	assert.Panics(t, func() {
		io.write(0, Long, 1)
	})
	assert.Equal(t, int32(0x01010101), io.read(0, Long))

	assert.Panics(t, func() {
		io.write(4, Long, 1)
	})
	assert.Equal(t, int32(0x01010101), io.read(4, Long))

	assert.Panics(t, func() {
		io.write(0, Long, 300)
	})
	assert.Equal(t, int32(0x01010101), io.read(0, Long))

	assert.Panics(t, func() {
		io.write(4, Long, 400)
	})
	assert.Equal(t, int32(0x01010101), io.read(4, Long))
}

func TestCache(t *testing.T) {
	page1 := NewRAMArea(1024)
	page2 := NewRAMArea(1024)
	page3 := NewRAMArea(1024)
	page4 := NewRAMArea(1024)

	io := NewAddressBusBuilder().AddArea(0, 1024, page1).AddArea(2048, 1024, page2).AddArea(4096, 1024, page3).AddArea(16000, 1024, page4).Build()
	for i := 1; i < 555; i++ {
		if i%2 == 0 {
			io.write(4096, Byte, 0)
		}
		if i%3 == 0 {
			io.write(2048, Word, 0)
		}
		if i%4 == 0 {
			io.write(0, Long, 0)
		}
		if i%5 == 0 {
			io.write(16000, Long, 0)
		}
	}
	io.reset()
}
