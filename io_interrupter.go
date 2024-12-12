package astiav

//#include "io_interrupter.h"
import "C"

type IOInterrupter interface {
	Interrupt()
	Resume()
	CB() *IOInterrupterCB
}

type IOInterrupterCB struct {
	c C.AVIOInterruptCB
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

func (i *defaultIOInterrupter) CB() *IOInterrupterCB {
	return &IOInterrupterCB{c: i.c}
}
