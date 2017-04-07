package main

import (
	"flag"
	"fmt"

	glog "github.com/golang/glog"
	m68k "github.com/jenska/atari2go/m68k"
)

func main() {
	flag.Parse()
	glog.Info("Starting atari2go...")
	m := m68k.NewMemoryHandler(1024 * 1024)
	cpu := m68k.NewM68k(m)
	cpu.SR.SetS(true)
	fmt.Println(cpu)

	var address uint32 = 0x1000

	for reg := 0; reg < 8; reg++ {
		for imm := 0; imm < 256; imm++ {
			opcode := 0x7000 + (reg << 9) + imm
			//fmt.Printf("set opcode $%04x imm=%02x \n", opcode, imm)
			cpu.Write(m68k.Word, address, uint32(opcode))
			address += m68k.Word.Size
		}
	}
	cpu.PC = 0x1000
	for reg := 0; reg < 8; reg++ {
		cpu.D[reg] = 1
		for imm := 0; imm < 256; imm++ {
			cpu.Execute() // moveq #0, D0
		}
	}
}
