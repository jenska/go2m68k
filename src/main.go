package src

import (
	"fmt"
	"cpu"
	glog "github.com/golang/glog"
)

func main() {
	glog.Info("Starting atari2go...")
	m := cpu.NewMemoryHandler(1024*1024, nil)
	c := cpu.NewM68k(m)
	c.SR.SetS(true)
	fmt.Println(c.SR)
}
