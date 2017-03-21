package mem

type PhysicalAddressSpace struct {
	mem []byte
}

func NewPhysicalAddressSpace(size uint32) *PhysicalAddressSpace {
	return &PhysicalAddressSpace{make([]byte, size)}
}