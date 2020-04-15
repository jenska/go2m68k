package cpu

import (
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

const (
	romTop = 0xfc0000

	initialSSP  = 0x1000
	initialPC   = romTop
	initialSize = 1024 * 1024
)

var (
	tcpu *M68K        = nil
	trom *AddressArea = nil
)

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)

	mem, err := ioutil.ReadFile("test/etos192us.img")
	if err != nil {
		panic(err)
	}
	trom = NewROMArea("etos192us", mem)
	builder := NewCPU()
	builder.AddBaseArea(initialSSP, initialPC, initialSize)
	builder.AddArea(romTop, 3*0x10000, trom)
	builder.InitISA68000()
	tcpu = builder.Go()
	log.Debug(tcpu)
	os.Exit(m.Run())
}
