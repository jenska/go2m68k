package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/jenska/atari2go/m68k"
	"github.com/jenska/atari2go/mem"
)

func main() {
	flag.Parse()
	glog.Info("Starting atari2go...")
	m := mem.NewMemoryHandler(1024 * 1024)
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

	cycles := 0

	start := time.Now()
	for i := 1; i < 10000; i++ {
		cpu.PC = 0x1000
		for reg := 0; reg < 8; reg++ {
			for imm := 0; imm < 256; imm++ {
				//	cycles += cpu.Execute() // moveq #imm, reg
			}
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("took %s to perform %d cycles\n", elapsed, cycles)

	atariKiloCycleSpeed := float64(1.0) / float64(8000.0)
	atari2goKiloCycleSpeed := elapsed.Seconds() / float64(cycles/1000)

	fmt.Printf("ST %fmsec EMU %fmsec => we are %4.2f times faster",
		atariKiloCycleSpeed, atari2goKiloCycleSpeed, atariKiloCycleSpeed/atari2goKiloCycleSpeed)
}
