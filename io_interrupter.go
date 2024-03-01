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

type IOInterrupter interface {
	Interrupt()
}

var _ IOInterrupter = (*ioInterrupter)(nil)

type ioInterrupter struct {
	c C.struct_AVIOInterruptCB
	i C.int
}

func newIOInterrupter() *ioInterrupter {
	cb := &ioInterrupter{}
	cb.c = C.astiavNewInterruptCallback(&cb.i)
	return cb
}

func (cb *ioInterrupter) Interrupt() {
	cb.i = 1
}
