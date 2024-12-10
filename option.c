#include <libavutil/opt.h>

int astiavOptionGet(void *obj, const char *name, const char **value, int flags)
{
    uint8_t *v = NULL;
    int ret = av_opt_get(obj, name, flags, &v);
    if (ret < 0) {
        return ret;
    }
    *value = (const char *)v;
    return 0;
}