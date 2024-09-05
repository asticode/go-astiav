#include <libavformat/avio.h>

int astiavInterruptCallback(void *ret);
AVIOInterruptCB astiavNewInterruptCallback(int *ret);