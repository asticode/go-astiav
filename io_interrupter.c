#include "_cgo_export.h"
#include <libavformat/avio.h>
#include <libavutil/mem.h>
#include <stdlib.h>

int astiavInterruptCallback(void *ret)
{
    return goAstiavAtomicLoadInt((int*)ret);
}

AVIOInterruptCB* astiavNewInterruptCallback(int *ret)
{
	AVIOInterruptCB* c = av_malloc(sizeof(AVIOInterruptCB));
	c->callback = astiavInterruptCallback;
	c->opaque = ret;
	return c;
}