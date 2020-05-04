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

// Build routine starter for M68K
func (cpu *M68K) Build() *M68K {
	if cpu.bus == nil {
		panic("bus must not be nil")
	}
	cpu.Reset()
	return cpu
}

// NewAddressBusBuilder returns a new AddressBusBuilder
func NewAddressBusBuilder() AddressBusBuilder {
	return make(addressAreaTable, 0x10000)
}

// AddArea adds an address area to the address bus
func (aat addressAreaTable) AddArea(address, size int32, area *AddressArea) AddressBusBuilder {
	if area == nil {
		panic("AdressArea must not be nil")
	}
	if address < 0 || address&0xffff != 0 {
		panic("adddress must not be negative and multiple of 0x10000")
	}
	if size <= 0 || size&0xffff != 0 {
		panic("size must be multiple of 0x1000 and not less or equal 0")
	}

	handler := &addressAreaHandler{area, address, size}
	for i := address >> 16; i < (address+size)>>16; i++ {
		if aat[i] != nil {
			panic("address area already in use")
		}
		aat[i] = handler
	}

	return aat
}

func (aat addressAreaTable) Build() AddressBus {
	return aat
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
