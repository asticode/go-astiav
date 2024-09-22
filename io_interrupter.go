package astiav

//#include "io_interrupter.h"
import "C"

type IOInterrupter interface {
	Interrupt()
	Resume()
}

type defaultIOInterrupter struct {
	c C.AVIOInterruptCB
	i C.int
}

func newDefaultIOInterrupter() *defaultIOInterrupter {
	i := &defaultIOInterrupter{}
	i.c = C.astiavNewInterruptCallback(&i.i)
	return i
}

func (i *defaultIOInterrupter) Interrupt() {
	i.i = 1
}

func (i *defaultIOInterrupter) Resume() {
	i.i = 0
}
