package astiav

//#cgo pkg-config: libavutil
//#include <libavutil/log.h>
//#include <stdio.h>
//#include <stdlib.h>
/*
extern void goAstiavLogCallback(int level, char* fmt, char* msg, char* parent);

static inline void astiavLogCallback(void *avcl, int level, const char *fmt, va_list vl)
{
	if (level > av_log_get_level()) return;
	AVClass* avc = avcl ? *(AVClass **) avcl : NULL;
	char parent[1024];
	if (avc) {
		sprintf(parent, "%p", avcl);
	}
	char msg[1024];
	vsprintf(msg, fmt, vl);
	goAstiavLogCallback(level, (char*)(fmt), msg, parent);
}
static inline void astiavSetLogCallback()
{
	av_log_set_callback(astiavLogCallback);
}
static inline void astiavResetLogCallback()
{
	av_log_set_callback(av_log_default_callback);
}
static inline void astiavLog(int level, const char *fmt)
{
	av_log(NULL, level, fmt, NULL);
}
*/
import "C"
import "unsafe"

type LogLevel int

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/log.h#L162
const (
	LogLevelQuiet   = LogLevel(C.AV_LOG_QUIET)
	LogLevelPanic   = LogLevel(C.AV_LOG_PANIC)
	LogLevelFatal   = LogLevel(C.AV_LOG_FATAL)
	LogLevelError   = LogLevel(C.AV_LOG_ERROR)
	LogLevelWarning = LogLevel(C.AV_LOG_WARNING)
	LogLevelInfo    = LogLevel(C.AV_LOG_INFO)
	LogLevelVerbose = LogLevel(C.AV_LOG_VERBOSE)
	LogLevelDebug   = LogLevel(C.AV_LOG_DEBUG)
)

func (l LogLevel) String() {
	return C.GoString(C.get_level_str(C.int(l)))
}

func SetLogLevel(l LogLevel) {
	C.av_log_set_level(C.int(l))
}

type LogCallback func(l LogLevel, fmt, msg, parent string)

var logCallback LogCallback

func SetLogCallback(c LogCallback) {
	logCallback = c
	C.astiavSetLogCallback()
}

//export goAstiavLogCallback
func goAstiavLogCallback(level C.int, fmt, msg, parent *C.char) {
	if logCallback == nil {
		return
	}
	logCallback(LogLevel(level), C.GoString(fmt), C.GoString(msg), C.GoString(parent))
}

func ResetLogCallback() {
	C.astiavResetLogCallback()
}

func Log(l LogLevel, msg string) {
	msgc := C.CString(msg)
	defer C.free(unsafe.Pointer(msgc))
	C.astiavLog(C.int(l), msgc)
}
