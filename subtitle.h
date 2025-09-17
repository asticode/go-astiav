#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>

// 字幕相关的包装函数
int astiavDecodeSubtitle2(AVCodecContext *avctx, AVSubtitle *sub, int *got_sub_ptr, AVPacket *avpkt);
int astiavEncodeSubtitle(AVCodecContext *avctx, uint8_t *buf, int buf_size, const AVSubtitle *sub);
void astiavSubtitleFree(AVSubtitle *sub);
AVSubtitle* astiavSubtitleAlloc(void);