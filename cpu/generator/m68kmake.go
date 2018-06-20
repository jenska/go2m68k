package main

import (
	"bytes"
	"flag"
	"log"
	"os"
)

var (
	input  = flag.String("input", "m68kin.go", "input file name")
	output = flag.String("output", "m68kops.go", "output file name")
)

type Generator struct {
	buf bytes.Buffer
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("m68kmake: ")
	flag.Parse()
	if len(*input) == 0 || len(*output == 0) {
		os.Exit(2)
	}

}
