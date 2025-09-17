#include <libavcodec/avcodec.h>

extern enum AVPixelFormat goAstiavCodecContextGetFormat(AVCodecContext *ctx, enum AVPixelFormat *pix_fmts, int pix_fmts_size);
enum AVPixelFormat astiavCodecContextGetFormat(AVCodecContext *ctx, const enum AVPixelFormat *pix_fmts);
void astiavSetCodecContextGetFormat(AVCodecContext *ctx);
void astiavResetCodecContextGetFormat(AVCodecContext *ctx);

// 包装函数声明用于CGO调用
int astiavCodecIsOpen(AVCodecContext *avctx);
int astiavCodecGetSupportedConfig(const AVCodecContext *avctx, const AVCodec *codec, 
                                  int config, unsigned flags, const void **out_configs, 
                                  int *out_num_configs);

// 视频尺寸对齐函数
void avcodec_align_dimensions(AVCodecContext *s, int *width, int *height);
void avcodec_align_dimensions2(AVCodecContext *s, int *width, int *height, int linesize_align[4]);

// 硬件加速函数
int avcodec_get_hw_frames_parameters(AVCodecContext *avctx, AVBufferRef *device_ref, enum AVPixelFormat hw_pix_fmt, AVBufferRef **out_frames_ref);
int avcodec_default_get_buffer2(AVCodecContext *s, AVFrame *frame, int flags);
int avcodec_default_get_encode_buffer(AVCodecContext *s, AVPacket *pkt, int flags);