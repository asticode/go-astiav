// To be able to use vasprintf in Cygwin
#define _GNU_SOURCE

#include "log.h"
#include <libavutil/log.h>
#include <stdio.h>

void astiavLogCallback(void *ptr, int level, const char *fmt, va_list vl)
{
	if (level > av_log_get_level()) return;
	char* msg = NULL;
	vasprintf(&msg, fmt, vl);
	goAstiavLogCallback(ptr, level, (char*)(fmt), msg);
}

void astiavSetLogCallback()
{
	av_log_set_callback(astiavLogCallback);
}

void astiavResetLogCallback()
{
	av_log_set_callback(av_log_default_callback);
}

void astiavLog(void* ptr, int level, const char *fmt, char* arg)
{
	av_log(ptr, level, fmt, arg);
}