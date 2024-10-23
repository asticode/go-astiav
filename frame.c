#include <errno.h>
#include <libavutil/avutil.h>
#include <libavutil/samplefmt.h>
#include <stdint.h>
#include <string.h>

int astiavSamplesCopyToBuffer(uint8_t* dst, int dst_size, const uint8_t * const src_data[8], int nb_channels, int nb_samples, enum AVSampleFormat sample_fmt, int align) {
    int linesize, buffer_size, nb_planes, i;
    
    buffer_size = av_samples_get_buffer_size(&linesize, nb_channels, nb_samples, sample_fmt, align);
    if (buffer_size > dst_size || buffer_size < 0) return AVERROR(EINVAL);

    nb_planes = buffer_size / linesize;

    for (i = 0; i < nb_planes; i++) {
        const uint8_t *src = src_data[i];
        memcpy(dst, src, linesize);
        dst += linesize;
        src += linesize;
    }
    return buffer_size;
}