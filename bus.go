package cpu

// TODO: when to re-sort the address areas for better performance?
//   do it thread safe and only if necessary
type (
	addressAreaHandler struct {
		area   *AddressArea
		offset int32
		size   int32
	}

	addressAreaTable []*addressAreaHandler
)

// Read returns a value at address or panics otherwise with a BusError
func (aat addressAreaTable) read(address int32, s *Size) int32 {
	if handler := aat[address>>16]; handler != nil {
		if read := handler.area.read; read != nil {
			return read(address-handler.offset, s)
		}
	}
	panic(BusError)
}

// Write writes a value to address or panics a BusError
func (aat addressAreaTable) write(address int32, s *Size, value int32) {
	if handler := aat[address>>16]; handler != nil {
		if write := handler.area.write; write != nil {
			write(address-handler.offset, s, value)
			return
		}
	}
	panic(BusError)
}

// Reset all areas
func (aat addressAreaTable) reset() {
	var prevHandler *addressAreaHandler = nil
	for _, handler := range aat {
		if handler != nil && handler != prevHandler {
			prevHandler = handler
			if handler.area.reset != nil {
				handler.area.reset()
			}
		}
	}
}
