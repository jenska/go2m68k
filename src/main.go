package main

import (
	"cpu"
	"fmt"
)

func main() {
	fmt.Println("Hello World!")
	m := cpu.NewMemoryHandler(1024*1024, nil)
	c := cpu.NewM68k(m)
	c.SR.SetS(true)
	fmt.Println(c.SR)
}
