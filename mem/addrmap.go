package mem

/// operands

type Operand struct {
    name string
    short string
    size int
    mask uint
    msb uint
}

var (
    Byte = &Operand{name: "Byte", short: ".b", size: 1, mask: 0x000000ff, msb: 0x00000080}
    Word = &Operand{name: "Word", short: ".w", size: 2, mask: 0x0000ffff, msb: 0x00008000}
    Long = &Operand{name: "Long", short: ".l", size: 4, mask: 0xffffffff, msb: 0x80000000}
)

func (o *Operand) write( slice []byte, index uint, value int) {
}

func (o *Operand) read( slice []byte, index uint) int {
    return 0
}


/// addresses

type (
    Address uint32
    MemoryReader func(Address, *Operand) (int, error)
    MemoryWriter func(Address, *Operand, int) error
    
    AddressBus interface {
        read(address Address, operand *Operand ) (int, error)
        write(address Address, operand *Operand, value int) error
    }
    
    AddressArea struct {
        start Address
        end Address
        read MemoryReader
        write MemoryWriter
    }

    addressMap struct {
        areas []AddressArea
        cache *AddressArea
    }
)

func (a *addressMap) findAddressArea(address Address) *AddressArea {
    if address >= a.cache.start && address < a.cache.end {
        return a.cache
    }
    for _, area := range a.areas {
        if address >= area.start && address < area.end {
            a.cache = &area
            return &area
        }
    }
    return nil
}

func (a *addressMap) read(address Address, operand *Operand) (int, error) {
    if area := a.findAddressArea(address); area != nil {
        if read := area.read; read != nil {
            return area.read(address, operand)
        } else {
            return 0, AddressError(address)
        }
    }
    return 0, BusError(address)
}

func (a *addressMap) write(address Address, operand *Operand, value int) error {
    if area := a.findAddressArea(address); area != nil {
        if write := area.write; write != nil {
            return area.write(address, operand, value)
        } else {
            return AddressError(address)
        }
    }
    return BusError(address)
}


func NewAddressBus(areas ...AddressArea) AddressBus {
    return &addressMap{areas: areas, cache: &areas[0]}
}
