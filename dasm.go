package cpu

import (
	"fmt"
	"sort"
	"strconv"
)

const (
	dasm68000 = 1
	dasm68010 = 2
	dasm68020 = 4
	dasm68030 = 8
	dasm68040 = 16
	dasmAll   = dasm68000 | dasm68010 | dasm68020 | dasm68030 | dasm68040
)

type (
	dasmOpcode struct {
		dasmHandler func(*dasmInfo) string
		mask        uint16
		match       uint16
		eaMask      uint16
		cpuTypes    uint16
	}

	dasmInfo struct {
		cpuType uint16
		bus     AddressBus
		pc      int32
		ir      uint16
		helper  string
	}
)

var (
	dasmTable          = make([]func(*dasmInfo) string, 0x10000)
	dasmSupportedTypes = make([]uint16, 0x10000)

	g3bitQDataTable     = []int32{8, 1, 2, 3, 4, 5, 6, 7}
	g5bitQDataTable     = []int32{32, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	gCCTable            = []string{"t", "f", "hi", "ls", "cc", "cs", "ne", "eq", "vc", "vs", "pl", "mi", "ge", "lt", "gt", "le"}
	gCPCCTable          = []string{"f", "eq", "ogt", "oge", "olt", "ole", "ogl", "or", "un", "ueq", "ugt", "uge", "ult", "ule", "ne", "t", "sf", "seq", "gt", "ge", "lt", "le", "glgle", "ngle", "ngl", "nle", "nlt", "nge", "ngt", "sne", "st", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"}
	gMMURegs            = []string{"tc", "drp", "srp", "crp", "cal", "val", "sccr", "acr"}
	gMMUCond            = []string{"bs", "bc", "ls", "lc", "ss", "sc", "as", "ac", "ws", "wc", "is", "ic", "gs", "gc", "cs", "cc"}
	gFPUDataFormatTable = []string{".l", ".s", ".x", ".p", ".w", ".d", ".b", ".p"}
)

func (d *dasmInfo) readImm(advance *Size) int32 {
	result := d.bus.read(d.pc, advance)
	d.pc += advance.align
	return result
}

func (d *dasmInfo) getImmStrSigned(size *Size) string {
	return size.SignedHexString(d.readImm(size))
}

func (d *dasmInfo) getImmStrUnsigned(size *Size) string {
	return size.HexString(d.readImm(size))
}

func (d *dasmInfo) getEaModeStr(size *Size) string {
	// Make string of effective address mode
	switch d.ir & 0x3f {
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07: // data register direct
		return fmt.Sprintf("D%d", y(d.ir))
	case 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f: // address register direct
		return fmt.Sprintf("A%d", y(d.ir))
	case 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17: // address register indirect
		return fmt.Sprintf("(A%d)", y(d.ir))
	case 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f: // address register indirect with postincrement
		return fmt.Sprintf("(A%d)+", y(d.ir))
	case 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27: // address register indirect with predecrement
		return fmt.Sprintf("-(A%d)", y(d.ir))
	case 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f: // address register indirect with displacement
		return fmt.Sprintf("(%s,A%d)", Word.SignedHexString(d.readImm(Word)), y(d.ir))
	case 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37:
		// address register indirect with index
		mode := ""
		extension := d.readImm(Word)
		if extension&0x100 != 0 {
			if extension&0xE4 == 0xC4 || extension&0xE2 == 0xC0 {
				return ""
			}
			base := int32(0)
			if extension&0x30 > 16 {
				if extension&0x30 == 0x30 {
					base = d.readImm(Long)
				} else {
					base = d.readImm(Word)
				}
			}
			outer := int32(0)
			if extension&3 > 1 && extension&71 < 68 {
				if extension&3 == 3 && extension&71 < 68 {
					outer = d.readImm(Long)
				} else {
					outer = d.readImm(Word)
				}
			}
			baseReg, indexReg := "", ""
			if extension&128 == 0 {
				baseReg = fmt.Sprintf("A%d", y(d.ir))
			}
			if extension&64 == 0 {
				indexReg = "D"
				if extension&0x8000 != 0 {
					indexReg = "A"
				}
				indexReg += strconv.Itoa(int((extension >> 12) & 7))
				if extension&0x800 != 0 {
					indexReg += Long.ext
				} else {
					indexReg += Word.ext
				}
				if extension>>uint64(9)&3 != 0 {
					indexReg += "*" + strconv.Itoa(int(1<<((extension>>9)&3)))
				}
			}
			preindex := extension&7 > 0 && extension&7 < 4
			postindex := extension&7 > 4
			comma := false
			if base != 0 {
				if extension&0x30 == 0x30 {
					mode += Long.SignedHexString(base)
				} else {
					mode += Word.SignedHexString(base)
				}
				comma = true
			}
			if baseReg[0] != 0 {
				mode += baseReg
				comma = true
			}
			if postindex {
				mode += "]"
				comma = true
			}
			if indexReg[0] != 0 {
				if comma {
					mode += ","
				}
				mode += indexReg
				comma = true
			}
			if preindex {
				mode += "]"
				comma = true
			}
			if outer != 0 {
				if comma {
					mode += ","
				}
				mode += Word.SignedHexString(outer)
			}
			return mode
		}
		if extension&0xFF == 0 {
			mode = fmt.Sprintf("(A%d,", y(d.ir))
			if extension&0x8000 != 0 {
				mode += "A"
			} else {
				mode += "D"
			}
			mode += strconv.Itoa(int((extension >> 12) & 7))
			if extension&0x800 != 0 {
				mode += Long.ext
			} else {
				mode += Word.ext
			}
		} else {
			mode = fmt.Sprintf("(%s,A%d,", Byte.SignedHexString(extension), y(d.ir))
			if extension&0x8000 != 0 {
				mode += "A"
			} else {
				mode += "D"
			}
			mode += strconv.Itoa(int((extension >> 12) & 7))
			if extension&0x800 != 0 {
				mode += Long.ext
			} else {
				mode += Word.ext
			}
		}
		if (extension>>9)&3 != 0 {
			mode += fmt.Sprintf("*%d", 1<<uint64(extension>>uint64(9)&3))
		}
		return mode
	case 56:
		// absolute short address
		return fmt.Sprintf("$%x.w", d.readImm(Word))
	case 57:
		// absolute long address
		return fmt.Sprintf("$%x.l", d.readImm(Long))
	case 58:
		// program counter with displacement
		tempValue := d.readImm(Word)
		d.helper = fmt.Sprintf("; ($%x)", tempValue+d.pc-2)
		return fmt.Sprintf("(%s,PC)", Word.SignedHexString(tempValue))
	case 59:
		// program counter with index
		mode := ""
		extension := d.readImm(Word)
		if extension&0x100 != 0 {
			if extension&0xE4 == 0xC4 || extension&0xE2 == 0xC0 {
				return ""
			}
			base := int32(0)
			if extension&0x30 > 0x10 {
				if extension&0x30 == 0x30 {
					base = d.readImm(Long)
				} else {
					base = d.readImm(Word)
				}
			}
			outer := int32(0)
			if extension&3 > 1 && extension&71 < 68 {
				if extension&3 == 3 && extension&71 < 68 {
					outer = d.readImm(Long)
				} else {
					outer = d.readImm(Word)
				}
			}
			baseReg := ""
			if extension&128 == 0 {
				baseReg = "PC"
			}
			indexReg := ""
			if extension&64 == 0 {
				indexReg = "D"
				if extension&0x8000 != 0 {
					indexReg = "A"
				}
				indexReg += strconv.Itoa(int((extension >> 12) & 7))
				if extension&0x800 != 0 {
					indexReg += Long.ext
				} else {
					indexReg += Word.ext
				}
				if extension>>uint64(9)&3 != 0 {
					indexReg += "*" + strconv.Itoa(int(1<<((extension>>9)&3)))
				}
			}
			preindex := extension&7 > 0 && extension&7 < 4
			postindex := extension&7 > 4
			comma := false
			mode = "("
			if preindex || postindex {
				mode += "["
			}
			if base != 0 {
				mode += Word.SignedHexString(base)
				comma = true
			}
			if baseReg != "" {
				if comma {
					mode += ","
				}
				mode += baseReg
				comma = true
			}
			if postindex {
				mode += "]"
				comma = true
			}
			if indexReg != "" {
				if comma {
					mode += ","
				}
				mode += indexReg
				comma = true
			}
			if preindex {
				mode += "]"
				comma = true
			}
			if outer != 0 {
				if comma {
					mode += ","
				}
				mode += Word.SignedHexString(outer)
			}
			return mode
		}
		if extension&0xFF == 0 {
			mode = "(PC,"
			if extension&0x8000 != 0 {
				mode += "A"
			} else {
				mode += "D"
			}
			mode += strconv.Itoa(int((extension >> 12) & 7))
			if extension&0x800 != 0 {
				mode += Long.ext
			} else {
				mode += Word.ext
			}
		} else {
			mode = fmt.Sprintf("(%s,PC,", Byte.SignedHexString(extension))
			if extension&0x8000 != 0 {
				mode += "A"
			} else {
				mode += "D"
			}
			mode += strconv.Itoa(int((extension >> 12) & 7))
			if extension&0x800 != 0 {
				mode += Long.ext
			} else {
				mode += Word.ext
			}
		}
		if (extension>>9)&3 != 0 {
			mode += fmt.Sprintf("*%d", 1<<uint64(extension>>uint64(9)&3))
		}
		return mode
	case 60:
		// Immediate
		return d.getImmStrUnsigned(size)
	default:
		return fmt.Sprintf("INVALID %x", d.ir&0x3f)
	}
}

func d68000Illegal(d *dasmInfo) string {
	return fmt.Sprintf("dc.w    $%04x; ILLEGAL", d.ir)
}

func d68000LineA(d *dasmInfo) string {
	return fmt.Sprintf("dc.w    $%04x; opcode 1010", d.ir)
}

func d68000LineF(d *dasmInfo) string {
	return fmt.Sprintf("dc.w    $%04x; opcode 1111", d.ir)
}

func d68000_abcd_rr(d *dasmInfo) string {
	return fmt.Sprintf("abcd    D%d, D%d", y(d.ir), x(d.ir))
}

func d68000_abcd_mm(d *dasmInfo) string {
	return fmt.Sprintf("abcd    -(A%d), -(A%d)", y(d.ir), x(d.ir))
}

func d68000_add_er_8(d *dasmInfo) string {
	return fmt.Sprintf("add.b   %s, D%d", d.getEaModeStr(Byte), x(d.ir))
}

func d68000_add_er_16(d *dasmInfo) string {
	return fmt.Sprintf("add.w   %s, D%d", d.getEaModeStr(Word), x(d.ir))
}

func d68000_add_er_32(d *dasmInfo) string {
	return fmt.Sprintf("add.l   %s, D%d", d.getEaModeStr(Long), x(d.ir))
}

func d68000_add_re_8(d *dasmInfo) string {
	return fmt.Sprintf("add.b   D%d, %s", x(d.ir), d.getEaModeStr(Byte))
}

func d68000_add_re_16(d *dasmInfo) string {
	return fmt.Sprintf("add.w   D%d, %s", x(d.ir), d.getEaModeStr(Word))
}

func d68000_add_re_32(d *dasmInfo) string {
	return fmt.Sprintf("add.l   D%d, %s", x(d.ir), d.getEaModeStr(Long))
}
func d68000_adda_16(d *dasmInfo) string {
	return fmt.Sprintf("adda.w  %s, A%d", d.getEaModeStr(Word), x(d.ir))
}

func d68000_adda_32(d *dasmInfo) string {
	return fmt.Sprintf("adda.l  %s, A%d", d.getEaModeStr(Long), x(d.ir))
}

func d68000_addi_8(d *dasmInfo) string {
	return fmt.Sprintf("addi.b  %s, %s", d.getImmStrSigned(Byte), d.getEaModeStr(Byte))
}

func d68000_addi_16(d *dasmInfo) string {
	return fmt.Sprintf("addi.w  %s, %s", d.getImmStrSigned(Word), d.getEaModeStr(Word))
}

func d68000_addi_32(d *dasmInfo) string {
	return fmt.Sprintf("addi.l  %s, %s", d.getImmStrSigned(Long), d.getEaModeStr(Long))
}

func d68000_addq_8(d *dasmInfo) string {
	return fmt.Sprintf("addq.b  #%d, %s", g3bitQDataTable[x(d.ir)], d.getEaModeStr(Byte))
}

func d68000_addq_16(d *dasmInfo) string {
	return fmt.Sprintf("addq.w  #%d, %s", g3bitQDataTable[x(d.ir)], d.getEaModeStr(Word))
}

func d68000_addq_32(d *dasmInfo) string {
	return fmt.Sprintf("addq.l  #%d, %s", g3bitQDataTable[x(d.ir)], d.getEaModeStr(Long))
}

func d68000_addx_rr_8(d *dasmInfo) string {
	return fmt.Sprintf("addx.b  D%d, D%d", y(d.ir), x(d.ir))
}

func d68000_addx_rr_16(d *dasmInfo) string {
	return fmt.Sprintf("addx.w  D%d, D%d", y(d.ir), x(d.ir))
}

func d68000_addx_rr_32(d *dasmInfo) string {
	return fmt.Sprintf("addx.l  D%d, D%d", y(d.ir), x(d.ir))
}

func d68000_addx_mm_8(d *dasmInfo) string {
	return fmt.Sprintf("addx.b  -(A%d), -(A%d)", y(d.ir), x(d.ir))
}

func d68000_addx_mm_16(d *dasmInfo) string {
	return fmt.Sprintf("addx.w  -(A%d), -(A%d)", y(d.ir), x(d.ir))
}

func d68000_addx_mm_32(d *dasmInfo) string {
	return fmt.Sprintf("addx.l  -(A%d), -(A%d)", y(d.ir), x(d.ir))
}

func d68000_and_er_8(d *dasmInfo) string {
	return fmt.Sprintf("and.b   %s, D%d", d.getEaModeStr(Byte), x(d.ir))
}

func d68000_and_er_16(d *dasmInfo) string {
	return fmt.Sprintf("and.w   %s, D%d", d.getEaModeStr(Word), x(d.ir))
}

func d68000_and_er_32(d *dasmInfo) string {
	return fmt.Sprintf("and.l   %s, D%d", d.getEaModeStr(Long), x(d.ir))
}

func d68000_and_re_8(d *dasmInfo) string {
	return fmt.Sprintf("and.b   D%d, %s", x(d.ir), d.getEaModeStr(Byte))
}

func d68000_and_re_16(d *dasmInfo) string {
	return fmt.Sprintf("and.w   D%d, %s", x(d.ir), d.getEaModeStr(Word))
}

func d68000_and_re_32(d *dasmInfo) string {
	return fmt.Sprintf("and.l   D%d, %s", x(d.ir), d.getEaModeStr(Long))
}

func d68000_andi_8(d *dasmInfo) string {
	return fmt.Sprintf("andi.b  %s, %s", d.getImmStrUnsigned(Byte), d.getEaModeStr(Byte))
}

func d68000_andi_16(d *dasmInfo) string {
	return fmt.Sprintf("andi.w  %s, %s", d.getImmStrUnsigned(Word), d.getEaModeStr(Word))
}

func d68000_andi_32(d *dasmInfo) string {
	return fmt.Sprintf("andi.l  %s, %s", d.getImmStrUnsigned(Long), d.getEaModeStr(Long))
}

func d68000_andi_to_ccr(d *dasmInfo) string {
	return fmt.Sprintf("andi    %s, CCR", d.getImmStrUnsigned(Byte))
}

func d68000_andi_to_sr(d *dasmInfo) string {
	return fmt.Sprintf("andi    %s, SR", d.getImmStrUnsigned(Word))
}

func d68000_asr_s_8(d *dasmInfo) string {
	return fmt.Sprintf("asr.b   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_asr_s_16(d *dasmInfo) string {
	return fmt.Sprintf("asr.w   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_asr_s_32(d *dasmInfo) string {
	return fmt.Sprintf("asr.l   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_asr_r_8(d *dasmInfo) string {
	return fmt.Sprintf("asr.b   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_asr_r_16(d *dasmInfo) string {
	return fmt.Sprintf("asr.w   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_asr_r_32(d *dasmInfo) string {
	return fmt.Sprintf("asr.l   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_asr_ea(d *dasmInfo) string {
	return fmt.Sprintf("asr.w   %s", d.getEaModeStr(Word))
}

func d68000_asl_s_8(d *dasmInfo) string {
	return fmt.Sprintf("asl.b   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_asl_s_16(d *dasmInfo) string {
	return fmt.Sprintf("asl.w   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_asl_s_32(d *dasmInfo) string {
	return fmt.Sprintf("asl.l   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_asl_r_8(d *dasmInfo) string {
	return fmt.Sprintf("asl.b   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_asl_r_16(d *dasmInfo) string {
	return fmt.Sprintf("asl.w   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_asl_r_32(d *dasmInfo) string {
	return fmt.Sprintf("asl.l   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_asl_ea(d *dasmInfo) string {
	return fmt.Sprintf("asl.w   %s", d.getEaModeStr(Word))
}

func d68000_bcc_8(d *dasmInfo) string {
	return fmt.Sprintf("b%-2s     $%x", gCCTable[(d.ir>>8)&15], d.pc+int32(int8(d.ir)))
}

func d68000_bcc_16(d *dasmInfo) string {
	return fmt.Sprintf("b%-2s     $%x", gCCTable[d.ir>>uint64(8)&15], d.pc+d.readImm(Word))
}

func d68020_bcc_32(d *dasmInfo) string {
	return fmt.Sprintf("b%-2s     $%x; (2+)", gCCTable[d.ir>>uint64(8)&15], d.pc+d.readImm(Long))
}

func d68000_bchg_r(d *dasmInfo) string {
	return fmt.Sprintf("bchg    D%d, %s", x(d.ir), d.getEaModeStr(Byte))
}

func d68000_bchg_s(d *dasmInfo) string {
	return fmt.Sprintf("bchg    %s, %s", d.getImmStrUnsigned(Byte), d.getEaModeStr(Byte))
}

func d68000_bclr_r(d *dasmInfo) string {
	return fmt.Sprintf("bclr    D%d, %s", x(d.ir), d.getEaModeStr(Byte))
}

func d68000_bclr_s(d *dasmInfo) string {
	return fmt.Sprintf("bclr    %s, %s", d.getImmStrUnsigned(Byte), d.getEaModeStr(Byte))
}

func d68010_bkpt(d *dasmInfo) string {
	return fmt.Sprintf("bkpt #%d; (1+)", y(d.ir))
}

func d68020_bfchg(d *dasmInfo) string {
	extension := d.readImm(Word)
	offset := ""
	width := ""
	if extension&0x800 != 0 {
		offset = fmt.Sprintf("D%d", (extension>>6)&7)
	} else {
		offset = fmt.Sprintf("%d", (extension>>6)&31)
	}
	if extension&32 != 0 {
		width = fmt.Sprintf("D%d", extension&7)
	} else {
		width = fmt.Sprintf("%d", g5bitQDataTable[extension&31])
	}
	return fmt.Sprintf("bfchg   %s {%s:%s}; (2+)", d.getEaModeStr(Byte), offset, width)
}

func d68020_bfclr(d *dasmInfo) string {
	extension := d.readImm(Word)
	offset := ""
	width := ""
	if extension&0x800 != 0 {
		offset = fmt.Sprintf("D%d", (extension>>6)&7)
	} else {
		offset = fmt.Sprintf("%d", (extension>>6)&31)
	}
	if extension&32 != 0 {
		width = fmt.Sprintf("D%d", extension&7)
	} else {
		width = fmt.Sprintf("%d", g5bitQDataTable[extension&31])
	}
	return fmt.Sprintf("bfclr   %s {%s:%s}; (2+)", d.getEaModeStr(Byte), offset, width)
}

func d68020_bfexts(d *dasmInfo) string {
	extension := d.readImm(Word)
	offset := ""
	width := ""
	if extension&0x800 != 0 {
		offset = fmt.Sprintf("D%d", (extension>>6)&7)
	} else {
		offset = fmt.Sprintf("%d", (extension>>6)&31)
	}
	if extension&32 != 0 {
		width = fmt.Sprintf("D%d", extension&7)
	} else {
		width = fmt.Sprintf("%d", g5bitQDataTable[extension&31])
	}
	return fmt.Sprintf("bfexts   %s {%s:%s}; (2+)", d.getEaModeStr(Byte), offset, width)
}

func d68020_bfextu(d *dasmInfo) string {
	extension := d.readImm(Word)
	offset := ""
	width := ""
	if extension&0x800 != 0 {
		offset = fmt.Sprintf("D%d", (extension>>6)&7)
	} else {
		offset = fmt.Sprintf("%d", (extension>>6)&31)
	}
	if extension&32 != 0 {
		width = fmt.Sprintf("D%d", extension&7)
	} else {
		width = fmt.Sprintf("%d", g5bitQDataTable[extension&31])
	}
	return fmt.Sprintf("bfextu   %s {%s:%s}; (2+)", d.getEaModeStr(Byte), offset, width)
}

func d68020_bfffo(d *dasmInfo) string {
	extension := d.readImm(Word)
	offset := ""
	width := ""
	if extension&0x800 != 0 {
		offset = fmt.Sprintf("D%d", (extension>>6)&7)
	} else {
		offset = fmt.Sprintf("%d", (extension>>6)&31)
	}
	if extension&32 != 0 {
		width = fmt.Sprintf("D%d", extension&7)
	} else {
		width = fmt.Sprintf("%d", g5bitQDataTable[extension&31])
	}
	return fmt.Sprintf("bfffo   %s {%s:%s}; (2+)", d.getEaModeStr(Byte), offset, width)
}

func d68020_bfins(d *dasmInfo) string {
	extension := d.readImm(Word)
	offset := ""
	width := ""
	if extension&0x800 != 0 {
		offset = fmt.Sprintf("D%d", (extension>>6)&7)
	} else {
		offset = fmt.Sprintf("%d", (extension>>6)&31)
	}
	if extension&32 != 0 {
		width = fmt.Sprintf("D%d", extension&7)
	} else {
		width = fmt.Sprintf("%d", g5bitQDataTable[extension&31])
	}
	return fmt.Sprintf("bfins   %s {%s:%s}; (2+)", d.getEaModeStr(Byte), offset, width)
}

func d68020_bfset(d *dasmInfo) string {
	extension := d.readImm(Word)
	offset := ""
	width := ""
	if extension&0x800 != 0 {
		offset = fmt.Sprintf("D%d", (extension>>6)&7)
	} else {
		offset = fmt.Sprintf("%d", (extension>>6)&31)
	}
	if extension&32 != 0 {
		width = fmt.Sprintf("D%d", extension&7)
	} else {
		width = fmt.Sprintf("%d", g5bitQDataTable[extension&31])
	}
	return fmt.Sprintf("bfset  %s {%s:%s}; (2+)", d.getEaModeStr(Byte), offset, width)
}

func d68020_bftst(d *dasmInfo) string {
	var offset string
	var width string
	extension := d.readImm(Word)
	if extension&0x800 != 0 {
		offset = fmt.Sprintf("D%d", (extension>>6)&7)
	} else {
		offset = fmt.Sprintf("%d", (extension>>6)&31)
	}
	if extension&32 != 0 {
		width = fmt.Sprintf("D%d", extension&7)
	} else {
		width = fmt.Sprintf("%d", g5bitQDataTable[extension&31])
	}
	return fmt.Sprintf("bftst   %s {%s:%s}; (2+)", d.getEaModeStr(Byte), offset, width)
}

func d68000_bra_8(d *dasmInfo) string {
	return fmt.Sprintf("bra     $%x", d.pc+int32(int8(d.ir)))
}

func d68000_bra_16(d *dasmInfo) string {
	return fmt.Sprintf("bra     $%x", d.pc+int32(int16(d.readImm(Word))))
}

func d68020_bra_32(d *dasmInfo) string {
	return fmt.Sprintf("bra     $%x; (2+)", d.pc+d.readImm(Long))
}

func d68000_bset_r(d *dasmInfo) string {
	return fmt.Sprintf("bset    D%d, %s", x(d.ir), d.getEaModeStr(Byte))
}

func d68000_bset_s(d *dasmInfo) string {
	return fmt.Sprintf("bset    %s, %s", d.getImmStrUnsigned(Byte), d.getEaModeStr(Byte))
}

func d68000_bsr_8(d *dasmInfo) string {
	return fmt.Sprintf("bsr     $%x", d.pc+int32(int8(d.ir)))
}

func d68000_bsr_16(d *dasmInfo) string {
	return fmt.Sprintf("bsr     $%x", d.pc+int32(int16(d.readImm(Word))))
}

func d68020_bsr_32(d *dasmInfo) string {
	return fmt.Sprintf("bsr     $%x; (2+)", d.pc+d.readImm(Long))
}

func d68000_btst_r(d *dasmInfo) string {
	return fmt.Sprintf("btst    D%d, %s", x(d.ir), d.getEaModeStr(Byte))
}

func d68000_btst_s(d *dasmInfo) string {
	return fmt.Sprintf("btst    %s, %s", d.getImmStrUnsigned(Byte), d.getEaModeStr(Byte))
}

func d68020_callm(d *dasmInfo) string {
	return fmt.Sprintf("callm   %s, %s; (2)", d.getImmStrUnsigned(Byte), d.getEaModeStr(Byte))
}

func d68020_cas_8(d *dasmInfo) string {
	extension := d.readImm(Word)
	return fmt.Sprintf("cas.b   D%d, D%d, %s; (2+)", extension&7, extension>>uint64(6)&7, d.getEaModeStr(Byte))
}

func d68020_cas_16(d *dasmInfo) string {
	extension := d.readImm(Word)
	return fmt.Sprintf("cas.w   D%d, D%d, %s; (2+)", extension&7, extension>>uint64(6)&7, d.getEaModeStr(Word))
}

func d68020_cas_32(d *dasmInfo) string {
	extension := d.readImm(Word)
	return fmt.Sprintf("cas.l   D%d, D%d, %s; (2+)", extension&7, (extension>>6)&7, d.getEaModeStr(Long))
}

func d68020_cas2_16(d *dasmInfo) string {
	extension := d.readImm(Long)
	return fmt.Sprintf("cas2.w  D%d:D%d, D%d:D%d, (%c%d):(%c%d); (2+)", (extension>>16)&7, extension&7, (extension>>22)&7, (extension>>6)&7,
		func() int32 {
			if Long.IsNegative(extension) {
				return int32('A')
			}
			return int32('D')
		}(), extension>>uint64(28)&7, func() int32 {
			if extension&0x8000 != 0 {
				return int32('A')
			}
			return int32('D')
		}(), extension>>uint64(12)&7)
}

func d68020_cas2_32(d *dasmInfo) string {
	extension := d.readImm(Long)
	return fmt.Sprintf("cas2.l  D%d:D%d, D%d:D%d, (%c%d):(%c%d); (2+)", extension>>uint64(16)&7, extension&7, extension>>uint64(22)&7, extension>>uint64(6)&7, func() int32 {
		if Long.IsNegative(extension) {
			return int32('A')
		}
		return int32('D')
	}(), extension>>uint64(28)&7, func() int32 {
		if extension&0x8000 != 0 {
			return int32('A')
		}
		return int32('D')
	}(), extension>>uint64(12)&7)
}

func d68000_chk_16(d *dasmInfo) string {
	return fmt.Sprintf("chk.w   %s, D%d", d.getEaModeStr(Word), x(d.ir))
}

func d68020_chk_32(d *dasmInfo) string {
	return fmt.Sprintf("chk.l   %s, D%d; (2+)", d.getEaModeStr(Long), x(d.ir))
}

func d68020_chk2_cmp2_8(d *dasmInfo) string {
	extension := d.readImm(Word)
	return fmt.Sprintf("%s.b  %s, %c%d; (2+)", func() string {
		if extension&0x800 != 0 {
			return "chk2"
		}
		return "cmp2"
	}(), d.getEaModeStr(Byte), func() int32 {
		if extension&0x8000 != 0 {
			return int32('A')
		}
		return int32('D')
	}(), (extension>>12)&7)
}

func d68020_chk2_cmp2_16(d *dasmInfo) string {
	extension := d.readImm(Word)
	return fmt.Sprintf("%s.w  %s, %c%d; (2+)", func() string {
		if extension&0x800 != 0 {
			return "chk2"
		}
		return "cmp2"
	}(), d.getEaModeStr(Word), func() int32 {
		if extension&0x8000 != 0 {
			return int32('A')
		}
		return int32('D')
	}(), extension>>uint64(12)&7)
}

func d68020_chk2_cmp2_32(d *dasmInfo) string {
	extension := d.readImm(Word)
	return fmt.Sprintf("%s.l  %s, %c%d; (2+)", func() string {
		if extension&0x800 != 0 {
			return "chk2"
		}
		return "cmp2"
	}(), d.getEaModeStr(Long), func() int32 {
		if extension&0x8000 != 0 {
			return int32('A')
		}
		return int32('D')
	}(), (extension>>12)&7)
}

func d68040_cinv(d *dasmInfo) string {
	switch (d.ir >> 3) & 3 {
	case 0:
		return "cinv (illegal scope); (4)"
	case 1:
		return fmt.Sprintf("cinvl   %d, (A%d); (4)", (d.ir>>6)&3, y(d.ir))
	case 2:
		return fmt.Sprintf("cinvp   %d, (A%d); (4)", (d.ir>>6)&3, y(d.ir))
	default:
		return fmt.Sprintf("cinva   %d; (4)", (d.ir>>6)&3)
	}
}

func d68000_clr_8(d *dasmInfo) string {
	return fmt.Sprintf("clr.b   %s", d.getEaModeStr(Byte))
}

func d68000_clr_16(d *dasmInfo) string {
	return fmt.Sprintf("clr.w   %s", d.getEaModeStr(Word))
}

func d68000_clr_32(d *dasmInfo) string {
	return fmt.Sprintf("clr.l   %s", d.getEaModeStr(Long))
}

func d68000_cmp_8(d *dasmInfo) string {
	return fmt.Sprintf("cmp.b   %s, D%d", d.getEaModeStr(Byte), x(d.ir))
}

func d68000_cmp_16(d *dasmInfo) string {
	return fmt.Sprintf("cmp.w   %s, D%d", d.getEaModeStr(Word), x(d.ir))
}

func d68000_cmp_32(d *dasmInfo) string {
	return fmt.Sprintf("cmp.l   %s, D%d", d.getEaModeStr(Long), x(d.ir))
}

func d68000_cmpa_16(d *dasmInfo) string {
	return fmt.Sprintf("cmpa.w  %s, A%d", d.getEaModeStr(Word), x(d.ir))
}

func d68000_cmpa_32(d *dasmInfo) string {
	return fmt.Sprintf("cmpa.l  %s, A%d", d.getEaModeStr(Long), x(d.ir))
}

func d68000_cmpi_8(d *dasmInfo) string {
	return fmt.Sprintf("cmpi.b  %s, %s", d.getImmStrSigned(Byte), d.getEaModeStr(Byte))
}

func d68020_cmpi_pcdi_8(d *dasmInfo) string {
	return fmt.Sprintf("cmpi.b  %s, %s; (2+)", d.getImmStrSigned(Byte), d.getEaModeStr(Byte))
}

func d68020_cmpi_pcix_8(d *dasmInfo) string {
	return fmt.Sprintf("cmpi.b  %s, %s; (2+)", d.getImmStrSigned(Byte), d.getEaModeStr(Byte))
}

func d68000_cmpi_16(d *dasmInfo) string {
	return fmt.Sprintf("cmpi.w  %s, %s", d.getImmStrSigned(Word), d.getEaModeStr(Word))
}

func d68020_cmpi_pcdi_16(d *dasmInfo) string {
	return fmt.Sprintf("cmpi.w  %s, %s; (2+)", d.getImmStrSigned(Word), d.getEaModeStr(Word))
}

func d68020_cmpi_pcix_16(d *dasmInfo) string {
	return fmt.Sprintf("cmpi.w  %s, %s; (2+)", d.getImmStrSigned(Word), d.getEaModeStr(Word))
}

func d68000_cmpi_32(d *dasmInfo) string {
	return fmt.Sprintf("cmpi.l  %s, %s", d.getImmStrSigned(Long), d.getEaModeStr(Long))
}

func d68020_cmpi_pcdi_32(d *dasmInfo) string {
	return fmt.Sprintf("cmpi.l  %s, %s; (2+)", d.getImmStrSigned(Long), d.getEaModeStr(Long))
}

func d68020_cmpi_pcix_32(d *dasmInfo) string {
	return fmt.Sprintf("cmpi.l  %s, %s; (2+)", d.getImmStrSigned(Long), d.getEaModeStr(Long))
}

func d68000_cmpm_8(d *dasmInfo) string {
	return fmt.Sprintf("cmpm.b  (A%d)+, (A%d)+", y(d.ir), x(d.ir))
}

func d68000_cmpm_16(d *dasmInfo) string {
	return fmt.Sprintf("cmpm.w  (A%d)+, (A%d)+", y(d.ir), x(d.ir))
}

func d68000_cmpm_32(d *dasmInfo) string {
	return fmt.Sprintf("cmpm.l  (A%d)+, (A%d)+", y(d.ir), x(d.ir))
}

func d68020_cpbcc_16(d *dasmInfo) string {
	extension := d.readImm(Word)
	new_pc := d.pc + int32(int16(d.readImm(Word)))
	return fmt.Sprintf("%db%-4s  %s; %x (extension = %x) (2-3)", x(d.ir), gCPCCTable[d.ir&63], d.getImmStrSigned(Word), new_pc, extension)
}

func d68020_cpbcc_32(d *dasmInfo) string {
	extension := d.readImm(Word)
	new_pc := d.pc + d.readImm(Long)
	return fmt.Sprintf("%db%-4s  %s; %x (extension = %x) (2-3)", x(d.ir), gCPCCTable[d.ir&63], d.getImmStrSigned(Word), new_pc, extension)
}

func d68020_cpdbcc(d *dasmInfo) string {
	extension1 := d.readImm(Word)
	extension2 := d.readImm(Word)
	new_pc := d.pc + int32(int16(d.readImm(Word)))
	return fmt.Sprintf("%ddb%-4s D%d,%s; %x (extension = %x) (2-3)", x(d.ir), gCPCCTable[extension1&63], y(d.ir), d.getImmStrSigned(Word), new_pc, extension2)
}

func d68020_cpgen(d *dasmInfo) string {
	return fmt.Sprintf("%dgen    %s; (2-3)", x(d.ir), d.getImmStrUnsigned(Long))
}

func d68020_cprestore(d *dasmInfo) string {
	if x(d.ir) == 1 {
		return fmt.Sprintf("frestore %s", d.getEaModeStr(Byte))
	} else {
		return fmt.Sprintf("%drestore %s; (2-3)", x(d.ir), d.getEaModeStr(Byte))
	}
}

func d68020_cpsave(d *dasmInfo) string {
	if x(d.ir) == 1 {
		return fmt.Sprintf("fsave   %s", d.getEaModeStr(Byte))
	} else {
		return fmt.Sprintf("%dsave   %s; (2-3)", x(d.ir), d.getEaModeStr(Byte))
	}
}

func d68020_cpscc(d *dasmInfo) string {
	extension1 := d.readImm(Word)
	extension2 := d.readImm(Word)
	return fmt.Sprintf("%ds%-4s  %s; (extension = %x) (2-3)", x(d.ir), gCPCCTable[extension1&63], d.getEaModeStr(Byte), extension2)
}

func d68020_cptrapcc_0(d *dasmInfo) string {
	extension1 := d.readImm(Word)
	extension2 := d.readImm(Word)
	return fmt.Sprintf("%dtrap%-4s; (extension = %x) (2-3)", x(d.ir), gCPCCTable[extension1&63], extension2)
}

func d68020_cptrapcc_16(d *dasmInfo) string {
	extension1 := d.readImm(Word)
	extension2 := d.readImm(Word)
	return fmt.Sprintf("%dtrap%-4s %s; (extension = %x) (2-3)", x(d.ir), gCPCCTable[extension1&63], d.getImmStrUnsigned(Word), extension2)
}

func d68020_cptrapcc_32(d *dasmInfo) string {
	extension1 := d.readImm(Word)
	extension2 := d.readImm(Word)
	return fmt.Sprintf("%dtrap%-4s %s; (extension = %x) (2-3)", x(d.ir), gCPCCTable[extension1&63], d.getImmStrUnsigned(Long), extension2)
}

func d68040_cpush(d *dasmInfo) string {
	switch (d.ir >> 3) & 3 {
	case 0:
		return fmt.Sprintf("cpush (illegal scope); (4)")
	case 1:
		return fmt.Sprintf("cpushl  %d, (A%d); (4)", d.ir>>uint64(6)&3, y(d.ir))
	case 2:
		return fmt.Sprintf("cpushp  %d, (A%d); (4)", d.ir>>uint64(6)&3, y(d.ir))
	default:
		return fmt.Sprintf("cpusha  %d; (4)", d.ir>>uint64(6)&3)
	}
}

func d68000_dbra(d *dasmInfo) string {
	return fmt.Sprintf("dbra    D%d, $%x", y(d.ir), d.pc+int32(int16(d.readImm(Word))))
}

func d68000_dbcc(d *dasmInfo) string {
	return fmt.Sprintf("db%-2s    D%d, $%x", gCCTable[d.ir>>uint64(8)&15], y(d.ir), d.pc+int32(int16(d.readImm(Word))))
}

func d68000_divs(d *dasmInfo) string {
	return fmt.Sprintf("divs.w  %s, D%d", d.getEaModeStr(Word), x(d.ir))
}

func d68000_divu(d *dasmInfo) string {
	return fmt.Sprintf("divu.w  %s, D%d", d.getEaModeStr(Word), x(d.ir))
}

func d68020_divl(d *dasmInfo) string {
	extension := d.readImm(Word)
	if extension&1024 != 0 {
		return fmt.Sprintf("div%c.l  %s, D%d:D%d; (2+)", func() int32 {
			if extension&0x800 != 0 {
				return int32('s')
			}
			return int32('u')
		}(), d.getEaModeStr(Long), extension&7, extension>>uint64(12)&7)
	} else if extension&7 == extension>>uint64(12)&7 {
		return fmt.Sprintf("div%c.l  %s, D%d; (2+)", func() int32 {
			if extension&0x800 != 0 {
				return int32('s')
			}
			return int32('u')
		}(), d.getEaModeStr(Long), extension>>uint64(12)&7)
	} else {
		return fmt.Sprintf("div%cl.l %s, D%d:D%d; (2+)", func() int32 {
			if extension&0x800 != 0 {
				return int32('s')
			}
			return int32('u')
		}(), d.getEaModeStr(Long), extension&7, extension>>uint64(12)&7)
	}
}

func d68000_eor_8(d *dasmInfo) string {
	return fmt.Sprintf("eor.b   D%d, %s", x(d.ir), d.getEaModeStr(Byte))
}
func d68000_eor_16(d *dasmInfo) string {
	return fmt.Sprintf("eor.w   D%d, %s", x(d.ir), d.getEaModeStr(Word))
}
func d68000_eor_32(d *dasmInfo) string {
	return fmt.Sprintf("eor.l   D%d, %s", x(d.ir), d.getEaModeStr(Long))
}

func d68000_eori_8(d *dasmInfo) string {
	return fmt.Sprintf("eori.b  %s, %s", d.getImmStrUnsigned(Byte), d.getEaModeStr(Byte))
}

func d68000_eori_16(d *dasmInfo) string {
	return fmt.Sprintf("eori.w  %s, %s", d.getImmStrUnsigned(Word), d.getEaModeStr(Word))
}

func d68000_eori_32(d *dasmInfo) string {
	return fmt.Sprintf("eori.l  %s, %s", d.getImmStrUnsigned(Long), d.getEaModeStr(Long))
}

func d68000_eori_to_ccr(d *dasmInfo) string {
	return fmt.Sprintf("eori    %s, CCR", d.getImmStrUnsigned(Byte))
}

func d68000_eori_to_sr(d *dasmInfo) string {
	return fmt.Sprintf("eori    %s, SR", d.getImmStrUnsigned(Word))
}

func d68000_exg_dd(d *dasmInfo) string {
	return fmt.Sprintf("exg     D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_exg_aa(d *dasmInfo) string {
	return fmt.Sprintf("exg     A%d, A%d", x(d.ir), y(d.ir))
}

func d68000_exg_da(d *dasmInfo) string {
	return fmt.Sprintf("exg     D%d, A%d", x(d.ir), y(d.ir))
}

func d68000_ext_16(d *dasmInfo) string {
	return fmt.Sprintf("ext.w   D%d", y(d.ir))
}

func d68000_ext_32(d *dasmInfo) string {
	return fmt.Sprintf("ext.l   D%d", y(d.ir))
}

func d68020_extb_32(d *dasmInfo) string {
	return fmt.Sprintf("extb.l  D%d; (2+)", y(d.ir))
}

var fpu_ukn = "FPU (?)"
var float_mnemonics = map[int32]string{
	0: "fmove", 1: "fint", 2: "fsinh", 3: "fintrz",
	4: "fsqrt", 6: "flognp1", 8: "fetoxml", 9: "ftanh1",
	10: "fatan", 12: "fasin", 13: "fatanh", 14: "fsin",
	15: "ftan", 16: "fetox", 17: "ftwotox", 18: "ftentox",
	20: "flogn", 21: "flog10", 22: "flog2", 24: "fabs",
	25: "fcosh", 26: "fneg", 28: "facos", 29: "fcos",
	30: "fgetexp", 31: "fgetman", 32: "fdiv", 33: "fmod",
	34: "fadd", 35: "fmul", 36: "fsgldiv", 37: "frem",
	38: "fscale", 39: "fsglmul", 40: "fsub", 0x30: "fsincos",
	49: "fsincos", 50: "fsincos", 51: "fsincos", 52: "fsincos",
	53: "fsincos", 54: "fsincos", 55: "fsincos", 56: "fcmp",
	58: "ftst", 65: "fssqrt", 69: "fdsqrt", 88: "fsabs",
	90: "fsneg", 92: "fdabs", 94: "fdneg", 96: "fsdiv",
	98: "fsadd", 99: "fsmul", 100: "fddiv", 102: "fdadd",
	103: "fdmul", 104: "fssub", 108: "fdsub",
}

func d68040_fpu(d *dasmInfo) string {
	w2 := d.readImm(Word)
	src := (w2 >> 10) & 7
	destReg := (w2 >> 7) & 7
	if (w2>>13)&7 == 2 && (w2>>10)&7 == 7 {
		// special override for FMOVECR
		return fmt.Sprintf("fmovecr   #$%0x, fp%d", w2&127, destReg)
	}
	switch (w2 >> 13) & 0x07 {
	case 0, 2:
		if mnemonic, ok := float_mnemonics[w2&0x7f]; ok {
			if w2&0x4000 != 0 {
				return fmt.Sprintf("%s%s   %s, FP%d", mnemonic, gFPUDataFormatTable[src], d.getEaModeStr(Long), destReg)
			}
			return fmt.Sprintf("%s.x   FP%d, FP%d", mnemonic, src, destReg)
		}
		return fpu_ukn
	case 3:
		switch (w2 >> 10) & 0x07 {
		case 3:
			// packed decimal w/fixed k-factor
			return fmt.Sprintf("fmove%s   FP%d, %s {#%d}", gFPUDataFormatTable[w2>>uint64(10)&7], destReg, d.getEaModeStr(Long), int32(int8(w2&127)))
		case 7:
			// packed decimal w/dynamic k-factor (register)
			return fmt.Sprintf("fmove%s   FP%d, %s {D%d}", gFPUDataFormatTable[w2>>uint64(10)&7], destReg, d.getEaModeStr(Long), (w2>>4)&7)
		default:
			return fmt.Sprintf("fmove%s   FP%d, %s", gFPUDataFormatTable[w2>>uint64(10)&7], destReg, d.getEaModeStr(Long))
		}
	case 4:
		// ea to control
		if w2&4096 != 0 {
			return fmt.Sprintf("fmovem.l   %s, fpcr", d.getEaModeStr(Long))
		}
		if w2&0x800 != 0 {
			return fmt.Sprintf("fmovem.l   %s, /fpsr", d.getEaModeStr(Long))
		}
		if w2&1024 != 0 {
			return fmt.Sprintf("fmovem.l   %s, /fpiar", d.getEaModeStr(Long))
		}
	case 5:
		// control to ea
		if w2&4096 != 0 {
			return fmt.Sprintf("fmovem.l   fpcr, %s", d.getEaModeStr(Long))
		}
		if w2&0x800 != 0 {
			return fmt.Sprintf("fmovem.l   /fpsr, %s", d.getEaModeStr(Long))
		}
		if w2&1024 != 0 {
			return fmt.Sprintf("fmovem.l   /fpiar, %s", d.getEaModeStr(Long))
		}
	case 6:
		if (w2>>11)&1 != 0 {
			// memory to FPU, list
			// dynamic register list
			return fmt.Sprintf("fmovem.x   %s, D%d", d.getEaModeStr(Long), (w2>>4)&7)
		} else {
			// static register list
			str := fmt.Sprintf("fmovem.x   %s, ", d.getEaModeStr(Long))
			for i := 0; i < 8; i++ {
				if w2&(1<<i) != 0 {
					if (w2>>12)&1 != 0 {
						// postincrement or control
						str += fmt.Sprintf("FP%d ", 7-i)
					} else {
						// predecrement
						str += fmt.Sprintf("FP%d ", i)
					}
				}
			}
			return str
		}
	case 7:
		if (w2>>11)&1 != 0 {
			// FPU to memory, list
			// dynamic register list
			return fmt.Sprintf("fmovem.x   D%d, %s", w2>>uint64(4)&7, d.getEaModeStr(Long))
		} else {
			// static register list
			str := "fmovem.x   "
			for i := 0; i < 8; i++ {
				if w2&(1<<i) != 0 {
					if (w2>>12)&1 != 0 {
						// postincrement or control
						str += fmt.Sprintf("FP%d ", 7-i)
					} else {
						// predecrement
						str += fmt.Sprintf("FP%d ", i)
					}
				}
			}
			return str + ", " + d.getEaModeStr(Long)
		}
	}
	return fpu_ukn
}

func d68000_jmp(d *dasmInfo) string {
	return fmt.Sprintf("jmp     %s", d.getEaModeStr(Long))
}

func d68000_jsr(d *dasmInfo) string {
	return fmt.Sprintf("jsr     %s", d.getEaModeStr(Long))
}

func d68000_lea(d *dasmInfo) string {
	return fmt.Sprintf("lea     %s, A%d", d.getEaModeStr(Long), x(d.ir))
}

func d68000_link_16(d *dasmInfo) string {
	return fmt.Sprintf("link    A%d, %s", y(d.ir), d.getImmStrSigned(Word))
}

func d68020_link_32(d *dasmInfo) string {
	return fmt.Sprintf("link    A%d, %s; (2+)", y(d.ir), d.getImmStrSigned(Long))
}

func d68000_lsr_s_8(d *dasmInfo) string {
	return fmt.Sprintf("lsr.b   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_lsr_s_16(d *dasmInfo) string {
	return fmt.Sprintf("lsr.w   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_lsr_s_32(d *dasmInfo) string {
	return fmt.Sprintf("lsr.l   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_lsr_r_8(d *dasmInfo) string {
	return fmt.Sprintf("lsr.b   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_lsr_r_16(d *dasmInfo) string {
	return fmt.Sprintf("lsr.w   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_lsr_r_32(d *dasmInfo) string {
	return fmt.Sprintf("lsr.l   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_lsr_ea(d *dasmInfo) string {
	return fmt.Sprintf("lsr.w   %s", d.getEaModeStr(Long))
}

func d68000_lsl_s_8(d *dasmInfo) string {
	return fmt.Sprintf("lsl.b   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_lsl_s_16(d *dasmInfo) string {
	return fmt.Sprintf("lsl.w   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_lsl_s_32(d *dasmInfo) string {
	return fmt.Sprintf("lsl.l   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_lsl_r_8(d *dasmInfo) string {
	return fmt.Sprintf("lsl.b   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_lsl_r_16(d *dasmInfo) string {
	return fmt.Sprintf("lsl.w   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_lsl_r_32(d *dasmInfo) string {
	return fmt.Sprintf("lsl.l   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_lsl_ea(d *dasmInfo) string {
	return fmt.Sprintf("lsl.w   %s", d.getEaModeStr(Long))
}

func d68000_move_8(d *dasmInfo) string {
	d.ir = x(d.ir) | d.ir>>3&0x38
	return fmt.Sprintf("move.b  %s, %s", d.getEaModeStr(Byte), d.getEaModeStr(Byte))
}

func d68000_move_16(d *dasmInfo) string {
	d.ir = x(d.ir) | d.ir>>3&0x38
	return fmt.Sprintf("move.w  %s, %s", d.getEaModeStr(Word), d.getEaModeStr(Word))
}

func d68000_move_32(d *dasmInfo) string {
	d.ir = x(d.ir) | d.ir>>3&0x38
	return fmt.Sprintf("move.l  %s, %s", d.getEaModeStr(Long), d.getEaModeStr(Long))
}

func d68000_movea_16(d *dasmInfo) string {
	return fmt.Sprintf("movea.w %s, A%d", d.getEaModeStr(Word), x(d.ir))
}

func d68000_movea_32(d *dasmInfo) string {
	return fmt.Sprintf("movea.l %s, A%d", d.getEaModeStr(Long), x(d.ir))
}

func d68000_move_to_ccr(d *dasmInfo) string {
	return fmt.Sprintf("move    %s, CCR", d.getEaModeStr(Byte))
}

func d68010_move_fr_ccr(d *dasmInfo) string {
	return fmt.Sprintf("move    CCR, %s; (1+)", d.getEaModeStr(Byte))
}

func d68000_move_fr_sr(d *dasmInfo) string {
	return fmt.Sprintf("move    SR, %s", d.getEaModeStr(Word))
}

func d68000_move_to_sr(d *dasmInfo) string {
	return fmt.Sprintf("move    %s, SR", d.getEaModeStr(Word))
}

func d68000_move_fr_usp(d *dasmInfo) string {
	return fmt.Sprintf("move    USP, A%d", y(d.ir))
}

func d68000_move_to_usp(d *dasmInfo) string {
	return fmt.Sprintf("move    A%d, USP", y(d.ir))
}

func d68010_movec(d *dasmInfo) string {
	var reg_name, processor string
	extension := d.readImm(Word)
	switch extension & 4095 {
	case (0):
		reg_name = "SFC"
		processor = "1+"
	case (1):
		reg_name = "DFC"
		processor = "1+"
	case (0x800):
		reg_name = "USP"
		processor = "1+"
	case (2049):
		reg_name = "VBR"
		processor = "1+"
	case (2):
		reg_name = "CACR"
		processor = "2+"
	case (2050):
		reg_name = "CAAR"
		processor = "2,3"
	case (2051):
		reg_name = "MSP"
		processor = "2+"
	case (2052):
		reg_name = "ISP"
		processor = "2+"
	case (3):
		reg_name = "TC"
		processor = "4+"
	case (4):
		reg_name = "ITT0"
		processor = "4+"
	case (5):
		reg_name = "ITT1"
		processor = "4+"
	case (6):
		reg_name = "DTT0"
		processor = "4+"
	case (7):
		reg_name = "DTT1"
		processor = "4+"
	case 2053:
		reg_name = "MMUSR"
		processor = "4+"
	case (2054):
		reg_name = "URP"
		processor = "4+"
	case (2055):
		reg_name = "SRP"
		processor = "4+"
	default:
		reg_name = Word.SignedHexString(extension & 4095)
		processor = "?"
	}
	if d.ir&1 != 0 {
		return fmt.Sprintf("movec %c%d, %s; (%s)", func() int32 {
			if extension&0x8000 != 0 {
				return int32('A')
			}
			return int32('D')
		}(), extension>>uint64(12)&7, reg_name, processor)
	} else {
		return fmt.Sprintf("movec %s, %c%d; (%s)", reg_name, func() int32 {
			if extension&0x8000 != 0 {
				return int32('A')
			}
			return int32('D')
		}(), extension>>uint64(12)&7, processor)
	}
}

func d68000_movem_pd_16(d *dasmInfo) string {
	data := d.readImm(Word)
	buffer := ""
	for i := 0; i < 8; i++ {
		if data&(1<<(15-i)) != 0 {
			first := i
			rl := 0
			for i < 7 && data&(1<<(15-(i+1))) != 0 {
				i++
				rl++
			}
			if len(buffer) > 0 {
				buffer += "/"
			}
			buffer += fmt.Sprintf("D%d", first)
			if rl > 0 {
				buffer += fmt.Sprintf("-D%d", first+rl)
			}
		}
	}
	for i := 0; i < 8; i++ {
		if data&(1<<(7-i)) != 0 {
			first := i
			rl := 0
			for i < 7 && data&(1<<(7-(i+1))) != 0 {
				i++
				rl++
			}
			if len(buffer) > 0 {
				buffer += "/"
			}
			buffer += fmt.Sprintf("A%d", first)
			if rl > 0 {
				buffer += fmt.Sprintf("-A%d", first+rl)
			}
		}
	}
	return fmt.Sprintf("movem.w %s, %s", buffer, d.getEaModeStr(Word))
}

func d68000_movem_pd_32(d *dasmInfo) string {
	data := d.readImm(Word)
	buffer := ""

	for i := 0; i < 8; i++ {
		if data&(1<<(15-i)) != 0 {
			first := i
			rl := 0
			for i < 7 && data&(1<<(15-(i+1))) != 0 {
				i++
				rl++
			}
			if len(buffer) > 0 {
				buffer += "/"
			}
			buffer += fmt.Sprintf("D%d", first)
			if rl > 0 {
				buffer += fmt.Sprintf("-D%d", first+rl)
			}
		}
	}
	for i := 0; i < 8; i++ {
		if data&(1<<(7-i)) != 0 {
			first := i
			rl := 0
			for i < 7 && data&(1<<(7-(i+1))) != 0 {
				i++
				rl++
			}
			if len(buffer) > 0 {
				buffer += "/"
			}
			buffer += fmt.Sprintf("A%d", first)
			if rl > 0 {
				buffer += fmt.Sprintf("-A%d", first+rl)
			}
		}
	}
	return fmt.Sprintf("movem.l %s, %s", buffer, d.getEaModeStr(Long))
}

// d68000_movem_er_16 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:0xE27
func d68000_movem_er_16(d *dasmInfo) string {
	data := d.readImm(Word)
	buffer := ""
	for i := 0; i < 8; i++ {
		if data&(1<<i) != 0 {
			first := i
			rl := 0
			for i < 7 && data&(1<<(i+1)) != 0 {
				i++
				rl++
			}
			if len(buffer) > 0 {
				buffer += "/"
			}
			buffer += fmt.Sprintf("D%d", first)
			if rl > 0 {
				buffer += fmt.Sprintf("-D%d", first+rl)
			}
		}
	}
	for i := 0; i < 8; i++ {
		if data&(1<<(i+8)) != 0 {
			first := i
			rl := 0
			for i < 7 && data&(1<<(i+8+1)) != 0 {
				i++
				rl++
			}
			if len(buffer) > 0 {
				buffer += "/"
			}
			buffer += fmt.Sprintf("A%d", first)
			if rl > 0 {
				buffer += fmt.Sprintf("-A%d", first+rl)
			}
		}
	}
	return fmt.Sprintf("movem.w %s, %s", d.getEaModeStr(Word), buffer)
}

// d68000_movem_er_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2315
func d68000_movem_er_32(d *dasmInfo) string {
	data := d.readImm(Word)
	buffer := ""
	for i := 0; i < 8; i++ {
		if data&(1<<i) != 0 {
			first := i
			rl := 0
			for i < 7 && data&(1<<(i+1)) != 0 {
				i++
				rl++
			}
			if len(buffer) > 0 {
				buffer += "/"
			}
			buffer += fmt.Sprintf("D%d", first)
			if rl > 0 {
				buffer += fmt.Sprintf("-D%d", first+rl)
			}
		}
	}
	for i := 0; i < 8; i++ {
		if data&(1<<(i+8)) != 0 {
			first := i
			rl := 0
			for i < 7 && data&(1<<uint64(i+8+1)) != 0 {
				i++
				rl++
			}
			if len(buffer) > 0 {
				buffer += "/"
			}
			buffer += fmt.Sprintf("A%d", first)
			if rl > 0 {
				buffer += fmt.Sprintf("-A%d", first+rl)
			}
		}
	}
	return fmt.Sprintf("movem.l %s, %s", d.getEaModeStr(Long), buffer)
}

func d68000_movem_re_16(d *dasmInfo) string {
	data := d.readImm(Word)
	buffer := ""
	for i := 0; i < 8; i++ {
		if data&(1<<i) != 0 {
			first := i
			rl := 0
			for i < 7 && data&(1<<uint64(i+1)) != 0 {
				i++
				rl++
			}
			if len(buffer) > 0 {
				buffer += "/"
			}
			buffer += fmt.Sprintf("D%d", first)
			if rl > 0 {
				buffer += fmt.Sprintf("-D%d", first+rl)
			}
		}
	}
	for i := 0; i < 8; i++ {
		if data&(1<<(i+8)) != 0 {
			first := i
			rl := 0
			for i < 7 && data&(1<<uint64(i+8+1)) != 0 {
				i++
				rl++
			}
			if len(buffer) > 0 {
				buffer += "/"
			}
			buffer += fmt.Sprintf("A%d", first)
			if rl > 0 {
				buffer += fmt.Sprintf("-A%d", first+rl)
			}
		}
	}
	return fmt.Sprintf("movem.w %s, %s", buffer, d.getEaModeStr(Word))
}

// d68000_movem_re_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2411
func d68000_movem_re_32(d *dasmInfo) string {
	data := d.readImm(Word)
	buffer := ""
	for i := 0; i < 8; i++ {
		if data&(1<<(i)) != 0 {
			first := i
			rl := 0
			for i < 7 && data&(1<<(i+1)) != 0 {
				i++
				rl++
			}
			if len(buffer) > 0 {
				buffer += "/"
			}
			buffer += fmt.Sprintf("D%d", first)
			if rl > 0 {
				buffer += fmt.Sprintf("-D%d", first+rl)
			}
		}
	}
	for i := 0; i < 8; i++ {
		if data&(1<<(i+8)) != 0 {
			first := i
			rl := 0
			for i < 7 && data&(1<<uint64(i+8+1)) != 0 {
				i++
				rl++
			}
			if len(buffer) > 0 {
				buffer += "/"
			}
			buffer += fmt.Sprintf("A%d", first)
			if rl > 0 {
				buffer += fmt.Sprintf("-A%d", first+rl)
			}
		}
	}
	return fmt.Sprintf("movem.l %s, %s", buffer, d.getEaModeStr(Long))
}

func d68000_movep_re_16(d *dasmInfo) string {
	return fmt.Sprintf("movep.w D%d, ($%x,A%d)", x(d.ir), d.readImm(Word), y(d.ir))
}

func d68000_movep_re_32(d *dasmInfo) string {
	return fmt.Sprintf("movep.l D%d, ($%x,A%d)", x(d.ir), d.readImm(Word), y(d.ir))
}

func d68000_movep_er_16(d *dasmInfo) string {
	return fmt.Sprintf("movep.w ($%x,A%d), D%d", d.readImm(Word), y(d.ir), x(d.ir))
}

func d68000_movep_er_32(d *dasmInfo) string {
	return fmt.Sprintf("movep.l ($%x,A%d), D%d", d.readImm(Word), y(d.ir), x(d.ir))
}

func d68010_moves_8(d *dasmInfo) string {
	extension := d.readImm(Word)
	if extension&0x800 != 0 {
		return fmt.Sprintf("moves.b %c%d, %s; (1+)", func() int32 {
			if extension&0x8000 != 0 {
				return int32('A')
			}
			return int32('D')
		}(), extension>>uint64(12)&7, d.getEaModeStr(Byte))
	} else {
		return fmt.Sprintf("moves.b %s, %c%d; (1+)", d.getEaModeStr(Byte), func() int32 {
			if extension&0x8000 != 0 {
				return int32('A')
			}
			return int32('D')
		}(), (extension>>12)&7)
	}
}

func d68010_moves_16(d *dasmInfo) string {
	extension := d.readImm(Word)
	if extension&0x800 != 0 {
		return fmt.Sprintf("moves.w %c%d, %s; (1+)", func() int32 {
			if extension&0x8000 != 0 {
				return int32('A')
			}
			return int32('D')
		}(), extension>>uint64(12)&7, d.getEaModeStr(Word))
	} else {
		return fmt.Sprintf("moves.w %s, %c%d; (1+)", d.getEaModeStr(Word), func() int32 {
			if extension&0x8000 != 0 {
				return int32('A')
			}
			return int32('D')
		}(), (extension>>12)&7)
	}
}

func d68010_moves_32(d *dasmInfo) string {
	extension := d.readImm(Word)
	if extension&0x800 != 0 {
		return fmt.Sprintf("moves.l %c%d, %s; (1+)", func() int32 {
			if extension&0x8000 != 0 {
				return int32('A')
			}
			return int32('D')
		}(), extension>>uint64(12)&7, d.getEaModeStr(Long))
	} else {
		return fmt.Sprintf("moves.l %s, %c%d; (1+)", d.getEaModeStr(Long), func() int32 {
			if extension&0x8000 != 0 {
				return int32('A')
			}
			return int32('D')
		}(), extension>>uint64(12)&7)
	}
}

// d68000_moveq - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2512
func d68000_moveq(d *dasmInfo) string {
	return fmt.Sprintf("moveq   #%s, D%d", Byte.SignedHexString(int32(d.ir)), x(d.ir))
}

// d68040_move16_pi_pi - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2517
func d68040_move16_pi_pi(d *dasmInfo) string {
	return fmt.Sprintf("move16  (A%d)+, (A%d)+; (4)", y(d.ir), d.readImm(Word)>>uint64(12)&7)
}

// d68040_move16_pi_al - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2523
func d68040_move16_pi_al(d *dasmInfo) string {
	return fmt.Sprintf("move16  (A%d)+, %s; (4)", y(d.ir), d.getImmStrUnsigned(Long))
}

// d68040_move16_al_pi - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2529
func d68040_move16_al_pi(d *dasmInfo) string {
	return fmt.Sprintf("move16  %s, (A%d)+; (4)", d.getImmStrUnsigned(Long), y(d.ir))
}

// d68040_move16_ai_al - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2535
func d68040_move16_ai_al(d *dasmInfo) string {
	return fmt.Sprintf("move16  (A%d), %s; (4)", y(d.ir), d.getImmStrUnsigned(Long))
}

// d68040_move16_al_ai - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2541
func d68040_move16_al_ai(d *dasmInfo) string {
	return fmt.Sprintf("move16  %s, (A%d); (4)", d.getImmStrUnsigned(Long), y(d.ir))
}

// d68000_muls - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2547
func d68000_muls(d *dasmInfo) string {
	return fmt.Sprintf("muls.w  %s, D%d", d.getEaModeStr(Word), x(d.ir))
}

// d68000_mulu - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:0xFF2
func d68000_mulu(d *dasmInfo) string {
	return fmt.Sprintf("mulu.w  %s, D%d", d.getEaModeStr(Word), x(d.ir))
}

// d68020_mull - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:0xFF7
func d68020_mull(d *dasmInfo) string {
	extension := d.readImm(Word)
	if extension&1024 != 0 {
		return fmt.Sprintf("mul%c.l %s, D%d:D%d; (2+)", func() int32 {
			if extension&0x800 != 0 {
				return int32('s')
			}
			return int32('u')
		}(), d.getEaModeStr(Long), extension&7, extension>>uint64(12)&7)
	} else {
		return fmt.Sprintf("mul%c.l  %s, D%d; (2+)", func() int32 {
			if extension&0x800 != 0 {
				return int32('s')
			}
			return int32('u')
		}(), d.getEaModeStr(Long), extension>>uint64(12)&7)
	}
}

// d68000_nbcd - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:0x1009
func d68000_nbcd(d *dasmInfo) string {
	return fmt.Sprintf("nbcd    %s", d.getEaModeStr(Byte))
}

// d68000_neg_8 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2574
func d68000_neg_8(d *dasmInfo) string {
	return fmt.Sprintf("neg.b   %s", d.getEaModeStr(Byte))
}

// d68000_neg_16 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2579
func d68000_neg_16(d *dasmInfo) string {
	return fmt.Sprintf("neg.w   %s", d.getEaModeStr(Word))
}

// d68000_neg_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2584
func d68000_neg_32(d *dasmInfo) string {
	return fmt.Sprintf("neg.l   %s", d.getEaModeStr(Long))
}

// d68000_negx_8 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2589
func d68000_negx_8(d *dasmInfo) string {
	return fmt.Sprintf("negx.b  %s", d.getEaModeStr(Byte))
}

// d68000_negx_16 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2594
func d68000_negx_16(d *dasmInfo) string {
	return fmt.Sprintf("negx.w  %s", d.getEaModeStr(Word))
}

// d68000_negx_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2599
func d68000_negx_32(d *dasmInfo) string {
	return fmt.Sprintf("negx.l  %s", d.getEaModeStr(Long))
}

// d68000_nop - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2604
func d68000_nop(d *dasmInfo) string {
	return fmt.Sprintf("nop")
}

// d68000_not_8 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2609
func d68000_not_8(d *dasmInfo) string {
	return fmt.Sprintf("not.b   %s", d.getEaModeStr(Byte))
}

// d68000_not_16 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2614
func d68000_not_16(d *dasmInfo) string {
	return fmt.Sprintf("not.w   %s", d.getEaModeStr(Word))
}

// d68000_not_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2619
func d68000_not_32(d *dasmInfo) string {
	return fmt.Sprintf("not.l   %s", d.getEaModeStr(Long))
}

// d68000_or_er_8 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2624
func d68000_or_er_8(d *dasmInfo) string {
	return fmt.Sprintf("or.b    %s, D%d", d.getEaModeStr(Byte), x(d.ir))
}

// d68000_or_er_16 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2629
func d68000_or_er_16(d *dasmInfo) string {
	return fmt.Sprintf("or.w    %s, D%d", d.getEaModeStr(Word), x(d.ir))
}

// d68000_or_er_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2634
func d68000_or_er_32(d *dasmInfo) string {
	return fmt.Sprintf("or.l    %s, D%d", d.getEaModeStr(Long), x(d.ir))
}

// d68000_or_re_8 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2639
func d68000_or_re_8(d *dasmInfo) string {
	return fmt.Sprintf("or.b    D%d, %s", x(d.ir), d.getEaModeStr(Byte))
}

// d68000_or_re_16 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2644
func d68000_or_re_16(d *dasmInfo) string {
	return fmt.Sprintf("or.w    D%d, %s", x(d.ir), d.getEaModeStr(Word))
}

// d68000_or_re_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2649
func d68000_or_re_32(d *dasmInfo) string {
	return fmt.Sprintf("or.l    D%d, %s", x(d.ir), d.getEaModeStr(Long))
}

// d68000_ori_8 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2654
func d68000_ori_8(d *dasmInfo) string {
	var str []byte
	return fmt.Sprintf("ori.b   %s, %s", str, d.getEaModeStr(Byte))
}

// d68000_ori_16 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2660
func d68000_ori_16(d *dasmInfo) string {
	return fmt.Sprintf("ori.w   %s, %s", d.getImmStrUnsigned(Word), d.getEaModeStr(Word))
}

// d68000_ori_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2666
func d68000_ori_32(d *dasmInfo) string {
	return fmt.Sprintf("ori.l   %s, %s", d.getImmStrUnsigned(Long), d.getEaModeStr(Long))
}

// d68000_ori_to_ccr - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2672
func d68000_ori_to_ccr(d *dasmInfo) string {
	return fmt.Sprintf("ori     %s, CCR", d.getImmStrUnsigned(Byte))
}

// d68000_ori_to_sr - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2677
func d68000_ori_to_sr(d *dasmInfo) string {
	return fmt.Sprintf("ori     %s, SR", d.getImmStrUnsigned(Word))
}

func d68020_pack_rr(d *dasmInfo) string {
	return fmt.Sprintf("pack    D%d, D%d, %s; (2+)", y(d.ir), x(d.ir), d.getImmStrUnsigned(Word))
}

func d68020_pack_mm(d *dasmInfo) string {
	return fmt.Sprintf("pack    -(A%d), -(A%d), %s; (2+)", y(d.ir), x(d.ir), d.getImmStrUnsigned(Word))
}

func d68000_pea(d *dasmInfo) string {
	return fmt.Sprintf("pea     %s", d.getEaModeStr(Long))
}

func d68040_pflush(d *dasmInfo) string {
	ext := ""
	if d.ir&8 == 0 {
		ext = "n"
	}
	if d.ir&16 != 0 {
		return fmt.Sprintf("pflusha%s", ext)
	}
	return fmt.Sprintf("pflush%s(A%d)", ext, y(d.ir))
}

func d68000_reset(d *dasmInfo) string {
	return "reset"
}

func d68000_ror_s_8(d *dasmInfo) string {
	return fmt.Sprintf("ror.b   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_ror_s_16(d *dasmInfo) string {
	return fmt.Sprintf("ror.w   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_ror_s_32(d *dasmInfo) string {
	return fmt.Sprintf("ror.l   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_ror_r_8(d *dasmInfo) string {
	return fmt.Sprintf("ror.b   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_ror_r_16(d *dasmInfo) string {
	return fmt.Sprintf("ror.w   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_ror_r_32(d *dasmInfo) string {
	return fmt.Sprintf("ror.l   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_ror_ea(d *dasmInfo) string {
	return fmt.Sprintf("ror.w   %s", d.getEaModeStr(Long))
}

func d68000_rol_s_8(d *dasmInfo) string {
	return fmt.Sprintf("rol.b   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_rol_s_16(d *dasmInfo) string {
	return fmt.Sprintf("rol.w   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_rol_s_32(d *dasmInfo) string {
	return fmt.Sprintf("rol.l   #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_rol_r_8(d *dasmInfo) string {
	return fmt.Sprintf("rol.b   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_rol_r_16(d *dasmInfo) string {
	return fmt.Sprintf("rol.w   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_rol_r_32(d *dasmInfo) string {
	return fmt.Sprintf("rol.l   D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_rol_ea(d *dasmInfo) string {
	return fmt.Sprintf("rol.w   %s", d.getEaModeStr(Long))
}

func d68000_roxr_s_8(d *dasmInfo) string {
	return fmt.Sprintf("roxr.b  #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_roxr_s_16(d *dasmInfo) string {
	return fmt.Sprintf("roxr.w  #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_roxr_s_32(d *dasmInfo) string {
	return fmt.Sprintf("roxr.l  #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_roxr_r_8(d *dasmInfo) string {
	return fmt.Sprintf("roxr.b  D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_roxr_r_16(d *dasmInfo) string {
	return fmt.Sprintf("roxr.w  D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_roxr_r_32(d *dasmInfo) string {
	return fmt.Sprintf("roxr.l  D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_roxr_ea(d *dasmInfo) string {
	return fmt.Sprintf("roxr.w  %s", d.getEaModeStr(Long))
}

func d68000_roxl_s_8(d *dasmInfo) string {
	return fmt.Sprintf("roxl.b  #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_roxl_s_16(d *dasmInfo) string {
	return fmt.Sprintf("roxl.w  #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_roxl_s_32(d *dasmInfo) string {
	return fmt.Sprintf("roxl.l  #%d, D%d", g3bitQDataTable[x(d.ir)], y(d.ir))
}

func d68000_roxl_r_8(d *dasmInfo) string {
	return fmt.Sprintf("roxl.b  D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_roxl_r_16(d *dasmInfo) string {
	return fmt.Sprintf("roxl.w  D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_roxl_r_32(d *dasmInfo) string {
	return fmt.Sprintf("roxl.l  D%d, D%d", x(d.ir), y(d.ir))
}

func d68000_roxl_ea(d *dasmInfo) string {
	return fmt.Sprintf("roxl.w  %s", d.getEaModeStr(Long))
}

func d68010_rtd(d *dasmInfo) string {
	return fmt.Sprintf("rtd     %s; (1+)", d.getImmStrSigned(Word))
}

func d68000_rte(d *dasmInfo) string {
	return "rte"
}

func d68020_rtm(d *dasmInfo) string {
	return fmt.Sprintf("rtm     %c%d; (2+)",
		func() int32 {
			if d.ir&8 != 0 {
				return int32('A')
			}
			return int32('D')
		}(), y(d.ir))
}

func d68000_rtr(d *dasmInfo) string {
	return "rtr"
}

func d68000_rts(d *dasmInfo) string {
	return "rts"
}

func d68000_sbcd_rr(d *dasmInfo) string {
	return fmt.Sprintf("sbcd    D%d, D%d", y(d.ir), x(d.ir))
}

func d68000_sbcd_mm(d *dasmInfo) string {
	return fmt.Sprintf("sbcd    -(A%d), -(A%d)", y(d.ir), x(d.ir))
}

func d68000_scc(d *dasmInfo) string {
	return fmt.Sprintf("s%-2s     %s", gCCTable[d.ir>>uint64(8)&15], d.getEaModeStr(Byte))
}

func d68000_stop(d *dasmInfo) string {
	return fmt.Sprintf("stop    %s", d.getImmStrSigned(Word))
}

func d68000_sub_er_8(d *dasmInfo) string {
	return fmt.Sprintf("sub.b   %s, D%d", d.getEaModeStr(Byte), x(d.ir))
}

func d68000_sub_er_16(d *dasmInfo) string {
	return fmt.Sprintf("sub.w   %s, D%d", d.getEaModeStr(Word), x(d.ir))
}

func d68000_sub_er_32(d *dasmInfo) string {
	return fmt.Sprintf("sub.l   %s, D%d", d.getEaModeStr(Long), x(d.ir))
}

func d68000_sub_re_8(d *dasmInfo) string {
	return fmt.Sprintf("sub.b   D%d, %s", x(d.ir), d.getEaModeStr(Byte))
}

// d68000_sub_re_16 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2932
func d68000_sub_re_16(d *dasmInfo) string {
	return fmt.Sprintf("sub.w   D%d, %s", x(d.ir), d.getEaModeStr(Word))
}

// d68000_sub_re_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2937
func d68000_sub_re_32(d *dasmInfo) string {
	return fmt.Sprintf("sub.l   D%d, %s", x(d.ir), d.getEaModeStr(Long))
}

// d68000_suba_16 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2942
func d68000_suba_16(d *dasmInfo) string {
	return fmt.Sprintf("suba.w  %s, A%d", d.getEaModeStr(Word), x(d.ir))
}

// d68000_suba_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2947
func d68000_suba_32(d *dasmInfo) string {
	return fmt.Sprintf("suba.l  %s, A%d", d.getEaModeStr(Long), x(d.ir))
}

// d68000_subi_8 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2952
func d68000_subi_8(d *dasmInfo) string {
	var str []byte
	return fmt.Sprintf("subi.b  %s, %s", str, d.getEaModeStr(Byte))
}

// d68000_subi_16 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2958
func d68000_subi_16(d *dasmInfo) string {
	return fmt.Sprintf("subi.w  %s, %s", d.getImmStrSigned(Word), d.getEaModeStr(Word))
}

// d68000_subi_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2964
func d68000_subi_32(d *dasmInfo) string {
	return fmt.Sprintf("subi.l  %s, %s", d.getImmStrSigned(Long), d.getEaModeStr(Long))
}

// d68000_subq_8 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2970
func d68000_subq_8(d *dasmInfo) string {
	return fmt.Sprintf("subq.b  #%d, %s", g3bitQDataTable[x(d.ir)], d.getEaModeStr(Byte))
}

// d68000_subq_16 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2975
func d68000_subq_16(d *dasmInfo) string {
	return fmt.Sprintf("subq.w  #%d, %s", g3bitQDataTable[x(d.ir)], d.getEaModeStr(Word))
}

// d68000_subq_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2980
func d68000_subq_32(d *dasmInfo) string {
	return fmt.Sprintf("subq.l  #%d, %s", g3bitQDataTable[x(d.ir)], d.getEaModeStr(Long))
}

// d68000_subx_rr_8 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2985
func d68000_subx_rr_8(d *dasmInfo) string {
	return fmt.Sprintf("subx.b  D%d, D%d", y(d.ir), x(d.ir))
}

// d68000_subx_rr_16 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2990
func d68000_subx_rr_16(d *dasmInfo) string {
	return fmt.Sprintf("subx.w  D%d, D%d", y(d.ir), x(d.ir))
}

// d68000_subx_rr_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:2995
func d68000_subx_rr_32(d *dasmInfo) string {
	return fmt.Sprintf("subx.l  D%d, D%d", y(d.ir), x(d.ir))
}

// d68000_subx_mm_8 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:3000
func d68000_subx_mm_8(d *dasmInfo) string {
	return fmt.Sprintf("subx.b  -(A%d), -(A%d)", y(d.ir), x(d.ir))
}

// d68000_subx_mm_16 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:3005
func d68000_subx_mm_16(d *dasmInfo) string {
	return fmt.Sprintf("subx.w  -(A%d), -(A%d)", y(d.ir), x(d.ir))
}

func d68000_subx_mm_32(d *dasmInfo) string {
	return fmt.Sprintf("subx.l  -(A%d), -(A%d)", y(d.ir), x(d.ir))
}

func d68000_swap(d *dasmInfo) string {
	return fmt.Sprintf("swap    D%d", y(d.ir))
}

func d68000_tas(d *dasmInfo) string {
	return fmt.Sprintf("tas     %s", d.getEaModeStr(Byte))
}

func d68000_trap(d *dasmInfo) string {
	return fmt.Sprintf("trap    #$%x", d.ir&15)
}

func d68020_trapcc_0(d *dasmInfo) string {
	return fmt.Sprintf("trap%-2s; (2+)", gCCTable[d.ir>>uint64(8)&15])
}

func d68020_trapcc_16(d *dasmInfo) string {
	return fmt.Sprintf("trap%-2s  %s; (2+)", gCCTable[d.ir>>uint64(8)&15], d.getImmStrUnsigned(Word))
}

func d68020_trapcc_32(d *dasmInfo) string {
	return fmt.Sprintf("trap%-2s  %s; (2+)", gCCTable[d.ir>>uint64(8)&15], d.getImmStrUnsigned(Long))
}

func d68000_trapv(d *dasmInfo) string {
	return "trapv"
}

func d68000_tst_8(d *dasmInfo) string {
	return fmt.Sprintf("tst.b   %s", d.getEaModeStr(Byte))
}

func d68020_tst_pcdi_8(d *dasmInfo) string {
	return fmt.Sprintf("tst.b   %s; (2+)", d.getEaModeStr(Byte))
}

func d68020_tst_pcix_8(d *dasmInfo) string {
	return fmt.Sprintf("tst.b   %s; (2+)", d.getEaModeStr(Byte))
}

func d68020_tst_i_8(d *dasmInfo) string {
	return fmt.Sprintf("tst.b   %s; (2+)", d.getEaModeStr(Byte))
}

func d68000_tst_16(d *dasmInfo) string {
	return fmt.Sprintf("tst.w   %s", d.getEaModeStr(Word))
}

func d68020_tst_a_16(d *dasmInfo) string {
	return fmt.Sprintf("tst.w   %s; (2+)", d.getEaModeStr(Word))
}

func d68020_tst_pcdi_16(d *dasmInfo) string {
	return fmt.Sprintf("tst.w   %s; (2+)", d.getEaModeStr(Word))
}

func d68020_tst_pcix_16(d *dasmInfo) string {
	return fmt.Sprintf("tst.w   %s; (2+)", d.getEaModeStr(Word))
}

func d68020_tst_i_16(d *dasmInfo) string {
	return fmt.Sprintf("tst.w   %s; (2+)", d.getEaModeStr(Word))
}

func d68000_tst_32(d *dasmInfo) string {
	return fmt.Sprintf("tst.l   %s", d.getEaModeStr(Long))
}

func d68020_tst_a_32(d *dasmInfo) string {
	return fmt.Sprintf("tst.l   %s; (2+)", d.getEaModeStr(Long))
}

// d68020_tst_pcdi_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:3120
func d68020_tst_pcdi_32(d *dasmInfo) string {
	return fmt.Sprintf("tst.l   %s; (2+)", d.getEaModeStr(Long))
}

// d68020_tst_pcix_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:3126
func d68020_tst_pcix_32(d *dasmInfo) string {
	return fmt.Sprintf("tst.l   %s; (2+)", d.getEaModeStr(Long))
}

// d68020_tst_i_32 - transpiled function from  /home/jens/projects/Musashi/m68kdasm.c:3132
func d68020_tst_i_32(d *dasmInfo) string {
	return fmt.Sprintf("tst.l   %s; (2+)", d.getEaModeStr(Long))
}

func d68000_unlk(d *dasmInfo) string {
	return fmt.Sprintf("unlk    A%d", y(d.ir))
}

func d68020_unpk_rr(d *dasmInfo) string {
	return fmt.Sprintf("unpk    D%d, D%d, %s; (2+)", y(d.ir), x(d.ir), d.getImmStrUnsigned(Word))
}

func d68020_unpk_mm(d *dasmInfo) string {
	return fmt.Sprintf("unpk    -(A%d), -(A%d), %s; (2+)", y(d.ir), x(d.ir), d.getImmStrUnsigned(Word))
}

func d68851_p000(d *dasmInfo) string {
	modes := d.readImm(Word)
	// PFLUSH:  001xxx0xxxxxxxxx
	// PLOAD:   001000x0000xxxxx
	// PVALID1: 0010100000000000
	// PVALID2: 0010110000000xxx
	// PMOVE 1: 010xxxx000000000
	// PMOVE 2: 011xxxx0000xxx00
	// PMOVE 3: 011xxxx000000000
	// PTEST:   100xxxxxxxxxxxxx
	// PFLUSHR:  1010000000000000
	// do this after fetching the second PMOVE word so we properly get the 3rd if necessary
	str := d.getEaModeStr(Long)
	if modes&0xfde0 == 0x2000 {
		if modes&200 != 0 {
			// PLOAD
			return fmt.Sprintf("pload  #%d, %s", (modes>>10)&7, str)
		}
		return fmt.Sprintf("pload  %s, #%d", str, (modes>>10)&7)
	}
	if modes&0xe200 == 0x2000 {
		// PFLUSH
		return fmt.Sprintf("pflushr %x, %x, %s", modes&31, (modes>>5)&15, str)
	}
	if modes == 0xa000 {
		// PFLUSHR
		return fmt.Sprintf("pflushr %s", str)
	}
	if modes == 0x2800 {
		// PVALID (FORMAT 1)
		return fmt.Sprintf("pvalid VAL, %s", str)
	}
	if modes&0xfff8 == 0x2c00 {
		// PVALID (FORMAT 2)
		return fmt.Sprintf("pvalid A%d, %s", modes&15, str)
	}
	if modes&0xe000 == 0x8000 {
		// PTEST
		return fmt.Sprintf("ptest #%d, %s", modes&31, str)
	}
	switch (modes >> 13) & 7 {
	case 0, 2:
		if modes&0x100 != 0 {
			if modes&0x200 != 0 {
				// MC68030/040 form with FD bit
				// MC68881 form, FD never set
				return fmt.Sprintf("pmovefd  %s, %s", gMMURegs[(modes>>10)&7], str)
			}
			return fmt.Sprintf("pmovefd  %s, %s", str, gMMURegs[(modes>>10)&7])
		} else {
			if modes&0x200 != 0 {
				return fmt.Sprintf("pmove  %s, %s", gMMURegs[modes>>uint64(10)&7], str)
			}
			return fmt.Sprintf("pmove  %s, %s", str, gMMURegs[modes>>uint64(10)&7])
		}
	case 3:
		if modes&512 != 0 {
			// MC68030 to/from status reg
			return fmt.Sprintf("pmove  mmusr, %s", str)
		}
		return fmt.Sprintf("pmove  %s, mmusr", str)
	default:
		return fmt.Sprintf("pmove [unknown form] %s", str)
	}
}

func d68851_pbcc16(d *dasmInfo) string {
	return fmt.Sprintf("pb%s %x", gMMUCond[d.ir&15], d.pc+int32(int16(d.readImm(Word))))
}

func d68851_pbcc32(d *dasmInfo) string {
	return fmt.Sprintf("pb%s %x", gMMUCond[d.ir&15], d.pc+int32(d.readImm(Long)))
}

func d68851_pdbcc(d *dasmInfo) string {
	return fmt.Sprintf("pb%s %x", gMMUCond[uint16(d.readImm(Word))&0x0f], d.pc+d.readImm(Word))
}

func d68851_p001(d *dasmInfo) string {
	// PScc:  0000000000xxxxxx
	return "MMU 001 group"
}

func compare_nof_true_bits(a uint16, b uint16) bool {
	// Used by qsort
	a = a&0xaaaa>>1 + a&0x5555
	a = a&0xcccc>>2 + a&0x3333
	a = a&0xf0f0>>4 + a&0x0f0f
	a = a&0xff00>>8 + a&0x00ff

	b = b&0xaaaa>>1 + b&0x5555
	b = b&0xcccc>>2 + b&0x3333
	b = b&0xf0f0>>4 + b&0x0f0f
	b = b&0xff00>>8 + b&0x00ff
	// reversed to get greatest to least sorting
	return b-a <= 0
}

func init() {
	// build the opcode handler jump table
	sort.SliceStable(dasmOpcodeTable, func(i, j int) bool {
		return compare_nof_true_bits(dasmOpcodeTable[i].mask, dasmOpcodeTable[j].mask)
	})

	for i := 0; i < 0x10000; i++ {
		// default to illegal
		dasmTable[i] = d68000Illegal
		opcode := uint16(i)
		// search through opcode info for a match
		for _, ostruct := range dasmOpcodeTable {
			if opcode&ostruct.mask == ostruct.match {
				// match opcode mask and allowed ea modes
				// Handle destination ea for move instructions
				if validEA(opcode, ostruct.eaMask) {
					dasmTable[i] = ostruct.dasmHandler
				}
			}
		}
	}
}

func Disassemble(cpuType Type, pc int32, bus AddressBus) (string, int32) {
	d := &dasmInfo{pc: pc, helper: "", bus: bus}
	d.ir = uint16(d.readImm(Word))

	if d.cpuType&dasmSupportedTypes[d.ir] == 0 {
		if d.ir&0xf000 == 0xf000 {
			return d68000LineF(d), d.pc - pc
		}
		return d68000Illegal(d), d.pc - pc
	}

	dasmTable[d.ir](d)
	return d.helper, d.pc - pc
}

var dasmOpcodeTable = []dasmOpcode{
	{d68000LineA, 0xf000, 0xa000, 0, dasmAll},
	{d68000LineF, 0xf000, 0xf000, 0, dasmAll},
	{d68000_abcd_rr, 0xf1f8, 49408, 0, dasmAll},
	{d68000_abcd_mm, 0xf1f8, 49416, 0, dasmAll},
	{d68000_add_er_8, 0xf1c0, 0xD000, 3071, dasmAll},
	{d68000_add_er_16, 0xf1c0, 53312, 4095, dasmAll},
	{d68000_add_er_32, 0xf1c0, 53376, 4095, dasmAll},
	{d68000_add_re_8, 0xf1c0, 53504, 1016, dasmAll},
	{d68000_add_re_16, 0xf1c0, 53568, 1016, dasmAll},
	{d68000_add_re_32, 0xf1c0, 53632, 1016, dasmAll},
	{d68000_adda_16, 0xf1c0, 53440, 4095, dasmAll},
	{d68000_adda_32, 0xf1c0, 53696, 4095, dasmAll},
	{d68000_addi_8, 0xffc0, 1536, 3064, dasmAll},
	{d68000_addi_16, 0xffc0, 1600, 3064, dasmAll},
	{d68000_addi_32, 0xffc0, 1664, 3064, dasmAll},
	{d68000_addq_8, 0xf1c0, 0x8000, 3064, dasmAll},
	{d68000_addq_16, 0xf1c0, 20544, 4088, dasmAll},
	{d68000_addq_32, 0xf1c0, 20608, 4088, dasmAll},
	{d68000_addx_rr_8, 0xf1f8, 53504, 0, dasmAll},
	{d68000_addx_rr_16, 0xf1f8, 53568, 0, dasmAll},
	{d68000_addx_rr_32, 0xf1f8, 53632, 0, dasmAll},
	{d68000_addx_mm_8, 0xf1f8, 53512, 0, dasmAll},
	{d68000_addx_mm_16, 0xf1f8, 53576, 0, dasmAll},
	{d68000_addx_mm_32, 0xf1f8, 53640, 0, dasmAll},
	{d68000_and_er_8, 0xf1c0, 49152, 3071, dasmAll},
	{d68000_and_er_16, 0xf1c0, 49216, 3071, dasmAll},
	{d68000_and_er_32, 0xf1c0, 49280, 3071, dasmAll},
	{d68000_and_re_8, 0xf1c0, 49408, 1016, dasmAll},
	{d68000_and_re_16, 0xf1c0, 49472, 1016, dasmAll},
	{d68000_and_re_32, 0xf1c0, 49536, 1016, dasmAll},
	{d68000_andi_to_ccr, 0xffff, 572, 0, dasmAll},
	{d68000_andi_to_sr, 0xffff, 636, 0, dasmAll},
	{d68000_andi_8, 0xffc0, 512, 3064, dasmAll},
	{d68000_andi_16, 0xffc0, 576, 3064, dasmAll},
	{d68000_andi_32, 0xffc0, 640, 3064, dasmAll},
	{d68000_asr_s_8, 0xf1f8, 57344, 0, dasmAll},
	{d68000_asr_s_16, 0xf1f8, 57408, 0, dasmAll},
	{d68000_asr_s_32, 0xf1f8, 57472, 0, dasmAll},
	{d68000_asr_r_8, 0xf1f8, 57376, 0, dasmAll},
	{d68000_asr_r_16, 0xf1f8, 57440, 0, dasmAll},
	{d68000_asr_r_32, 0xf1f8, 57504, 0, dasmAll},
	{d68000_asr_ea, 0xffc0, 57536, 1016, dasmAll},
	{d68000_asl_s_8, 0xf1f8, 57600, 0, dasmAll},
	{d68000_asl_s_16, 0xf1f8, 57664, 0, dasmAll},
	{d68000_asl_s_32, 0xf1f8, 57728, 0, dasmAll},
	{d68000_asl_r_8, 0xf1f8, 57632, 0, dasmAll},
	{d68000_asl_r_16, 0xf1f8, 57696, 0, dasmAll},
	{d68000_asl_r_32, 0xf1f8, 57760, 0, dasmAll},
	{d68000_asl_ea, 0xffc0, 57792, 1016, dasmAll},
	{d68000_bcc_8, 0xf000, 0x6000, 0, dasmAll},
	{d68000_bcc_16, 0xf0ff, 0x6000, 0, dasmAll},
	{d68020_bcc_32, 0xf0ff, 24831, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_bchg_r, 0xf1c0, 320, 3064, dasmAll},
	{d68000_bchg_s, 0xffc0, 2112, 3064, dasmAll},
	{d68000_bclr_r, 0xf1c0, 384, 3064, dasmAll},
	{d68000_bclr_s, 0xffc0, 2176, 3064, dasmAll},
	{d68020_bfchg, 0xffc0, 60096, 2680, dasm68020 | dasm68030 | dasm68040},
	{d68020_bfclr, 0xffc0, 60608, 2680, dasm68020 | dasm68030 | dasm68040},
	{d68020_bfexts, 0xffc0, 60352, 2683, dasm68020 | dasm68030 | dasm68040},
	{d68020_bfextu, 0xffc0, 59840, 2683, dasm68020 | dasm68030 | dasm68040},
	{d68020_bfffo, 0xffc0, 60864, 2683, dasm68020 | dasm68030 | dasm68040},
	{d68020_bfins, 0xffc0, 61376, 2680, dasm68020 | dasm68030 | dasm68040},
	{d68020_bfset, 0xffc0, 61120, 2680, dasm68020 | dasm68030 | dasm68040},
	{d68020_bftst, 0xffc0, 59584, 2683, dasm68020 | dasm68030 | dasm68040},
	{d68010_bkpt, 0xfff8, 18504, 0, dasm68010 | dasm68020 | dasm68030 | dasm68040},
	{d68000_bra_8, 0xff00, 0x6000, 0, dasmAll},
	{d68000_bra_16, 0xffff, 0x6000, 0, dasmAll},
	{d68020_bra_32, 0xffff, 24831, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_bset_r, 0xf1c0, 4048, 3064, dasmAll},
	{d68000_bset_s, 0xffc0, 2240, 3064, dasmAll},
	{d68000_bsr_8, 0xff00, 24832, 0, dasmAll},
	{d68000_bsr_16, 0xffff, 24832, 0, dasmAll},
	{d68020_bsr_32, 0xffff, 25087, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_btst_r, 0xf1c0, 0x100, 3071, dasmAll},
	{d68000_btst_s, 0xffc0, 0x800, 3067, dasmAll},
	{d68020_callm, 0xffc0, 1728, 635, dasm68020 | dasm68030 | dasm68040},
	{d68020_cas_8, 0xffc0, 2752, 1016, dasm68020 | dasm68030 | dasm68040},
	{d68020_cas_16, 0xffc0, 3264, 1016, dasm68020 | dasm68030 | dasm68040},
	{d68020_cas_32, 0xffc0, 3776, 1016, dasm68020 | dasm68030 | dasm68040},
	{d68020_cas2_16, 0xffff, 3324, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_cas2_32, 0xffff, 3836, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_chk_16, 0xf1c0, 16768, 3071, dasmAll},
	{d68020_chk_32, 0xf1c0, 16640, 3071, dasm68020 | dasm68030 | dasm68040},
	{d68020_chk2_cmp2_8, 0xffc0, 0xC0, 635, dasm68020 | dasm68030 | dasm68040},
	{d68020_chk2_cmp2_16, 0xffc0, 704, 635, dasm68020 | dasm68030 | dasm68040},
	{d68020_chk2_cmp2_32, 0xffc0, 1216, 635, dasm68020 | dasm68030 | dasm68040},
	{d68040_cinv, 0xff20, 62464, 0, dasm68040},
	{d68000_clr_8, 0xffc0, 16896, 3064, dasmAll},
	{d68000_clr_16, 0xffc0, 16960, 3064, dasmAll},
	{d68000_clr_32, 0xffc0, 17024, 3064, dasmAll},
	{d68000_cmp_8, 0xf1c0, 45056, 3071, dasmAll},
	{d68000_cmp_16, 0xf1c0, 45120, 4095, dasmAll},
	{d68000_cmp_32, 0xf1c0, 45184, 4095, dasmAll},
	{d68000_cmpa_16, 0xf1c0, 45248, 4095, dasmAll},
	{d68000_cmpa_32, 0xf1c0, 45504, 4095, dasmAll},
	{d68000_cmpi_8, 0xffc0, 3072, 3064, dasmAll},
	{d68020_cmpi_pcdi_8, 0xffff, 3130, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_cmpi_pcix_8, 0xffff, 3131, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_cmpi_16, 0xffc0, 3136, 3064, dasmAll},
	{d68020_cmpi_pcdi_16, 0xffff, 3194, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_cmpi_pcix_16, 0xffff, 3195, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_cmpi_32, 0xffc0, 3200, 3064, dasmAll},
	{d68020_cmpi_pcdi_32, 0xffff, 3258, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_cmpi_pcix_32, 0xffff, 3259, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_cmpm_8, 0xf1f8, 45320, 0, dasmAll},
	{d68000_cmpm_16, 0xf1f8, 45384, 0, dasmAll},
	{d68000_cmpm_32, 0xf1f8, 45448, 0, dasmAll},
	{d68020_cpbcc_16, 0xf1c0, 61568, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_cpbcc_32, 0xf1c0, 0xf0c0, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_cpdbcc, 0xf1f8, 61512, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_cpgen, 0xf1c0, 0xf000, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_cprestore, 0xf1c0, 61760, 895, dasm68020 | dasm68030 | dasm68040},
	{d68020_cpsave, 0xf1c0, 61696, 760, dasm68020 | dasm68030 | dasm68040},
	{d68020_cpscc, 0xf1c0, 61504, 3064, dasm68020 | dasm68030 | dasm68040},
	{d68020_cptrapcc_0, 61951, 61564, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_cptrapcc_16, 61951, 61562, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_cptrapcc_32, 61951, 61563, 0, dasm68020 | dasm68030 | dasm68040},
	{d68040_cpush, 0xff20, 62496, 0, dasm68040},
	{d68000_dbcc, 0xf0f8, 20680, 0, dasmAll},
	{d68000_dbra, 0xfff8, 20936, 0, dasmAll},
	{d68000_divs, 0xf1c0, 33216, 3071, dasmAll},
	{d68000_divu, 0xf1c0, 32960, 3071, dasmAll},
	{d68020_divl, 0xffc0, 19520, 3071, dasm68020 | dasm68030 | dasm68040},
	{d68000_eor_8, 0xf1c0, 45312, 3064, dasmAll},
	{d68000_eor_16, 0xf1c0, 45376, 3064, dasmAll},
	{d68000_eor_32, 0xf1c0, 45440, 3064, dasmAll},
	{d68000_eori_to_ccr, 0xffff, 2620, 0, dasmAll},
	{d68000_eori_to_sr, 0xffff, 2684, 0, dasmAll},
	{d68000_eori_8, 0xffc0, 0x1000, 3064, dasmAll},
	{d68000_eori_16, 0xffc0, 2624, 3064, dasmAll},
	{d68000_eori_32, 0xffc0, 2688, 3064, dasmAll},
	{d68000_exg_dd, 0xf1f8, 49472, 0, dasmAll},
	{d68000_exg_aa, 0xf1f8, 49480, 0, dasmAll},
	{d68000_exg_da, 0xf1f8, 49544, 0, dasmAll},
	{d68020_extb_32, 0xfff8, 18880, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_ext_16, 0xfff8, 18560, 0, dasmAll},
	{d68000_ext_32, 0xfff8, 18624, 0, dasmAll},
	{d68040_fpu, 0xffc0, 61952, 0, dasm68040},
	{d68000Illegal, 0xffff, 19196, 0, dasmAll},
	{d68000_jmp, 0xffc0, 20160, 635, dasmAll},
	{d68000_jsr, 0xffc0, 20096, 635, dasmAll},
	{d68000_lea, 0xf1c0, 16832, 635, dasmAll},
	{d68000_link_16, 0xfff8, 20048, 0, dasmAll},
	{d68020_link_32, 0xfff8, 18440, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_lsr_s_8, 0xf1f8, 57352, 0, dasmAll},
	{d68000_lsr_s_16, 0xf1f8, 57416, 0, dasmAll},
	{d68000_lsr_s_32, 0xf1f8, 57480, 0, dasmAll},
	{d68000_lsr_r_8, 0xf1f8, 57384, 0, dasmAll},
	{d68000_lsr_r_16, 0xf1f8, 57448, 0, dasmAll},
	{d68000_lsr_r_32, 0xf1f8, 57512, 0, dasmAll},
	{d68000_lsr_ea, 0xffc0, 58048, 1016, dasmAll},
	{d68000_lsl_s_8, 0xf1f8, 57608, 0, dasmAll},
	{d68000_lsl_s_16, 0xf1f8, 57672, 0, dasmAll},
	{d68000_lsl_s_32, 0xf1f8, 57736, 0, dasmAll},
	{d68000_lsl_r_8, 0xf1f8, 57640, 0, dasmAll},
	{d68000_lsl_r_16, 0xf1f8, 57704, 0, dasmAll},
	{d68000_lsl_r_32, 0xf1f8, 57768, 0, dasmAll},
	{d68000_lsl_ea, 0xffc0, 58304, 1016, dasmAll},
	{d68000_move_8, 0xf000, 4096, 3071, dasmAll},
	{d68000_move_16, 0xf000, 10228, 4095, dasmAll},
	{d68000_move_32, 0xf000, 8192, 4095, dasmAll},
	{d68000_movea_16, 0xf1c0, 12352, 4095, dasmAll},
	{d68000_movea_32, 0xf1c0, 8255, 4095, dasmAll},
	{d68000_move_to_ccr, 0xffc0, 17600, 3071, dasmAll},
	{d68010_move_fr_ccr, 0xffc0, 17088, 3064, dasmAll},
	{d68000_move_to_sr, 0xffc0, 18112, 3071, dasmAll},
	{d68000_move_fr_sr, 0xffc0, 16576, 3064, dasmAll},
	{d68000_move_to_usp, 0xfff8, 20064, 0, dasmAll},
	{d68000_move_fr_usp, 0xfff8, 20072, 0, dasmAll},
	{d68010_movec, 65534, 20090, 0, dasm68010 | dasm68020 | dasm68030 | dasm68040},
	{d68000_movem_pd_16, 0xfff8, 18592, 0, dasmAll},
	{d68000_movem_pd_32, 0xfff8, 18656, 0, dasmAll},
	{d68000_movem_re_16, 0xffc0, 18560, 760, dasmAll},
	{d68000_movem_re_32, 0xffc0, 18624, 760, dasmAll},
	{d68000_movem_er_16, 0xffc0, 19584, 891, dasmAll},
	{d68000_movem_er_32, 0xffc0, 19648, 891, dasmAll},
	{d68000_movep_er_16, 0xf1f8, 264, 0, dasmAll},
	{d68000_movep_er_32, 0xf1f8, 328, 0, dasmAll},
	{d68000_movep_re_16, 0xf1f8, 392, 0, dasmAll},
	{d68000_movep_re_32, 0xf1f8, 456, 0, dasmAll},
	{d68010_moves_8, 0xffc0, 3584, 1016, dasm68010 | dasm68020 | dasm68030 | dasm68040},
	{d68010_moves_16, 0xffc0, 3648, 1016, dasm68010 | dasm68020 | dasm68030 | dasm68040},
	{d68010_moves_32, 0xffc0, 3712, 1016, dasm68010 | dasm68020 | dasm68030 | dasm68040},
	{d68000_moveq, 61696, 28672, 0, dasmAll},
	{d68040_move16_pi_pi, 0xfff8, 63008, 0, dasm68040},
	{d68040_move16_pi_al, 0xfff8, 62976, 0, dasm68040},
	{d68040_move16_al_pi, 0xfff8, 62984, 0, dasm68040},
	{d68040_move16_ai_al, 0xfff8, 62992, 0, dasm68040},
	{d68040_move16_al_ai, 0xfff8, 63000, 0, dasm68040},
	{d68000_muls, 0xf1c0, 49600, 3071, dasmAll},
	{d68000_mulu, 0xf1c0, 49344, 3071, dasmAll},
	{d68020_mull, 0xffc0, 19456, 3071, dasm68020},
	{d68000_nbcd, 0xffc0, 18432, 3064, dasmAll},
	{d68000_neg_8, 0xffc0, 17408, 3064, dasmAll},
	{d68000_neg_16, 0xffc0, 17472, 3064, dasmAll},
	{d68000_neg_32, 0xffc0, 17536, 3064, dasmAll},
	{d68000_negx_8, 0xffc0, 16384, 3064, dasmAll},
	{d68000_negx_16, 0xffc0, 16448, 3064, dasmAll},
	{d68000_negx_32, 0xffc0, 16512, 3064, dasmAll},
	{d68000_nop, 0xffff, 20081, 0, dasmAll},
	{d68000_not_8, 0xffc0, 17920, 3064, dasmAll},
	{d68000_not_16, 0xffc0, 17984, 3064, dasmAll},
	{d68000_not_32, 0xffc0, 18048, 3064, dasmAll},
	{d68000_or_er_8, 0xf1c0, 0x8000, 3071, dasmAll},
	{d68000_or_er_16, 0xf1c0, 32832, 3071, dasmAll},
	{d68000_or_er_32, 0xf1c0, 32896, 3071, dasmAll},
	{d68000_or_re_8, 0xf1c0, 33024, 1016, dasmAll},
	{d68000_or_re_16, 0xf1c0, 33088, 1016, dasmAll},
	{d68000_or_re_32, 0xf1c0, 33152, 1016, dasmAll},
	{d68000_ori_to_ccr, 0xffff, 60, 0, dasmAll},
	{d68000_ori_to_sr, 0xffff, 124, 0, dasmAll},
	{d68000_ori_8, 0xffc0, 0, 3064, dasmAll},
	{d68000_ori_16, 0xffc0, 64, 3064, dasmAll},
	{d68000_ori_32, 0xffc0, 128, 3064, dasmAll},
	{d68020_pack_rr, 0xf1f8, 33088, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_pack_mm, 0xf1f8, 33096, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_pea, 0xffc0, 18496, 635, dasmAll},
	{d68040_pflush, 0xffe0, 62720, 0, dasm68040},
	{d68000_reset, 0xffff, 20080, 0, dasmAll},
	{d68000_ror_s_8, 0xf1f8, 57368, 0, dasmAll},
	{d68000_ror_s_16, 0xf1f8, 57432, 0, dasmAll},
	{d68000_ror_s_32, 0xf1f8, 57496, 0, dasmAll},
	{d68000_ror_r_8, 0xf1f8, 57400, 0, dasmAll},
	{d68000_ror_r_16, 0xf1f8, 57464, 0, dasmAll},
	{d68000_ror_r_32, 0xf1f8, 57528, 0, dasmAll},
	{d68000_ror_ea, 0xffc0, 59072, 1016, dasmAll},
	{d68000_rol_s_8, 0xf1f8, 57624, 0, dasmAll},
	{d68000_rol_s_16, 0xf1f8, 57688, 0, dasmAll},
	{d68000_rol_s_32, 0xf1f8, 57752, 0, dasmAll},
	{d68000_rol_r_8, 0xf1f8, 57656, 0, dasmAll},
	{d68000_rol_r_16, 0xf1f8, 57720, 0, dasmAll},
	{d68000_rol_r_32, 0xf1f8, 57784, 0, dasmAll},
	{d68000_rol_ea, 0xffc0, 59328, 1016, dasmAll},
	{d68000_roxr_s_8, 0xf1f8, 57360, 0, dasmAll},
	{d68000_roxr_s_16, 0xf1f8, 57424, 0, dasmAll},
	{d68000_roxr_s_32, 0xf1f8, 57488, 0, dasmAll},
	{d68000_roxr_r_8, 0xf1f8, 57392, 0, dasmAll},
	{d68000_roxr_r_16, 0xf1f8, 57456, 0, dasmAll},
	{d68000_roxr_r_32, 0xf1f8, 57520, 0, dasmAll},
	{d68000_roxr_ea, 0xffc0, 58560, 1016, dasmAll},
	{d68000_roxl_s_8, 0xf1f8, 57616, 0, dasmAll},
	{d68000_roxl_s_16, 0xf1f8, 57680, 0, dasmAll},
	{d68000_roxl_s_32, 0xf1f8, 57744, 0, dasmAll},
	{d68000_roxl_r_8, 0xf1f8, 57648, 0, dasmAll},
	{d68000_roxl_r_16, 0xf1f8, 57712, 0, dasmAll},
	{d68000_roxl_r_32, 0xf1f8, 57776, 0, dasmAll},
	{d68000_roxl_ea, 0xffc0, 58816, 1016, dasmAll},
	{d68010_rtd, 0xffff, 20084, 0, dasm68010 | dasm68020 | dasm68030 | dasm68040},
	{d68010_rtd, 0xffff, 20084, 0, dasm68010 | dasm68020 | dasm68030 | dasm68040},
	{d68000_rte, 0xffff, 20083, 0, dasmAll},
	{d68020_rtm, 0xfff0, 1728, 0, dasm68020},
	{d68000_rtr, 0xffff, 20087, 0, dasmAll},
	{d68000_rts, 0xffff, 20085, 0, dasmAll},
	{d68000_sbcd_rr, 0xf1f8, 33024, 0, dasmAll},
	{d68000_sbcd_mm, 0xf1f8, 33032, 0, dasmAll},
	{d68000_scc, 0xf0c0, 20672, 3064, dasmAll},
	{d68000_stop, 0xffff, 20082, 0, dasmAll},
	{d68000_sub_er_8, 0xf1c0, 36864, 3071, dasmAll},
	{d68000_sub_er_16, 0xf1c0, 36928, 4095, dasmAll},
	{d68000_sub_er_32, 0xf1c0, 36992, 4095, dasmAll},
	{d68000_sub_re_8, 0xf1c0, 37120, 1016, dasmAll},
	{d68000_sub_re_16, 0xf1c0, 37184, 1016, dasmAll},
	{d68000_sub_re_32, 0xf1c0, 3748, 1016, dasmAll},
	{d68000_suba_16, 0xf1c0, 37056, 4095, dasmAll},
	{d68000_suba_32, 0xf1c0, 37312, 4095, dasmAll},
	{d68000_subi_8, 0xffc0, 1024, 3064, dasmAll},
	{d68000_subi_16, 0xffc0, 1088, 3064, dasmAll},
	{d68000_subi_32, 0xffc0, 1152, 3064, dasmAll},
	{d68000_subq_8, 0xf1c0, 20736, 3064, dasmAll},
	{d68000_subq_16, 0xf1c0, 20800, 4088, dasmAll},
	{d68000_subq_32, 0xf1c0, 20864, 4088, dasmAll},
	{d68000_subx_rr_8, 0xf1f8, 37120, 0, dasmAll},
	{d68000_subx_rr_16, 0xf1f8, 37184, 0, dasmAll},
	{d68000_subx_rr_32, 0xf1f8, 37248, 0, dasmAll},
	{d68000_subx_mm_8, 0xf1f8, 37128, 0, dasmAll},
	{d68000_subx_mm_16, 0xf1f8, 37192, 0, dasmAll},
	{d68000_subx_mm_32, 0xf1f8, 37255, 0, dasmAll},
	{d68000_swap, 0xfff8, 18496, 0, dasmAll},
	{d68000_tas, 0xffc0, 19136, 3064, dasmAll},
	{d68000_trap, 0xfff0, 20032, 0, dasmAll},
	{d68020_trapcc_0, 0xf0ff, 20732, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_trapcc_16, 0xf0ff, 20730, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_trapcc_32, 0xf0ff, 20731, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_trapv, 0xffff, 20086, 0, dasmAll},
	{d68000_tst_8, 0xffc0, 18944, 3064, dasmAll},
	{d68020_tst_pcdi_8, 0xffff, 19002, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_tst_pcix_8, 0xffff, 19003, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_tst_i_8, 0xffff, 19004, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_tst_16, 0xffc0, 19008, 3064, dasmAll},
	{d68020_tst_a_16, 0xfff8, 19016, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_tst_pcdi_16, 0xffff, 19066, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_tst_pcix_16, 0xffff, 19067, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_tst_i_16, 0xffff, 19068, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_tst_32, 0xffc0, 19072, 3064, dasmAll},
	{d68020_tst_a_32, 0xfff8, 19080, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_tst_pcdi_32, 0xffff, 19130, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_tst_pcix_32, 0xffff, 19131, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_tst_i_32, 0xffff, 19132, 0, dasm68020 | dasm68030 | dasm68040},
	{d68000_unlk, 0xfff8, 20056, 0, dasmAll},
	{d68020_unpk_rr, 0xf1f8, 33152, 0, dasm68020 | dasm68030 | dasm68040},
	{d68020_unpk_mm, 0xf1f8, 33160, 0, dasm68020 | dasm68030 | dasm68040},
	{d68851_p000, 0xffc0, 0xf000, 0, dasmAll},
	{d68851_pbcc16, 0xffc0, 61568, 0, dasmAll},
	{d68851_pbcc32, 0xffc0, 0xf0c0, 0, dasmAll},
	{d68851_pdbcc, 0xfff8, 61512, 0, dasmAll},
	{d68851_p001, 0xffc0, 61504, 0, dasmAll}}
