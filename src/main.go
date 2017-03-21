package main

import (
	"cpu"
	"fmt"
	"mem"
)

var Memory = mem.NewPhysicalAddressSpace(1024*1024)
var CPU = cpu.NewM68k(Memory)

func main() {
	fmt.Println("Hello World!")
	var c cpu.M68k

	c.SR.Get()
}
