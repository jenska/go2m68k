package cpu

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Registers represents the programmer visible registers of the 68000 CPU.
type Registers struct {
	D  [8]uint32
	A  [8]uint32
	PC uint32
	SR uint16
}

// Memory is a byte addressable memory that is used by the CPU.
type Memory struct {
	data []byte
}

// Instruction represents a single opcode implementation.
type Instruction func(regs *Registers, mem *Memory) error

// CPU models the execution core and allows single-step execution.
type CPU struct {
	regs         Registers
	memory       Memory
	instructions map[uint16]Instruction
}

// New creates a new CPU with the given memory size in bytes.
func New(memorySize int) *CPU {
	c := &CPU{
		memory:       Memory{data: make([]byte, memorySize)},
		instructions: make(map[uint16]Instruction),
	}
	return c
}

// Reset resets the CPU state to a known value.
func (c *CPU) Reset(pc uint32) {
	c.regs = Registers{PC: pc}
}

// Registers returns a copy of the current register set.
func (c *CPU) Registers() Registers {
	return c.regs
}

// LoadProgram copies program bytes into memory starting at the provided address.
func (c *CPU) LoadProgram(addr uint32, code []byte) error {
	if int(addr)+len(code) > len(c.memory.data) {
		return fmt.Errorf("program does not fit into memory: 0x%x", addr+uint32(len(code)))
	}
	copy(c.memory.data[addr:], code)
	return nil
}

// RegisterInstruction adds an opcode handler to the CPU.
func (c *CPU) RegisterInstruction(opcode uint16, handler Instruction) error {
	if handler == nil {
		return errors.New("instruction handler must not be nil")
	}
	if _, exists := c.instructions[opcode]; exists {
		return fmt.Errorf("instruction 0x%04x already registered", opcode)
	}
	c.instructions[opcode] = handler
	return nil
}

// ExecuteInstruction runs an instruction without fetching it from memory. This allows
// callers to execute single instructions directly through the API.
func (c *CPU) ExecuteInstruction(opcode uint16) error {
	handler, ok := c.instructions[opcode]
	if !ok {
		return fmt.Errorf("unknown opcode 0x%04x", opcode)
	}
	return handler(&c.regs, &c.memory)
}

// Step fetches the next opcode at the program counter and executes it.
func (c *CPU) Step() error {
	opcode, err := c.fetchOpcode()
	if err != nil {
		return err
	}
	return c.ExecuteInstruction(opcode)
}

func (c *CPU) fetchOpcode() (uint16, error) {
	if int(c.regs.PC)+2 > len(c.memory.data) {
		return 0, fmt.Errorf("pc out of range: 0x%x", c.regs.PC)
	}
	opcode := binary.BigEndian.Uint16(c.memory.data[c.regs.PC:])
	c.regs.PC += 2
	return opcode, nil
}

// ReadWord returns a 16-bit value from memory at the given address.
func (m *Memory) ReadWord(addr uint32) (uint16, error) {
	if int(addr)+2 > len(m.data) {
		return 0, fmt.Errorf("read out of bounds at 0x%x", addr)
	}
	return binary.BigEndian.Uint16(m.data[addr:]), nil
}

// WriteWord writes a 16-bit value to memory at the given address.
func (m *Memory) WriteWord(addr uint32, value uint16) error {
	if int(addr)+2 > len(m.data) {
		return fmt.Errorf("write out of bounds at 0x%x", addr)
	}
	binary.BigEndian.PutUint16(m.data[addr:], value)
	return nil
}
