package main

import (
	"github.com/jenska/atari2go/mem"
)

func main() {
	mem.NewAddressBus(mem.NewRAM(0, 1024*1024))
}
