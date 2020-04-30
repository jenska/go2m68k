package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	assert.NotNil(t, BusError.Error())
	assert.NotNil(t, AdressError.Error())
}
