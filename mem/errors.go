package mem

import (
"fmt"
)

type AddressError Address
type BusError Address

func (e AddressError) Error() string {
    return fmt.Sprintf("adress error: $%X", e)
}

func (e BusError) Error() string {
    return fmt.Sprintf("bus error: $%X", e)
}
