package core

import "testing"

func Test_sr_Uint16(t *testing.T) {
	tests := []struct {
		name string
		sr   StatusRegister
		want uint16
	}{
		{"reset", NewStatusRegister(0x2700), 0x2700},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sr.Word(); got != tt.want {
				t.Errorf("sr.Uint16() = %v, want %v", got, tt.want)
			}
		})
	}
}
