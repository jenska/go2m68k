package m68k

import (
	"sync"
	"testing"
	"time"

	_ "github.com/jenska/m68k/internal/instructions"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	bc := NewBusController(BaseRAM(0x1000, 0xfc0000, 1024*1024), ROM(0xFC0000, make([]byte, 3*64*1024)))
	cpu := New(M68000, bc)
	assert.NotNil(t, cpu)
}

func TestReset(t *testing.T) {
	bc := NewBusController(BaseRAM(0x1000, 0xfc0000, 1024*1024), ROM(0xFC0000, make([]byte, 3*64*1024)))
	cpu := New(M68000, bc)
	assert.NotNil(t, cpu)
	signals := make(chan uint16, 2)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cpu.Execute(signals)
	}()

	select {
	case signals <- ResetSignal:
	case <-time.After(time.Second):
		t.Fatalf("sending reset signal timed out")
	}

	select {
	case signals <- HaltSignal:
	case <-time.After(time.Second):
		t.Fatalf("sending halt signal timed out")
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("CPU did not halt in time")
	}
}

func TestNewBusController(t *testing.T) {
	NewBusController(BaseRAM(0x1000, 0xfc0000, 1024*1024), ROM(0xFC0000, make([]byte, 3*64*1024)))
}
