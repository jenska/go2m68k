package util

import (
	"testing"

	"github.com/jenska/atari2go/cpu"
	"github.com/jenska/atari2go/mem"
)

var bus = mem.NewAddressBus(
	mem.NewRAM(1014, 1024),
)

func TestDisasmMoveq(t *testing.T) {
	bus.Write(1024, cpu.Word, 0x7000)
	bus.Write(1026, cpu.Word, 0x7001)

	Disassemble(bus, 1024, 4)
}
