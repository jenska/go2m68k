package cpu

const (
	PageSize     = 256
	OffsetMask   = 0xff
	AddressShift = 8
)

type (
	Reader func(uint32, *Size) int
	Writer func(uint32, *Size, int)
	Reset  func()

	IOManager struct {
		cpu   *M68000
		pages []*Page
	}

	Page struct {
		protected bool
		Read      Reader
		Write     Writer
		Reset     Reset
	}
)

func NewIOManager(pages []*Page) *IOManager {
	return &IOManager{pages: pages}
}

func (io *IOManager) AttachCPU(cpu *M68000) {
	io.cpu = cpu
}

func (io *IOManager) Read(address uint32, s *Size) int {
	if index := address >> AddressShift; int(index) < len(io.pages) {
		if page := io.pages[index]; page != nil {
			if read := page.Read; read != nil {
				return read(address&OffsetMask, s)
			}
		}
	}
	if io.cpu != nil {
		io.cpu.Error <- BusErrorVector
	}
	return 0
}

func (io *IOManager) Write(address uint32, s *Size, value int) {
	if index := address >> AddressShift; int(index) < len(io.pages) {
		if page := io.pages[index]; page != nil {
			if write := page.Write; write != nil {
				if io.cpu != nil && !io.cpu.SR.S && page.protected {
					io.cpu.Error <- PrivilegeViolationError
				}
				write(address&OffsetMask, s, value)
				return
			}
		}
	}
	if io.cpu != nil {
		io.cpu.Error <- BusErrorVector
	}
}

func (io *IOManager) Reset() {
	for _, page := range io.pages {
		if page != nil {
			page.Reset()
		}
	}
}

func ProtectPage(page *Page) *Page {
	page.protected = true
	return page
}

func NewROMPage(mem []byte) *Page {
	rom := mem
	return &Page{
		Read: func(offset uint32, s *Size) int {
			return s.Read(rom[offset:])
		},
		Reset: func() {
		},
	}
}

func NewRAMPage() *Page {
	ram := make([]byte, PageSize)
	return &Page{
		Write: func(offset uint32, s *Size, value int) {
			s.Write(ram[offset:], value)
		},
		Read: func(offset uint32, s *Size) int {
			return s.Read(ram[offset:])
		},
		Reset: func() {
			for i := range ram {
				ram[i] = 0
			}
		},
	}
}
