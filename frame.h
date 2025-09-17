#include <libavcodec/avcodec.h>
#include <libavutil/samplefmt.h>
#include <stdint.h>

int astiavSamplesCopyToBuffer(uint8_t* dst, int dst_size, const uint8_t * const src_data[8], int nb_channels, int nb_samples, enum AVSampleFormat sample_fmt, int align);
int astiavFillAudioFrame(AVFrame *frame, int nb_channels, enum AVSampleFormat sample_fmt, const uint8_t *buf, int buf_size, int align);