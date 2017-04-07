package m68k

var shortExecutionTime = [][]int{{4, 4, 8, 8, 8, 12, 14, 12, 16},
	{4, 4, 8, 8, 8, 12, 14, 12, 16}, {8, 8, 12, 12, 12, 16, 18, 16, 20},
	{8, 8, 12, 12, 12, 16, 18, 16, 20}, {10, 10, 14, 14, 14, 18, 20, 18, 22},
	{12, 12, 16, 16, 16, 20, 22, 20, 24}, {14, 14, 18, 18, 18, 22, 24, 22, 26},
	{12, 12, 16, 16, 16, 20, 22, 20, 24}, {16, 16, 20, 20, 20, 24, 26, 24, 28},
	{12, 12, 16, 16, 16, 20, 22, 20, 24}, {14, 14, 18, 18, 18, 22, 24, 22, 26},
	{8, 8, 12, 12, 12, 16, 18, 16, 20}}

var longExecutionTime = [][]int{{4, 4, 12, 12, 12, 16, 18, 16, 20},
	{4, 4, 12, 12, 12, 16, 18, 16, 20}, {12, 12, 20, 20, 20, 24, 26, 24, 28},
	{12, 12, 20, 20, 20, 24, 26, 24, 28}, {14, 14, 22, 22, 22, 26, 28, 26, 30},
	{16, 16, 24, 24, 24, 28, 30, 28, 32}, {18, 18, 26, 26, 26, 30, 32, 30, 34},
	{16, 16, 24, 24, 24, 28, 30, 28, 32}, {20, 20, 28, 28, 28, 32, 34, 32, 36},
	{16, 16, 24, 24, 24, 28, 30, 28, 32}, {18, 18, 26, 26, 26, 30, 32, 30, 34},
	{12, 12, 20, 20, 20, 24, 26, 24, 28}}

var m2RTiming = []int{0, 0, 12, 12, 0, 16, 18, 16, 20, 16, 18}
var r2MTiming = []int{0, 0, 8, 0, 8, 12, 14, 12, 16}

func move(cpu *M68k) int {
	return 0
}

/*


	opcode := cpu.ir
	reg := (opcode >> 9) & 0x7
	data := int(int8(opcode & 0xff))
*/
type moveq struct {
	value  uint32
	target *uint32
}

func (m *moveq) Execute(cpu *M68k) int {
	*m.target = m.value
	cpu.SR.N, cpu.SR.Z = m.value < 0, m.value == 0
	cpu.SR.C, cpu.SR.V = false, false
	return 4
}
