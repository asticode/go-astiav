package astiav

//#include "io_interrupter.h"
//#include <libavutil/mem.h>
//#include <stdlib.h>
import "C"
import (
	"sync/atomic"
	"unsafe"
)

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
	if i.c != nil {
		C.av_free(unsafe.Pointer(i.c))
		i.c = nil
	}
}

func (i *IOInterrupter) Interrupt() {
	atomic.StoreInt32((*int32)(&i.i), 1)
}

func (i *IOInterrupter) Interrupted() bool {
	return atomic.LoadInt32((*int32)(&i.i)) > 0
}

func (i *IOInterrupter) Resume() {
	atomic.StoreInt32((*int32)(&i.i), 0)
}

//export goAstiavAtomicLoadInt
func goAstiavAtomicLoadInt(i *C.int) C.int {
	return C.int(atomic.LoadInt32((*int32)(i)))
}
