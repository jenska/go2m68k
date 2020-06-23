package cpu

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const (
	romTop      int32 = 0xfc0000
	initialSSP  int32 = 0x1000 // 0x602e0104
	initialPC   int32 = 0xfc0030
	initialSize       = 1024 * 1024
)

var (
	tcpu *M68K        = nil
	tbus AddressBus   = nil
	trom *AddressArea = nil
	tram *AddressArea = nil
)

func twrite(opcodes ...uint16) {
	for _, opcode := range opcodes {
		tcpu.write(tcpu.pc, Word, int32(opcode))
		tcpu.pc += 2
	}
}

func trun(start int32) {
	twrite(0x4e72, 0x2700) // stop #$27000
	tcpu.write(int32(IllegalInstruction)<<2, Long, tcpu.pc)
	twrite(0x7e00 + uint16(IllegalInstruction)) // moveq #IllegalInstruction, d7
	twrite(0x4e73)                              // rte
	tcpu.write(int32(PrivilegeViolationError)<<2, Long, tcpu.pc)
	twrite(0x7e00 + uint16(PrivilegeViolationError)) // moveq #PrivilegeViolationError, d7
	twrite(0x4e73)                                   // rte
	tcpu.write(int32(UnintializedInterrupt)<<2, Long, tcpu.pc)
	twrite(0x7e00 + uint16(UnintializedInterrupt)) // moveq #UnintializedInterrupt, d7
	twrite(0x4e73)                                 // rte

	tcpu.pc = start
	signals := make(chan Signal)
	tcpu.Run(signals)
}

func TestMain(m *testing.M) {
	mem, err := ioutil.ReadFile("testdata/etos192us.img")
	if err != nil {
		panic(err)
	}
	tram = NewBaseArea(initialSSP, initialPC, initialSize)
	trom = NewROMArea(mem)
	tbus = NewAddressBusBuilder().AddArea(0, initialSize, tram).AddArea(romTop, int32(len(mem)), trom).Build()

	builder := NewBuilder()
	builder.SetBus(tbus)
	builder.SetISA68000()
	tcpu = builder.Build()
	log.Println(tcpu)
	result := m.Run()
	log.Println(tcpu)
	os.Exit(result)
}
