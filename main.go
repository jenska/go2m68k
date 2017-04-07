package main

import (
	"flag"
	"fmt"
	"m68k"

	glog "github.com/golang/glog"
)

func main() {
	flag.Parse()
	glog.Info("Starting atari2go...")
	m := m68k.NewMemoryHandler(1024 * 1024)
	c := m68k.NewM68k(m)
	c.SR.SetS(true)
	fmt.Println(c)
}
