package main

import (
	"os"
	"time"

	"github.com/jenska/m68k"
)

func main() {
	if rom, err := os.ReadFile("./emutos-192k-1.2.1/etos192de.img"); err == nil {
		busController := m68k.NewBusController(
			m68k.BaseRAM(0x1000, 0xfc0000, 1024*1024),
			m68k.ROM(0xfc0000, rom),
		)
		cpu := m68k.New(m68k.M68000, busController)

		signals := make(chan uint16)
		go cpu.Execute(signals)
		time.Sleep(1 * time.Second)
		signals <- m68k.ResetSignal
	} else {
		panic(err)
	}
}
