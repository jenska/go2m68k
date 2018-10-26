package mem

import (
	"fmt"

	"github.com/jenska/atari2go/cpu"
)

type AddressError cpu.Address
type BusError cpu.Address

func (e AddressError) Error() string {
	return fmt.Sprintf("adress error: $%X", e)
}

func (e BusError) Error() string {
	return fmt.Sprintf("bus error: $%X", e)
}
