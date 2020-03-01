package cpu

import "fmt"

type SuperVisorException Address

func (e SuperVisorException) Error() string {
	return fmt.Sprintf("adress error: $%X", e)
}
