package m68k

type shiftFunc func(cpu *M68K, o *operand, shiftCount uint32, data *uint32)

func shiftLeft(cpu *M68K, o *operand, shiftCount uint32, data *uint32) {
	sr := cpu.SR
	carry := false
	if shiftCount >= o.Bits {
		if shiftCount == o.Bits {
			carry = (*data & 1) != 0
		} else {
			carry = false
		}
	} else if shiftCount > 0 {
		*data <<= shiftCount - 1
		carry = (*data & o.Msb) != 0
		*data <<= 1
	}
	sr.setFlags(flagLogical, o, *data, 0, 0)
	sr.C = carry
	if shiftCount != 0 {
		sr.X = sr.C
	}
}

func shiftLeftArithmetic(cpu *M68K, o *operand, shiftCount uint32, data *uint32) {
	sr := cpu.SR
	carry, overflow := false, false
	if shiftCount >= o.Bits {
		overflow = *data != 0
		if shiftCount == o.Bits {
			carry = (*data & 1) != 0
		} else {
			carry = false
		}
	} else if shiftCount > 0 {
		mask := o.Mask << (o.Bits - 1 - shiftCount) & o.Mask
		overflow = (*data&mask) != mask && (*data&mask) != 0
		*data <<= shiftCount - 1
		carry = (*data & o.Msb) != 0
		*data <<= 1
	}
	sr.setFlags(flagLogical, o, *data, 0, 0)
	sr.C = carry
	sr.V = overflow
	if shiftCount != 0 {
		sr.X = sr.C
	}
}

func shiftRight(cpu *M68K, o *operand, shiftCount uint32, data *uint32) {
	sr := cpu.SR
	carry, sign := false, (*data&o.Msb) != 0
	if shiftCount > o.Bits {
		*data = 0
		if shiftCount == o.Bits {
			carry = sign
		} else {
			carry = false
		}
	} else if shiftCount > 0 {
		*data >>= shiftCount - 1
		carry = (*data & 1) != 0
		*data >>= 1
	}

	sr.setFlags(flagLogical, o, *data, 0, 0)
	sr.C = carry
	if shiftCount != 0 {
		sr.X = sr.C
	}
}

func shiftRightArithmetic(cpu *M68K, o *operand, shiftCount uint32, data *uint32) {
	sr := cpu.SR
	carry, sign := false, (*data&o.Msb) != 0
	if shiftCount > o.Bits {
		if sign {
			*data = o.Mask
		} else {
			*data = 0
		}
		carry = sign
	} else if shiftCount > 0 {
		*data >>= shiftCount - 1
		carry = (*data & 1) != 0
		*data >>= 1
		mask := o.Mask << (o.Bits - shiftCount)
		if sign {
			mask &= o.Mask
		}
		*data |= mask
	}

	sr.setFlags(flagLogical, o, *data, 0, 0)
	sr.C = carry
	if shiftCount != 0 {
		sr.X = sr.C
	}
}

func (cpu *M68K) shift(shifter shiftFunc, opcode uint16) {

}
