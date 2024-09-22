#include <libavcodec/avcodec.h>

extern enum AVPixelFormat goAstiavCodecContextGetFormat(AVCodecContext *ctx, enum AVPixelFormat *pix_fmts, int pix_fmts_size);
enum AVPixelFormat astiavCodecContextGetFormat(AVCodecContext *ctx, const enum AVPixelFormat *pix_fmts);
void astiavSetCodecContextGetFormat(AVCodecContext *ctx);
void astiavResetCodecContextGetFormat(AVCodecContext *ctx);