package cpu

import (
	"reflect"
	"testing"
)

func Test_int32ToFloatx80(t *testing.T) {
	type args struct {
		a int32
	}
	tests := []struct {
		name string
		args args
		want floatx80
	}{
		{"zero", args{0}, floatx80{0, 0}},
		{"one", args{1}, floatx80{0x3fff, 0x8000000000000000}},
		{"two", args{2}, floatx80{0x4000, 0x8000000000000000}},
		{"three", args{3}, floatx80{0x4000, 0xc000000000000000}},
		{"-one", args{-1}, floatx80{0xbfff, 0x8000000000000000}},
		{"-two", args{-2}, floatx80{0xc000, 0x8000000000000000}},
		{"-three", args{-3}, floatx80{0xc000, 0xc000000000000000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := int32ToFloatx80(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("int32ToFloatx80() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_float32ToFloatx80(t *testing.T) {
	type args struct {
		a float32
	}
	tests := []struct {
		name string
		args args
		want floatx80
	}{
		{"zero", args{0.0}, floatx80{0, 0}},
		{"one", args{1.0}, floatx80{0x3fff, 0x8000000000000000}},
		{"two", args{2.0}, floatx80{0x4000, 0x8000000000000000}},
		{"three", args{3.0}, floatx80{0x4000, 0xc000000000000000}},
		{"-one", args{-1.0}, floatx80{0xbfff, 0x8000000000000000}},
		{"-two", args{-2.0}, floatx80{0xc000, 0x8000000000000000}},
		{"-three", args{-3.0}, floatx80{0xc000, 0xc000000000000000}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := float32ToFloatx80(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("float32ToFloatx80() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_float64ToFloatx80(t *testing.T) {
	type args struct {
		a float64
	}
	tests := []struct {
		name string
		args args
		want floatx80
	}{
		{"zero", args{0.0}, floatx80{0, 0}},
		{"one", args{1.0}, floatx80{0x3fff, 0x8000000000000000}},
		{"two", args{2.0}, floatx80{0x4000, 0x8000000000000000}},
		{"three", args{3.0}, floatx80{0x4000, 0xc000000000000000}},
		{"-one", args{-1.0}, floatx80{0xbfff, 0x8000000000000000}},
		{"-two", args{-2.0}, floatx80{0xc000, 0x8000000000000000}},
		{"-three", args{-3.0}, floatx80{0xc000, 0xc000000000000000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := float64ToFloatx80(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("float64ToFloatx80() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_floatx80_toInt32(t *testing.T) {
	type fields struct {
		high uint16
		low  uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   int32
	}{
		{"zero", fields{0, 0}, 0},
		{"one", fields{0x3fff, 0x8000000000000000}, 1},
		{"two", fields{0x4000, 0x8000000000000000}, 2},
		{"three", fields{0x4000, 0xc000000000000000}, 3},
		{"-one", fields{0xbfff, 0x8000000000000000}, -1},
		{"-two", fields{0xc000, 0x8000000000000000}, -2},
		{"-three", fields{0xc000, 0xc000000000000000}, -3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := floatx80{
				high: tt.fields.high,
				low:  tt.fields.low,
			}
			if got := a.toInt32(); got != tt.want {
				t.Errorf("floatx80.toInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}
