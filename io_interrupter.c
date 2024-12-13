#include <libavformat/avio.h>
#include <stdlib.h>

int astiavInterruptCallback(void *ret)
{
    return *((int*)ret);
}

AVIOInterruptCB* astiavNewInterruptCallback(int *ret)
{
	AVIOInterruptCB* c = malloc(sizeof(AVIOInterruptCB));
	c->callback = astiavInterruptCallback;
	c->opaque = ret;
	return c;
}