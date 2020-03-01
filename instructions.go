package cpu

func (cpu M68000) initISA() {

}

func moveq(cpu M68000) {
	cpu.D[(cpu.IR>>9)&0x7] = int32(cpu.IR & 0xff)
}
