#include <libavformat/avio.h>
#include <libavutil/mem.h>
#include <stdatomic.h>
#include <stdlib.h>

int astiavInterruptCallback(void *ret)
{
    return atomic_load((atomic_int*)ret);
}

AVIOInterruptCB* astiavNewInterruptCallback(atomic_int *ret)
{
	AVIOInterruptCB* c = av_malloc(sizeof(AVIOInterruptCB));
	c->callback = astiavInterruptCallback;
	c->opaque = ret;
	return c;
}