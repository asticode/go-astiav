int astiavOptionGet(void *obj, const char *name, const char **value, int flags);
int astiavOptSetInt(void *obj, const char *name, long long val, int search_flags);
int astiavOptSetArray(void *obj, const char *name, int search_flags, 
                      unsigned start, unsigned count, int type, const void *val);