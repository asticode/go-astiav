#include "codec_context.h"
#include <libavcodec/avcodec.h>
#include <stdlib.h>

enum AVPixelFormat astiavCodecContextGetFormat(AVCodecContext *ctx, const enum AVPixelFormat *pix_fmts)
{
	int pix_fmts_size = 0;
	while (*pix_fmts != AV_PIX_FMT_NONE) {
		pix_fmts_size++;
		pix_fmts++;
	}
	pix_fmts -= pix_fmts_size;
	return goAstiavCodecContextGetFormat(ctx, (enum AVPixelFormat*)(pix_fmts), pix_fmts_size);
}
void astiavSetCodecContextGetFormat(AVCodecContext *ctx)
{
	ctx->get_format = astiavCodecContextGetFormat;
}
void astiavResetCodecContextGetFormat(AVCodecContext *ctx)
{
	ctx->get_format = NULL;
}

int astiavCodecIsOpen(AVCodecContext *avctx) {
    return avcodec_is_open(avctx);
}

int astiavCodecGetSupportedConfig(const AVCodecContext *avctx, const AVCodec *codec, 
                                  int config, unsigned flags, const void **out_configs, 
                                  int *out_num_configs) {
    return avcodec_get_supported_config(avctx, codec, config, flags, out_configs, out_num_configs);
}