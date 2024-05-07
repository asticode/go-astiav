#include "io_context.h"
#include <stdint.h>

int astiavIOContextReadFunc(void *opaque, uint8_t *buf, int buf_size)
{
    return goAstiavIOContextReadFunc(opaque, buf, buf_size);
}

int64_t astiavIOContextSeekFunc(void *opaque, int64_t offset, int whence)
{
    return goAstiavIOContextSeekFunc(opaque, offset, whence);
}

int astiavIOContextWriteFunc(void *opaque, uint8_t *buf, int buf_size)
{
    return goAstiavIOContextWriteFunc(opaque, buf, buf_size);
}