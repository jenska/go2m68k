package main

import (
	"cpu"
	"flag"
	"fmt"
	glog "github.com/golang/glog"
)

func main() {
	flag.Parse()
	glog.Info("Starting atari2go...")
	m := cpu.NewMemoryHandler(1024 * 1024)
	c := cpu.NewM68k(m)
	c.SR.SetS(true)
	fmt.Println(c)
}
