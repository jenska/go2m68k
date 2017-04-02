package cpu

//import glog "github.com/golang/glog"

type Instruction func() int
type InstructionFactory interface {
	Instruction(cpu *M68k, opcode uint16) Instruction
}

type M68k struct {
	A        [8]uint32
	D        [8]uint32
	SR       StatusRegister
	SSP, USP uint32
	PC       uint32

	memory       AddressHandler
	instructions []Instruction
	eaVec        []EA
}

func NewM68k(memory AddressHandler) *M68k {
	cpu := &M68k{}
	cpu.memory = memory
	cpu.SR = NewStatusRegister(cpu)
	cpu.eaVec = NewEAVectors(cpu)
	cpu.instructions = []Instruction{move}
	cpu.Reset()
	return cpu
}

func (cpu *M68k) Reset() {
	// TODO reset
}

func (cpu *M68k) Execute() {
	/*	defer func() {
		//r := recover()

	}n*/

}

func (cpu *M68k) String() string {
	// TODO
	return "M68000"
}

func (cpu *M68k) read(o *Operand, address uint32) uint32 {
	address &= 0x00ffffff
	if v, ok := cpu.memory.Mem(o, address); ok {
		return v
	}
	// TODO raise exception
	return 0
}

func (cpu *M68k) write(o *Operand, address uint32, value uint32) {
	address &= 0x00ffffff
	if cpu.memory.setMem(o, address, value) {
		return
	}
	// TODO raise exception
}

func (cpu *M68k) popPC(o *Operand) uint32 {
	result := cpu.read(o, cpu.PC)
	cpu.PC += o.Size
	return result
}

func (cpu *M68k) pushPC(o *Operand, v uint32) {
	cpu.PC -= o.Size
	cpu.write(o, cpu.PC, v)
}
