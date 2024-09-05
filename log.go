package astiav

//#include "log.h"
//#include <libavutil/log.h>
//#include <stdlib.h>
import "C"
import (
	"unsafe"
)

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

func SetLogLevel(l LogLevel) {
	C.av_log_set_level(C.int(l))
}

func GetLogLevel() LogLevel {
	return LogLevel(C.av_log_get_level())
}

type LogCallback func(c Classer, l LogLevel, fmt, msg string)

var logCallback LogCallback

func SetLogCallback(c LogCallback) {
	logCallback = c
	C.astiavSetLogCallback()
}

//export goAstiavLogCallback
func goAstiavLogCallback(ptr unsafe.Pointer, level C.int, fmt, msg *C.char) {
	// No callback
	if logCallback == nil {
		return
	}

	// Get classer
	var c Classer
	if ptr != nil {
		var ok bool
		if c, ok = classers.get(ptr); !ok {
			c = newUnknownClasser(ptr)
		}
	}

	// Callback
	logCallback(c, LogLevel(level), C.GoString(fmt), C.GoString(msg))
}

func ResetLogCallback() {
	C.astiavResetLogCallback()
}

func Log(c Classer, l LogLevel, fmt string, args ...string) {
	fmtc := C.CString(fmt)
	defer C.free(unsafe.Pointer(fmtc))
	argc := (*C.char)(nil)
	if len(args) > 0 {
		argc = C.CString(args[0])
		defer C.free(unsafe.Pointer(argc))
	}
	var ptr unsafe.Pointer
	if c != nil {
		if cl := c.Class(); cl != nil {
			ptr = cl.ptr
		}
	}
	C.astiavLog(ptr, C.int(l), fmtc, argc)
}
