package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexString(t *testing.T) {
	str1 := Byte.HexString(0x01)
	assert.Equal(t, "$01", str1)
	str2 := Word.HexString(0x01)
	assert.Equal(t, "$0001", str2)
	str3 := Long.HexString(0x01)
	assert.Equal(t, "$00000001", str3)
}

func TestSignedHexString(t *testing.T) {
	str1 := Byte.SignedHexString(-0x01)
	assert.Equal(t, "-$01", str1)
	str2 := Word.SignedHexString(-0x01)
	assert.Equal(t, "-$0001", str2)
	str3 := Long.SignedHexString(-0x01)
	assert.Equal(t, "-$00000001", str3)
	str4 := Long.SignedHexString(0x01)
	assert.Equal(t, "$00000001", str4)
}

func TestSignedHexStringMsb(t *testing.T) {
	str1 := Byte.SignedHexString(0x80)
	assert.Equal(t, "$80", str1)
	str2 := Word.SignedHexString(0x8000)
	assert.Equal(t, "$8000", str2)
	str3 := Long.SignedHexString(0x80000000)
	assert.Equal(t, "$80000000", str3)
}
