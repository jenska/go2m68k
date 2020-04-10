package cpu

// TODO: when to re-sort the address areas for better performance?
//   do it thread safe and only if necessary
type (
	Reader func(uint32, *Size) int
	Writer func(uint32, *Size, int)
	Reset  func()
	Raw    func() *[]byte

	AddressArea struct {
		Name  string
		Read  Reader
		Write Writer
		Reset Reset
		Raw   Raw
	}

	addressAreaHandler struct {
		area      *AddressArea
		offset    uint32
		size      uint32
		accessCnt uint32
	}

	addressAreaQueue []*addressAreaHandler

	IOManager struct {
		areas *addressAreaQueue
	}
)

// Sort Impl.
func (h addressAreaQueue) Len() int {
	return len(h)
}

func (h addressAreaQueue) Less(i, j int) bool {
	return h[i].accessCnt > h[j].accessCnt // bigger is lesser
}

func (h addressAreaQueue) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *addressAreaQueue) findArea(address uint32, s *Size) (*AddressArea, uint32) {
	for _, handler := range *h {
		start := handler.offset
		end := start + handler.size
		if address >= start && address+s.size < end {
			handler.accessCnt++
			return handler.area, start
		}
	}
	return nil, 0
}

// ---------------------------------------

func NewIOManager(size uint32, baseArea *AddressArea) *IOManager {
	return (&IOManager{&addressAreaQueue{}}).AddArea(0, size, baseArea)
}

func (io *IOManager) AddArea(offset, size uint32, area *AddressArea) *IOManager {
	if area == nil {
		panic("AdressArea must not be nil")
	}
	if size == 0 {
		panic("size must not be 0")
	}
	*io.areas = append(*io.areas, &addressAreaHandler{area, offset, size, 0})
	return io
}

func (io *IOManager) Read(address uint32, s *Size) int {
	if area, offset := io.areas.findArea(address, s); area != nil {
		if read := area.Read; read != nil {
			return read(address-offset, s)
		}
	}
	panic(BusError)
}

func (io *IOManager) Write(address uint32, s *Size, value int) {
	if area, offset := io.areas.findArea(address, s); area != nil {
		if write := area.Write; write != nil {
			write(address-offset, s, value)
			return
		}
	}
	panic(BusError)
}

func (io *IOManager) Reset() {
	for _, handler := range *io.areas {
		if handler.area.Reset != nil {
			handler.area.Reset()
		}
	}
}

func NewROMArea(title string, mem []byte) *AddressArea {
	rom := mem
	return &AddressArea{
		Name: title,
		Read: func(offset uint32, s *Size) int {
			return s.Read(rom[offset:])
		},
		Raw: func() *[]byte { return &rom },
	}
}

func NewRAMArea(title string, size uint32) *AddressArea {
	ram := make([]byte, size)
	return &AddressArea{
		Name: title,
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
		Raw: func() *[]byte { return &ram },
	}
}

func NewBaseArea(title string, cpu *M68K, ssp, pc, size uint32) *AddressArea {
	ram := make([]byte, size)
	Long.Write(ram[0:], int(ssp))
	Long.Write(ram[4:], int(pc))
	sv := &cpu.SR.S

	return &AddressArea{
		Name: title,
		Write: func(offset uint32, s *Size, value int) {
			if offset < 8 {
				panic(BusError)
			}
			if !*sv && offset < 1024 {
				panic(PrivilegeViolationError)
			}
			s.Write(ram[offset:], value)
		},
		Read: func(offset uint32, s *Size) int {
			return s.Read(ram[offset:])
		},
		Reset: func() {
			for i := 8; i < len(ram); i++ {
				ram[i] = 0
			}
		},
		Raw: func() *[]byte { return &ram },
	}
}
