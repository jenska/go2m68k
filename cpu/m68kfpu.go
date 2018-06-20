package cpu

type bits16 uint16
type bits32 float32
type bits64 float64

type floatx80 struct {
	high bits16
	low  bits64
}
