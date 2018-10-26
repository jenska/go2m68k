package mem

func NewRAM(start Address, size uint) AddressArea {
    ram := make([]byte, size)
    end := start + Address(size)
    
    return AddressArea{
        start: start,
        end: end,
        write: func(address Address, operand *Operand, value int) error {
            if address >= start && address < end {
                operand.write(ram[:], uint(address - start), value) 
                return nil
            }
            return BusError(address)
        },
        read: func(address Address, operand *Operand) (int, error) {
            if address >= start && address < end {
                return operand.read(ram[:], uint(address - start)), nil
            }
            return 0, BusError(address)
        },
    }
}