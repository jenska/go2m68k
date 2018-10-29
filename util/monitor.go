package util

import (
	"fmt"

	"github.com/jenska/atari2go/cpu"
)

func Dump(bus cpu.AddressBus, start cpu.Address, size int) {
	for i, j := 0, 0; i < size; i += 4 {
		b := start + cpu.Address(i)
		a, err := bus.Read(b, cpu.Long)
		if err != nil {
			panic(fmt.Sprintf("invalid read %s", err))
		}
		fmt.Printf("$%08x #%08x\t", b, a)
		j++
		if j == 4 {
			fmt.Println()
			j = 0
		}
	}

}
