#include <libavutil/log.h>

char* astiavClassItemName(AVClass* c, void* ptr);
AVClassCategory astiavClassCategory(AVClass* c, void* ptr);
AVClass** astiavClassParent(AVClass* c, void* ptr);