package cpu

import "log"

// Address Bus Mask for 68000 CPU
const busMask = 0x00ffffff

// InitISA68000 Instruction Set Archtiecture for M68000
func (cpu *M68K) SetISA68000() Builder {

	cpu.read = func(a int32, s *Size) int32 {
		if a&1 == 1 && s != Byte {
			panic(AdressError)
		}
		return cpu.bus.read(a&busMask, s)
	}

	cpu.write = func(a int32, s *Size, value int32) {
		if a&1 == 1 && s != Byte {
			panic(AdressError)
		}
		if !cpu.sr.S && a < 0x800 {
			panic(PrivilegeViolationError)
		}
		cpu.bus.write(a&busMask, s, value)
	}

	// // TODO: reusable method for other cpu types
	counter := 0
	for i := range dasmTable {
		opcode := uint16(i)
		dasmTable[i] = dasmIllegal
		for _, info := range opcodeTable {
			if info.cpu == M68000 {
				if (opcode & info.mask) == info.match {
					if info.move && !validEA(eam(uint16(opcode)), 0xbf8) {
						continue
					}
					if validEA(opcode, info.eaMask) {
						cpu.instructions[i] = info.instr
						counter++
						break
					}
				}
			}
		}
	}
	log.Printf("added %d cpu instructions", counter)

	return cpu
}
