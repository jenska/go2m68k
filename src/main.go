package main

import (
	"fmt"
	"cpu"
)

func main() {
	fmt.Println("Hello World!")
	var c cpu.M68k

	c.SR.Get()
}
