#include <libavformat/avio.h>
#include <stdatomic.h>

int astiavInterruptCallback(void *ret);
AVIOInterruptCB* astiavNewInterruptCallback(atomic_int *ret);