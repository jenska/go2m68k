package cpu

import "log"

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------

func bra8(c *M68K) {
	c.pc += int32(int8(c.ir))
}

func bra16(c *M68K) {
	c.pc += c.popPC(Word)
}

func dbra(c *M68K) {
	reg := c.ir & 7
	count := int16(c.d[reg]) - 1
	Word.set(int32(count), &c.d[reg])
	if count != -1 {
		c.pc += c.popPC(Word) - 2
	} else {
		c.pc += Word.size
	}
}

func moveq(c *M68K) {
	reg := (c.ir >> 9) & 7
	data := int32(int8(c.ir))
	c.d[reg] = data
	c.sr.setLogicalFlags(data)
}

func stop(c *M68K) {
	if c.sr.S {
		c.sr.setbits(c.popPC(Word))
		c.stopped = true
	} else {
		panic(PrivilegeViolationError)
	}
}

func nop(c *M68K) {}

//------------------------------------------------------------------------------

// SetISA68000 Instruction Set Archtiecture for M68000
func (cpu *M68K) SetISA68000() Builder {
	// Address Bus Mask for 68000 CPU
	const busMask = 0x00ffffff
	c := cpu

	cpu.read = func(a int32, s *Size) int32 {
		if a&1 == 1 && s != Byte {
			panic(AdressError)
		}
		return c.bus.read(a&busMask, s)
	}

	cpu.write = func(a int32, s *Size, value int32) {
		if a&1 == 1 && s != Byte {
			panic(AdressError)
		}
		if !c.sr.S && a < 0x800 {
			panic(PrivilegeViolationError)
		}
		c.bus.write(a&busMask, s, value)
	}

	// // TODO: reusable method for other cpu types
	counter := 0
	for i := range dasmTable {
		opcode := uint16(i)
		// todo: cpu.opcodeTable[i] = dasmIllegal
		for _, info := range opcodeTable {
			if info.cpu == M68000 {
				if (opcode & info.mask) == info.match {
					if info.move && !validEA(eam(opcode), 0xbf8) {
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
