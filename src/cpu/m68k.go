package cpu

import (
	"log"
)


type M68k struct {
	 A [8]uint32
	 D [8]int32


	log log.Logger
}





func (cpu *M68k) SetSupervisorMode(mode bool) {

}