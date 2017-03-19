package cpu

type Operand struct {
	Size        int
	AlignedSize int
	Msb         int
	Mask        int
	Ext         string
	formatter   string
}

var Byte = &Operand{1, 2, 0x80, 0xff, ".b", "%02x"}
var Word = &Operand{2, 2, 0x8000, 0xffff, ".w", "%04x"}
var Long = &Operand{4, 4, 0x80000000, 0xffffffff, ".l", "%08x"}
