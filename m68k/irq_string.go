// Code generated by "stringer -type=IRQ"; DO NOT EDIT.

package m68k

import "fmt"

const _IRQ_name = "IRQNoneIRQ1IRQ2IRQ3IRQ4IRQ5IRQ6IRQ7"

var _IRQ_index = [...]uint8{0, 7, 11, 15, 19, 23, 27, 31, 35}

func (i IRQ) String() string {
	if i >= IRQ(len(_IRQ_index)-1) {
		return fmt.Sprintf("IRQ(%d)", i)
	}
	return _IRQ_name[_IRQ_index[i]:_IRQ_index[i+1]]
}
