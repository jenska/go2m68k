package main

import (
	"cpu"
	"fmt"
)

func main() {
	fmt.Println("Hello World!")
	var c cpu.M68k
	c.SR.Get()
}
