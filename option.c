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

int astiavOptSetInt(void *obj, const char *name, long long val, int search_flags) {
    return av_opt_set_int(obj, name, val, search_flags);
}

int astiavOptSetArray(void *obj, const char *name, int search_flags, 
                      unsigned start, unsigned count, int type, const void *val) {
    return av_opt_set_array(obj, name, search_flags, start, count, type, val);
}