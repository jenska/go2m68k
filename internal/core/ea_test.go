package core

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tinyCore = func() *Core {
	var mem [1024 * 16]byte
	return Build(
		func(address uint32, o Operand) uint32 {
			return o.Read(mem[address:])
		},
		func(address uint32, o Operand, v uint32) {
			o.Write(v, mem[address:])
		},
		func() {

		})(M68000InstructionSet)
}

func TestCore_resolveEA(t *testing.T) {
	c := tinyCore()
	type args struct {
		mode uint16
		xn   uint16
		o    Operand
	}
	tests := []struct {
		name string
		c    *Core
		args args
		want EA
	}{
		{"data register d0.b", c, args{ModeDN, 0, Byte}, eaRegister{&c.D[0], Byte}},
		{"data register d0.w", c, args{ModeDN, 0, Word}, eaRegister{&c.D[0], Word}},
		{"data register d0.l", c, args{ModeDN, 0, Long}, eaRegister{&c.D[0], Long}},
		{"data register d1.b", c, args{ModeDN, 1, Byte}, eaRegister{&c.D[1], Byte}},
		{"data register d1.w", c, args{ModeDN, 1, Word}, eaRegister{&c.D[1], Word}},
		{"data register d1.l", c, args{ModeDN, 1, Long}, eaRegister{&c.D[1], Long}},
		{"address register a0.b", c, args{ModeAN, 0, Byte}, eaRegister{&c.A[0], Byte}},
		{"address register a0.w", c, args{ModeAN, 0, Word}, eaRegister{&c.A[0], Word}},
		{"address register a0.l", c, args{ModeAN, 0, Long}, eaRegister{&c.A[0], Long}},
		{"address register a1.b", c, args{ModeAN, 1, Byte}, eaRegister{&c.A[1], Byte}},
		{"address register a1.w", c, args{ModeAN, 1, Word}, eaRegister{&c.A[1], Word}},
		{"address register a1.l", c, args{ModeAN, 1, Long}, eaRegister{&c.A[1], Long}},
		{"address register indirect (a0).b", c, args{ModeAI, 0, Byte}, eaAddress{1, c, Byte}},
		{"address register indirect (a0).w", c, args{ModeAI, 0, Word}, eaAddress{1, c, Word}},
		{"address register indirect (a0).l", c, args{ModeAI, 0, Long}, eaAddress{1, c, Long}},
		{"address register indirect (a1).b", c, args{ModeAI, 1, Byte}, eaAddress{1, c, Byte}},
		{"address register indirect (a1).w", c, args{ModeAI, 1, Word}, eaAddress{1, c, Word}},
		{"address register indirect (a1).l", c, args{ModeAI, 1, Long}, eaAddress{1, c, Long}},
		{"post increment .b (a0)+", c, args{ModePI, 0, Byte}, eaAddress{1, c, Byte}},
		{"post increment .w (a0)+", c, args{ModePI, 0, Word}, eaAddress{2, c, Word}},
		{"post increment .l (a0)+", c, args{ModePI, 0, Long}, eaAddress{4, c, Long}},
		{"post increment .b (a1)+", c, args{ModePI, 1, Byte}, eaAddress{1, c, Byte}},
		{"post increment .w (a1)+", c, args{ModePI, 1, Word}, eaAddress{2, c, Word}},
		{"post increment .l (a1)+", c, args{ModePI, 1, Long}, eaAddress{4, c, Long}},
		{"pre decrement .b -(a0)", c, args{ModePD, 0, Byte}, eaAddress{7, c, Byte}},
		{"pre decrement .w -(a0)", c, args{ModePD, 0, Word}, eaAddress{5, c, Word}},
		{"pre decrement .l -(a0)", c, args{ModePD, 0, Long}, eaAddress{1, c, Long}},
		{"pre decrement .b -(a1)", c, args{ModePD, 1, Byte}, eaAddress{7, c, Byte}},
		{"pre decrement .w -(a1)", c, args{ModePD, 1, Word}, eaAddress{5, c, Word}},
		{"pre decrement .l -(a1)", c, args{ModePD, 1, Long}, eaAddress{1, c, Long}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.FetchEA(tt.args.mode, tt.args.xn, tt.args.o)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Core.resolveEA() = %v, want %v", got, tt.want)
			}
			got.Write(1)
			assert.Equal(t, uint32(1), got.Read())
		})
	}
}
