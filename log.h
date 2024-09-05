#include <stdarg.h>

extern void goAstiavLogCallback(void* ptr, int level, char* fmt, char* msg);
void astiavLogCallback(void *ptr, int level, const char *fmt, va_list vl);
void astiavSetLogCallback();
void astiavResetLogCallback();
void astiavLog(void* ptr, int level, const char *fmt, char* arg);