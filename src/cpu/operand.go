package cpu

import "fmt"

type Operand struct {
	Size        int
	AlignedSize int
	Msb         uint32
	Mask        uint32
	Ext         string
	formatter   string
}

var Byte = &Operand{1, 2, 0x80, 0xff, ".b", "%02x"}
var Word = &Operand{2, 2, 0x8000, 0xffff, ".w", "%04x"}
var Long = &Operand{4, 4, 0x80000000, 0xffffffff, ".l", "%08x"}

func (o *Operand) toHex(value uint32) string {
	return fmt.Sprintf(o.formatter, value&o.Mask)
}

func (o *Operand) Byte() {

}
