// Code generated by "stringer -type=IntVector"; DO NOT EDIT.

package m68k

import "fmt"

const _IntVector_name = "IntAckSpuriousIntAckAutovector"

var _IntVector_index = [...]uint8{0, 14, 30}

func (i IntVector) String() string {
	i -= 4294967294
	if i >= IntVector(len(_IntVector_index)-1) {
		return fmt.Sprintf("IntVector(%d)", i+4294967294)
	}
	return _IntVector_name[_IntVector_index[i]:_IntVector_index[i+1]]
}
