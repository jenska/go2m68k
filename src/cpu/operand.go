package cpu

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

func (o *Operand) isNegative(value uint32) bool {
	return (o.Msb & value) != 0
}

func (o *Operand) set(target *uint32, value uint32) {
	*target = (*target & ^o.Mask) | (value & o.Mask)
}

func (o *Operand) getSigned(value uint32) int32 {
	v := uint32(value)
	if o.isNegative(v) {
		return int32(v | ^o.Mask)
	}
	return int32(v & o.Mask)
}

func (o *Operand) get(value uint32) uint32 {
	return value & o.Mask
}