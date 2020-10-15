package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFPUTable(t *testing.T) {
	assert.Equal(t, float_mnemonics[33], "fmod")
	assert.Equal(t, float_mnemonics[48], "fsincos")
	assert.Equal(t, float_mnemonics[56], "fcmp")
}
