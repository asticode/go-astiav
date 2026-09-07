#include <libavformat/avformat.h>

extern int goAstiavFormatContextIOOpen(AVFormatContext *s, AVIOContext **pb, char *url, int flags, AVDictionary **options);
extern int goAstiavFormatContextIOClose(AVFormatContext *s, AVIOContext *pb);

int astiavFormatContextIOOpen(AVFormatContext *s, AVIOContext **pb, const char *url, int flags, AVDictionary **options);
int astiavFormatContextIOClose(AVFormatContext *s, AVIOContext *pb);