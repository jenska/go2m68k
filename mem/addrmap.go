package mem

type OperandType struct {
    name string
    short string
    mask int
}

var Word = OperandType("word", ".w", 0xffff)



type AddressBus interface {
}
