package astiav

//#cgo pkg-config: libavformat
//#include <libavformat/avio.h>
/*
int astiavInterruptCallback(void *ret)
{
    return *((int*)ret);
}
AVIOInterruptCB astiavNewInterruptCallback(int *ret)
{
	AVIOInterruptCB c = { astiavInterruptCallback, ret };
	return c;
}
*/
import "C"

type IOInterrupter struct {
	c C.struct_AVIOInterruptCB
	i C.int
}

func newIOInterrupter() *IOInterrupter {
	cb := &IOInterrupter{}
	cb.c = C.astiavNewInterruptCallback(&cb.i)
	return cb
}

func (cb *IOInterrupter) Interrupt() {
	cb.i = 1
}
