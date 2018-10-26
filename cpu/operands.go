package cpu

type Operand struct {
	name  string
	short string
	size  int
	mask  uint
	bits  uint
	msb   uint
}

var (
	Byte = &Operand{name: "Byte", short: ".b", size: 1, mask: 0x000000ff, bits: 8, msb: 0x00000080}
	Word = &Operand{name: "Word", short: ".w", size: 2, mask: 0x0000ffff, bits: 16, msb: 0x00008000}
	Long = &Operand{name: "Long", short: ".l", size: 4, mask: 0xffffffff, bits: 32, msb: 0x80000000}
)

func (o *Operand) Write(slice []byte, index uint, value int) {

}

func (o *Operand) Read(slice []byte, index uint) int {
	return 0
}

func (o *Operand) IsNegative(value int) bool {
	return o.msb&uint(value) != 0
}
