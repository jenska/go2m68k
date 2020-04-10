package cpu

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRAM(t *testing.T) {
	page1 := NewRAMArea("TestRAM1", 1024*1024)
	page2 := NewRAMArea("TestRAM2", 1024*1024)
	assert.Equal(t, "TestRAM1", page1.Name)
	assert.NotNil(t, page1.Read)
	assert.NotNil(t, page1.Write)
	assert.NotNil(t, page1.Reset)
	assert.NotNil(t, page1.Raw)

	io := NewIOManager(1024*1024, page1).AddArea(3*1024*1024, 1024*1024, page2)
	io.Write(0, Long, 0)
	assert.Equal(t, 0, io.Read(0, Long))

	io.Write(0, Long, 1)
	assert.Equal(t, 1, io.Read(0, Long))

	io.Write(4, Long, 1)
	assert.Equal(t, 1, io.Read(4, Long))

	io.Write(0, Long, 300)
	assert.Equal(t, 300, io.Read(0, Long))

	io.Write(4, Long, 400)
	assert.Equal(t, 400, io.Read(4, Long))

	assert.Panics(t, func() {
		io.Write(1024*1024, Byte, 2)
	})
	assert.Panics(t, func() {
		io.Write(1024*1024-1, Word, 2)
	})
	assert.Panics(t, func() {
		io.Write(1024*1024-3, Long, 2)
	})

	io.Reset()
	assert.Equal(t, 0, io.Read(0, Long))
}

func TestROM(t *testing.T) {
	rom := make([]byte, 1024)
	for i := range rom {
		rom[i] = 1
	}
	page1 := NewROMArea("TestROM", rom)

	io := NewIOManager(1024, page1)
	assert.Panics(t, func() {
		io.Write(0, Long, 0)
	})
	assert.Equal(t, 0x01010101, io.Read(0, Long))

	assert.Panics(t, func() {
		io.Write(0, Long, 1)
	})
	assert.Equal(t, 0x01010101, io.Read(0, Long))

	assert.Panics(t, func() {
		io.Write(4, Long, 1)
	})
	assert.Equal(t, 0x01010101, io.Read(4, Long))

	assert.Panics(t, func() {
		io.Write(0, Long, 300)
	})
	assert.Equal(t, 0x01010101, io.Read(0, Long))

	assert.Panics(t, func() {
		io.Write(4, Long, 400)
	})
	assert.Equal(t, 0x01010101, io.Read(4, Long))
}

func TestCache(t *testing.T) {
	page1 := NewRAMArea("TestRAM1", 1024)
	page2 := NewRAMArea("TestRAM2", 1024)
	page3 := NewRAMArea("TestRAM3", 1024)
	page4 := NewRAMArea("TestRAM4", 1024)

	io := NewIOManager(1024, page1).AddArea(2048, 1024, page2).AddArea(4096, 1024, page3).AddArea(16000, 1024, page4)
	for i := 1; i < 555; i++ {
		if i%2 == 0 {
			io.Write(4096, Byte, 0)
		}
		if i%3 == 0 {
			io.Write(2048, Word, 0)
		}
		if i%4 == 0 {
			io.Write(0, Long, 0)
		}
		if i%5 == 0 {
			io.Write(16000, Long, 0)
		}
	}
	sort.Sort(io.areas)
	//	for _, h := range *io.areas {
	//		fmt.Println(h.accessCnt)
	//	}
	assert.True(t, sort.IsSorted(io.areas))
}
