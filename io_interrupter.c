#include <libavformat/avio.h>

int astiavInterruptCallback(void *ret)
{
    return *((int*)ret);
}

AVIOInterruptCB astiavNewInterruptCallback(int *ret)
{
	AVIOInterruptCB c = { astiavInterruptCallback, ret };
	return c;
}