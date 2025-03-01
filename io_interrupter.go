package astiav

//#include "atomic.h"
//#include "io_interrupter.h"
//#include <libavutil/mem.h>
//#include <stdlib.h>
import "C"
import "unsafe"

type IOInterrupter struct {
	c *C.AVIOInterruptCB
	i C.atomic_int
}

func NewIOInterrupter() *IOInterrupter {
	i := &IOInterrupter{}
	i.c = C.astiavNewInterruptCallback(&i.i)
	return i
}

func (i *IOInterrupter) Free() {
	if i.c != nil {
		C.av_free(unsafe.Pointer(i.c))
		i.c = nil
	}
}

func (i *IOInterrupter) Interrupt() {
	C.astiavAtomicStoreInt(&i.i, 1)
}

func (i *IOInterrupter) Interrupted() bool {
	return C.astiavAtomicLoadInt(&i.i) == 1
}

func (i *IOInterrupter) Resume() {
	C.astiavAtomicStoreInt(&i.i, 0)
}
