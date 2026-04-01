#include "log.h"
#include <libavutil/log.h>
#include <stdio.h>

#ifndef ASTIAV_LOG_BUF_SIZE
#define ASTIAV_LOG_BUF_SIZE 4096
#endif

void astiavLogCallback(void *ptr, int level, const char *fmt, va_list vl)
{
	if (level > av_log_get_level()) return;
	char msg[ASTIAV_LOG_BUF_SIZE];
	vsprintf(msg, fmt, vl);
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