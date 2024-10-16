#include <libavutil/samplefmt.h>
#include <stdint.h>

int astiavSamplesCopyToBuffer(uint8_t* dst, int dst_size, const uint8_t * const src_data[8], int nb_channels, int nb_samples, enum AVSampleFormat sample_fmt, int align);