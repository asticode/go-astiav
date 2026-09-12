package astiav

//#include "io_interrupter.h"
//#include <libavutil/mem.h>
//#include <stdlib.h>
import "C"
import "unsafe"

type IOInterrupter struct {
	c *C.AVIOInterruptCB
}

func NewIOInterrupter() *IOInterrupter {
	return &IOInterrupter{c: C.astiavNewInterruptCallback()}
}

func (i *IOInterrupter) Free() {
	if i.c != nil {
		C.av_free(unsafe.Pointer(i.c.opaque))
		C.av_free(unsafe.Pointer(i.c))
		i.c = nil
	}
}

func (i *IOInterrupter) Interrupt() {
	C.astiavInterruptCallbackStore(i.c, 1)
}

func (i *IOInterrupter) Interrupted() bool {
	return C.astiavInterruptCallbackLoad(i.c) == 1
}

func (i *IOInterrupter) Resume() {
	C.astiavInterruptCallbackStore(i.c, 0)
}
