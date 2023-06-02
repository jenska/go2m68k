package main

import (
	"time"

	"github.com/jenska/m68k"
)

func main() {
	busController := m68k.NewBusController(m68k.BaseRAM(0x1000, 0xfc0000, 1024*1024))
	cpu := m68k.New(m68k.M68000, busController)

	signals := make(chan uint16)
	go cpu.Execute(signals)
	time.Sleep(1 * time.Second)
	signals <- m68k.ResetSignal
}
