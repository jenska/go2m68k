// Code generated by "stringer -type Register"; DO NOT EDIT.

package m68k

import "fmt"

const _Register_name = "RegD0RegD1RegD2RegD3RegD4RegD5RegD6RegD7RegA0RegA1RegA2RegA3RegA4RegA5RegA6RegA7RegPCRegSRRegSPRegUSPRegISPRegMSPRegSFCRegDFCRegVBRRegCACRRegCAARRegPrefAddrRegPrefDataRegPPCRegIRRegCPUType"

var _Register_index = [...]uint8{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 101, 107, 113, 119, 125, 131, 138, 145, 156, 167, 173, 178, 188}

func (i Register) String() string {
	if i >= Register(len(_Register_index)-1) {
		return fmt.Sprintf("Register(%d)", i)
	}
	return _Register_name[_Register_index[i]:_Register_index[i+1]]
}
