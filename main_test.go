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
