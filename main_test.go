package m68k

import (
	"os"
	"sync"
	"testing"
)

type (
	DummyReader int
	DummyWriter int
)

func (DummyReader) Read8(address uint32) uint8 {
	return 1
}

func (DummyReader) Read16(address uint32) uint16 {
	return 2
}

func (DummyReader) Read32(address uint32) uint32 {
	return 3
}

func (DummyWriter) Write8(address uint32, value uint8) {
}

func (DummyWriter) Write16(address uint32, value uint16) {
}

func (DummyWriter) Write32(address uint32, value uint32) {
}

func TestBootEnvironment(t *testing.T) {
	var wg sync.WaitGroup
	var dummyReader DummyReader
	var dummyWriter DummyWriter
	if rom, err := os.ReadFile("./emutos-192k-1.2.1/etos192de.img"); err == nil {
		busController := NewBusController(
			BaseRAM(0x1000, 0xfc0000, 1024*1024),
			ROM(0xfc0000, rom),
			ChipArea(0xff8000, 4096, dummyReader, dummyWriter, nil),
		)
		cpu := New(M68000, busController)

		signals := make(chan uint16)
		wg.Add(1)
		go func() {
			defer wg.Done()
			cpu.Execute(signals)
		}()
		wg.Wait()
	} else {
		panic(err)
	}
}
