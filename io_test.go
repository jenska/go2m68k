package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRAM(t *testing.T) {
	page1 := NewRAMPage()
	page2 := NewRAMPage()
	io := NewIOManager([]*Page{page1, page2})
	io.Write(0, Long, 0)
	assert.Equal(t, 0, io.Read(0, Long))

	io.Write(0, Long, 1)
	assert.Equal(t, 1, io.Read(0, Long))

	io.Write(4, Long, 1)
	assert.Equal(t, 1, io.Read(4, Long))
}

func TestReader(t *testing.T) {
	page1 := ProtectPage(NewRAMPage())
	page2 := NewRAMPage()
	io := NewIOManager([]*Page{page1, page2})
	io.Read(0, Word)
}
