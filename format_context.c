#include "format_context.h"

int astiavFormatContextIOOpen(AVFormatContext *s, AVIOContext **pb, const char *url, int flags, AVDictionary **options)
{
    return goAstiavFormatContextIOOpen(s, pb, (char*)(url), flags, options);
}

int astiavFormatContextIOClose(AVFormatContext *s, AVIOContext *pb)
{
    return goAstiavFormatContextIOClose(s, pb);
}