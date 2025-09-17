#include "subtitle.h"
#include <stdlib.h>

// 包装avcodec_decode_subtitle2函数
int astiavDecodeSubtitle2(AVCodecContext *avctx, AVSubtitle *sub, int *got_sub_ptr, AVPacket *avpkt) {
    return avcodec_decode_subtitle2(avctx, sub, got_sub_ptr, avpkt);
}

// 包装avcodec_encode_subtitle函数
int astiavEncodeSubtitle(AVCodecContext *avctx, uint8_t *buf, int buf_size, const AVSubtitle *sub) {
    return avcodec_encode_subtitle(avctx, buf, buf_size, sub);
}

// 包装avsubtitle_free函数
void astiavSubtitleFree(AVSubtitle *sub) {
    avsubtitle_free(sub);
}

// 分配AVSubtitle结构体
AVSubtitle* astiavSubtitleAlloc(void) {
    AVSubtitle *sub = (AVSubtitle*)av_mallocz(sizeof(AVSubtitle));
    return sub;
}