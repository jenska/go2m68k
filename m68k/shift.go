package m68k

/*
//shift/rotate instruction
func xsx(cpu *M68K, opcode uint16, left, memory, arithmetic bool) {
	var data uint32
	if memory {
		data = cpu.loadEA(Word, opcode & 0x3F)
        prefetch()
		sign := 0x8000 & data
		carry = data
		if left {

		}
		bool carry = left ? (0x8000 & data) != 0 : data & 1;
		data = left ? data << 1 : data >> 1;
		if (arithmetic && !left) data |= sign;
		u16 sign2 = 0x8000 & data;
        setFlags(flag_logical, SizeWord, data);
		reg_s.x = reg_s.c = carry;
		if (arithmetic && left) reg_s.v = sign != sign2;
	} else {
        prefetch();
        size = (opcode >> 6) & 3;
		bool ir = (opcode >> 5) & 1;
		u8 multi = (opcode >> 9) & 7;
		u8 regPos = opcode & 7;
        u8 shiftCount = !ir ? ( multi == 0 ? 8 : multi ) : reg_d[multi] & 63;
        data = LoadEA(size, regPos); //register direct
		switch (size) {
            case SizeByte: left ? shift_left<arithmetic, SizeByte>(shiftCount, data)
                                : shift_right<arithmetic, SizeByte>(shiftCount, data); break;
            case SizeWord: left ? shift_left<arithmetic, SizeWord>(shiftCount, data)
                                : shift_right<arithmetic, SizeWord>(shiftCount, data); break;
            case SizeLong: left ? shift_left<arithmetic, SizeLong>(shiftCount, data)
                                : shift_right<arithmetic, SizeLong>(shiftCount, data); break;
		}
		sync(2 + (size == SizeLong ? 2 : 0) + shiftCount * 2);
	}
    writeEA(size, data, true);
}
*/
