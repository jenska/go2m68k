package core

import (
	"fmt"
	"log"
	"strings"
)

const (
	PageShift     = 12
	PageSize      = 1 << PageShift
	Address24Mask = 0x00ffffff

	XBusError              = 2
	XAdressError           = 3
	XIllegalInstruction    = 4
	XZeroDivide            = 5
	XPrivilegeViolation    = 8
	XUnintializedInterrupt = 15
	XTrapBase              = 32

	HaltSignal = iota
	ResetSignal
)

type (
	AddressError uint32
	BusError     uint32

	Reader func(address uint32, o Operand) uint32
	Writer func(address uint32, o Operand, v uint32)

	Instruction func(*Core)

	Core struct {
		Regs []uint32 // all address and data registers D0..A7
		D    []uint32 // D0, D1 ... D7
		A    []uint32 // A0, A1 ... A7

		PC, PC0  uint32 // program counter, beginnig of the currently executed instruction
		USP, SSP uint32 // user, interrupt
		SR       StatusRegister
		SP       *uint32

		readRaw      Reader
		writeRaw     Writer
		resetHandler func()

		halted bool

		// prefetch quueue
		IRC uint16 // The most recent word prefetched from memory
		IRD uint16 // The instruction currently being executed

		instructions [0x10000]Instruction
	}
)

func (ae AddressError) Error() string {
	return fmt.Sprintf("address error at $%x", uint32(ae))
}

func (be BusError) Error() string {
	return fmt.Sprintf("bus error at $%x", uint32(be))
}

func (c *Core) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "PC %08x USP %08x SR %s\n", c.PC, c.USP, c.SR)
	fmt.Fprintf(&b, "A0 %08x %08x %08x %08x\n", c.A[0], c.A[1], c.A[2], c.A[3])
	fmt.Fprintf(&b, "A4 %08x %08x %08x %08x\n", c.A[4], c.A[5], c.A[6], c.SP)
	fmt.Fprintf(&b, "D0 %08x %08x %08x %08x\n", c.D[0], c.D[1], c.D[2], c.D[3])
	fmt.Fprintf(&b, "D4 %08x %08x %08x %08x\n", c.D[4], c.D[5], c.D[6], c.D[7])
	return b.String()
}

func (c *Core) Reset() {
	for i := range c.Regs {
		c.Regs[i] = 0
	}
	c.SR = NewStatusRegister(0x2700)
	*c.SP = c.readRaw(0, Long)
	c.SSP = *c.SP
	c.PC = c.readRaw(4, Long)
	c.PC0 = c.PC
	// Fill the prefetch queue
	c.IRC = uint16(c.readRaw(c.PC, Word))

	c.resetHandler()

	c.halted = false
}

func (c *Core) Execute(signals <-chan uint16) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case BusError:
				c.ProcessException(XBusError)
			case AddressError:
				c.ProcessException(XAdressError)
			default:
				panic(r)
			}
		}
	}()

	for !c.halted {
		select {
		case signal := <-signals:
			log.Printf("signal %d PC %08x INS %04x\n", signal, c.PC0, c.IRC)
			if signal == ResetSignal {
				c.Reset()
			} else if signal == HaltSignal {
				c.halted = true
			}
		default:
			c.PC0 = c.PC
			c.IRC = uint16(c.PopPc(Word))
			log.Printf("PC %08x INS %04x\n", c.PC0, c.IRC)
			if instruction := c.instructions[c.IRC]; instruction != nil {
				instruction(c)
			} else {
				c.ProcessException(XIllegalInstruction)
			}
		}
	}
}

func (c *Core) ProcessException(x uint16) {
	oldSR := c.SR
	if !c.SR.S {
		c.SR.S = true
		c.USP = *c.SP
		*c.SP = c.SSP
	}
	c.Push(Long, c.PC)
	c.Push(Word, uint32(oldSR.Word()))

	if xaddr := c.readRaw(uint32(x)<<2, Long); xaddr != 0 {
		c.PC = xaddr
	} else if xaddr = c.readRaw(XUnintializedInterrupt<<2, Long); xaddr != 0 {
		c.PC = xaddr
	} else {
		c.halted = true
		log.Printf("Interrupt vector not set for uninitialised interrupt vector from 0x%08x", c.PC)
		log.Print(c)
	}
}
