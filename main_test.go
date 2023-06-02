package m68k

import (
	"encoding/binary"
	"testing"
)

type TestEnv struct {
	c   CPU
	ram []byte
}

func (e *TestEnv) Read8(offset uint32) uint8 {
	return e.ram[offset]
}

func (e *TestEnv) Read16(offset uint32) uint16 {
	return binary.BigEndian.Uint16(e.ram[offset:])
}

func (e *TestEnv) Read32(offset uint32) uint32 {
	return binary.BigEndian.Uint32(e.ram[offset:])
}

func (e *TestEnv) Write8(offset uint32, v uint8) {
	e.ram[offset] = v
}

func (e *TestEnv) Write16(offset uint32, v uint16) {
	binary.BigEndian.PutUint16(e.ram[offset:], v)
}

func (e *TestEnv) Write32(offset uint32, v uint32) {
	binary.BigEndian.PutUint32(e.ram[offset:], v)
}

func buildTestEnv(m Model) *TestEnv {
	env := &TestEnv{}
	env.ram = make([]byte, 0x10000)
	offset := uint32(len(env.ram))
	bus := NewBusController(BaseRAM(0x10000, 0x10000, 0x10000), ChipArea(0x10000, offset, env, env, nil))
	env.c = New(m, bus)
	return env
}

func TestExecute(t *testing.T) {
	env := buildTestEnv(M68000)
	env.Write16(0, 0x7000)
	env.Write16(2, 0x7201)
	signals := make(chan uint16)
	env.c.Execute(signals)
	signals <- HaltSignal
}
