package astiav

//#include "io_interrupter.h"
//#include <stdlib.h>
import "C"
import "unsafe"

type IOInterrupter struct {
	c *C.AVIOInterruptCB
	i C.int
}

func NewIOInterrupter() *IOInterrupter {
	i := &IOInterrupter{}
	i.c = C.astiavNewInterruptCallback(&i.i)
	return i
}

func (i *IOInterrupter) Free() {
	C.free(unsafe.Pointer(i.c))
}

func (i *IOInterrupter) Interrupt() {
	i.i = 1
}

func (i *IOInterrupter) Interrupted() bool {
	return i.i == 1
}

func (i *IOInterrupter) Resume() {
	i.i = 0
}
