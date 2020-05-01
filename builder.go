package cpu

type (
	// AddressBusBuilder builds a AddressBus for M68K CPU
	AddressBusBuilder interface {
		Build() AddressBus
		AddArea(offset, size int32, area *AddressArea) AddressBusBuilder
	}

	// Builder builds a M68K CPU
	Builder interface {
		Build() *M68K
		SetBus(AddressBus) Builder
		SetISA68000() Builder
		// SetResetHandler(resetHandler chan<- struct{})
	}
)

// NewBuilder creates a new CPU Builder object
func NewBuilder() Builder {
	return &M68K{sr: ssr{S: true, Interrupts: 7}}
}

// SetBus sets the CPU's address bus
func (cpu *M68K) SetBus(bus AddressBus) Builder {
	if bus == nil {
		panic("bus must not be nil")
	}
	bus.read(0, Long) // check if the minimum amount of memory is available
	bus.read(4, Long)
	cpu.bus = bus
	return cpu
}

// Go routine starter for M68K
func (cpu *M68K) Build() *M68K {
	if cpu.bus == nil {
		panic("bus must not be nil")
	}
	cpu.Reset()
	return cpu
}

// NewAddressBusBuilder returns a new AddressBusBuilder
func NewAddressBusBuilder() AddressBusBuilder {
	return &addressAreaQueue{}
}

// AddArea adds an address area to the address bus
func (aaq *addressAreaQueue) AddArea(offset, size int32, area *AddressArea) AddressBusBuilder {
	if area == nil {
		panic("AdressArea must not be nil")
	}
	if offset < 0 {
		panic("offset must not be negative")
	}
	if size <= 0 {
		panic("size must not be less or equal 0")
	}

	for _, handler := range *aaq {
		start := handler.offset
		end := start + handler.size
		if offset >= start && offset < end {
			panic("address area already in use")
		}
	}
	*aaq = append(*aaq, &addressAreaHandler{area, offset, size, 0})
	return aaq
}

func (aaq *addressAreaQueue) Build() AddressBus {
	return aaq
}

// NewAddressArea creates a new AddressArea with accessors read, write and resetHandler
// Build behaviour:
// read accessor is mandatory
// Runtime behaviour:
// If a write accessor is nil, the access to the address area will panic a BusError
// If reset is nil, no reset will be perfomed
func NewAddressArea(read Reader, write Writer, reset Reset) *AddressArea {
	if read == nil {
		panic("read accessor is mandatory")
	}
	return &AddressArea{
		read:  read,
		write: write,
		reset: reset,
	}
}

// NewROMArea returns a Read-Only, Pre-filled address area
func NewROMArea(mem []byte) *AddressArea {
	if mem == nil || len(mem) == 0 {
		panic("mem must not be nil or of size 0")
	}
	rom := mem
	return NewAddressArea(
		func(offset int32, s *Size) int32 {
			return s.read(rom[offset:])
		},
		nil,
		nil)
}

// NewRAMArea returns a RAM address area
func NewRAMArea(size uint32) *AddressArea {
	ram := make([]byte, size)
	return NewAddressArea(
		func(offset int32, s *Size) int32 {
			return s.read(ram[offset:])
		},
		func(offset int32, s *Size, value int32) {
			s.write(ram[offset:], value)
		},
		func() {
			for i := range ram {
				ram[i] = 0
			}
		},
	)
}

// NewBaseArea returns a new M68K base address area starting at address 0
func NewBaseArea(ssp, pc int32, size uint32) *AddressArea {
	if size < 8 {
		panic("size must be at least 8 bytes")
	}
	ram := make([]byte, size)
	Long.write(ram[0:], ssp)
	Long.write(ram[4:], pc)

	return NewAddressArea(
		func(offset int32, s *Size) int32 {
			return s.read(ram[offset:])
		},
		func(offset int32, s *Size, value int32) {
			if offset < 8 {
				panic(BusError)
			}
			s.write(ram[offset:], value)
		},
		func() {
			for i := 8; i < len(ram); i++ {
				ram[i] = 0
			}
		},
	)
}
