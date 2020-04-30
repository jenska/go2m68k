package cpu

import "sort"

// TODO: when to re-sort the address areas for better performance?
//   do it thread safe and only if necessary
type (
	addressAreaHandler struct {
		area      *AddressArea
		offset    int32
		size      int32
		accessCnt int32
	}

	addressAreaQueue []*addressAreaHandler
)

// Sort Impl.
func (aaq addressAreaQueue) Len() int {
	return len(aaq)
}

func (aaq addressAreaQueue) Less(i, j int) bool {
	return aaq[i].accessCnt > aaq[j].accessCnt // bigger is lesser
}

func (aaq addressAreaQueue) Swap(i, j int) {
	aaq[i], aaq[j] = aaq[j], aaq[i]
}

func (aaq *addressAreaQueue) findArea(address int32, s *Size) (*AddressArea, int32) {
	for _, handler := range *aaq {
		start := handler.offset
		end := start + handler.size
		if address >= start && address+s.size < end {
			handler.accessCnt++
			return handler.area, start
		}
	}
	return nil, 0
}

// Read returns a value at address or panics otherwise with a BusError
func (aaq *addressAreaQueue) read(address int32, s *Size) int32 {
	if area, offset := aaq.findArea(address, s); area != nil {
		if read := area.read; read != nil {
			return read(address-offset, s)
		}
	}
	panic(BusError)
}

// Write writes a value to address or panics a BusError
func (aaq *addressAreaQueue) write(address int32, s *Size, value int32) {
	if area, offset := aaq.findArea(address, s); area != nil {
		if write := area.write; write != nil {
			write(address-offset, s, value)
			return
		}
	}
	panic(BusError)
}

// Reset all areas
func (aaq *addressAreaQueue) reset() {
	sort.Sort(aaq)
	for _, handler := range *aaq {
		if handler.area.reset != nil {
			handler.area.reset()
		}
	}
}
