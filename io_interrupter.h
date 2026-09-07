#include <libavformat/avio.h>
#include <stdatomic.h>

int astiavInterruptCallback(void *ret);
AVIOInterruptCB* astiavNewInterruptCallback();
int astiavInterruptCallbackLoad(AVIOInterruptCB* c);
void astiavInterruptCallbackStore(AVIOInterruptCB* c, int v);