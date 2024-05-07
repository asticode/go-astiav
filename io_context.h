#include <stdint.h>

extern int goAstiavIOContextReadFunc(void *opaque, uint8_t *buf, int buf_size);
extern int64_t goAstiavIOContextSeekFunc(void *opaque, int64_t offset, int whence);
extern int goAstiavIOContextWriteFunc(void *opaque, uint8_t *buf, int buf_size);

int astiavIOContextReadFunc(void *opaque, uint8_t *buf, int buf_size);
int64_t astiavIOContextSeekFunc(void *opaque, int64_t offset, int whence);
int astiavIOContextWriteFunc(void *opaque, uint8_t *buf, int buf_size);