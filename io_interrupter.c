#include <libavformat/avio.h>
#include <libavutil/mem.h>
#include <stdatomic.h>
#include <stdlib.h>

int astiavInterruptCallback(void *ret)
{
    return atomic_load((atomic_int*)ret);
}

AVIOInterruptCB* astiavNewInterruptCallback()
{
	AVIOInterruptCB* c = av_malloc(sizeof(AVIOInterruptCB));
	atomic_int* ret = av_malloc(sizeof(atomic_int));
	atomic_store(ret, 0);
	c->callback = astiavInterruptCallback;
	c->opaque = ret;
	return c;
}

int astiavInterruptCallbackLoad(AVIOInterruptCB* c)
{
    return atomic_load((atomic_int*)(c->opaque));
}

void astiavInterruptCallbackStore(AVIOInterruptCB* c, int v)
{
    return atomic_store((atomic_int*)(c->opaque), v);
}