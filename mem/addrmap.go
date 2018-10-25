package mem

/// operands

type operandType struct {
    name string
    short string
    size int
    mask uint
}

var (
    Byte = &operandType{name: "Byte", short: ".b", size: 1, mask: 0xff}
    Word = &operandType{name: "Word", short: ".w", size: 2, mask: 0xffff}
    Long = &operandType{name: "Long", short: ".l", size: 4, mask: 0xffffffff}
)

/// addresses

type (
    Address uint32
    Reader func(Address, *operandType) (int, error)
    Writer func(Address, *operandType, int) error
    
    AddressBus interface {
        read(address Address, operand *operandType ) (int, error)
        write(address Address, operand *operandType, value int) error
    }
    
    AddressBlock struct {
        start Address
        end Address
        read Reader
        write Writer
    }
)



type addressMap struct {
    blocks []*AddressBlock
    cache *AddressBlock
}

func (a addressMap) findAddressBlock(address Address) *AddressBlock {
    if address >= a.cache.start && address < a.cache.end {
        return a.cache
    }
    for _, block := range a.blocks {
        if address >= block.start && address < block.end {
            a.cache = block
            return block
        }
    }
    return nil
}

func (a addressMap) read(address Address, operand *operandType) (int, error) {
    if block := a.findAddressBlock(address); block != nil {
        if read := block.read; read != nil {
            return read(address, operand)
        } else {
            return 0, AddressError(address)
        }
    }
    return 0, BusError(address)
}

func (a addressMap) write(address Address, operand *operandType, value int) error {
    if block := a.findAddressBlock(address); block != nil {
        if write := block.write; write != nil {
            return write(address, operand, value)
        } else {
            return AddressError(address)
        }
    }
    return BusError(address)
}


func NewAddressBus(blocks ...*AddressBlock) AddressBus {
    return addressMap{blocks: blocks, cache: blocks[0]}
}

type ram struct {
    AddressBlock
    mem []byte
}

func NewRAM(start, length uint) AddressBlock {
    return ram{
        start: start,
        end: start + length,
        mem: [length]byte,
    }
}

