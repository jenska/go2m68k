package mem

func NewROM(start Address, rom []byte, size uint) AddressArea {
    end := start + Address(size)
    
    return AddressArea{
        start: start,
        end: end,
        read: func(address Address, operand *Operand) (int, error) {
            if address >= start && address < end {
                return operand.read(rom, uint(address - start)), nil
            }
            return 0, BusError(address)
        },
    }
}