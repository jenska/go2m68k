package cpu

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

const (
	romTop      int32 = 0xfc0000
	initialSSP  int32 = 0x1000 // 0x602e0104
	initialPC   int32 = 0xfc0030
	initialSize       = 1024 * 1024

	gnuAs      = "m68k-linux-gnu-as"
	gnuObjDump = "m68k-linux-gnu-objdump"
	gnuOjCopy  = "m68k-linux-gnu-objcopy"
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

func assemble(t *testing.T, asmFile string) []byte {
	tempDir := os.TempDir()
	sourcePath := fmt.Sprintf("testdata/%s.s", asmFile)
	objectPath := fmt.Sprintf("%s/%s.o", tempDir, asmFile)
	binPath := fmt.Sprintf("%s/%s.bin", tempDir, asmFile)

	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	var binaryModificationTime int64
	binaryInfo, err := os.Stat(binPath)
	if err == nil {
		binaryModificationTime = binaryInfo.ModTime().UnixNano()
	}
	if binaryModificationTime < sourceInfo.ModTime().UnixNano() {
		// fmt.Println("Assembler run")
		errOut := bytes.Buffer{}
		gnuAsCmd := exec.Command(gnuAs, sourcePath, fmt.Sprintf("-o%s", objectPath))
		gnuAsCmd.Stderr = &errOut
		if err := gnuAsCmd.Run(); err != nil {
			t.Fatalf("%s: %q\n", err, errOut.String())
		}

		gnuOjCopyCmd := exec.Command(gnuOjCopy, "-Obinary", objectPath, binPath)
		gnuOjCopyCmd.Stderr = &errOut
		if err := gnuOjCopyCmd.Run(); err != nil {
			t.Fatalf("%s: %q\n", err, errOut.String())
		}
	}
	res, err := ioutil.ReadFile(binPath)
	if err != nil {
		t.Fatal(err)
	}

	copy(trom.raw, res)

	return res
}

func _TestEnv(t *testing.T) {
	res := assemble(t, "imisc_test")
	fmt.Println(hex.Dump(res))
}

func TestMain(m *testing.M) {
	tram = NewBaseArea(initialSSP, initialPC, initialSize)
	trom = NewROMArea(make([]byte, 1024*1024))
	tbus = NewAddressBusBuilder().AddArea(0, initialSize, tram).AddArea(romTop, 1024*1024, trom).Build()

	builder := NewBuilder()
	builder.SetBus(tbus)
	builder.SetISA68000()
	tcpu = builder.Build()
	os.Exit(m.Run())
}
