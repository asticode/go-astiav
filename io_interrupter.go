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
	Resume()
}

type defaultIOInterrupter struct {
	c C.struct_AVIOInterruptCB
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
