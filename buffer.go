package astiav

// #include <libavutil/buffer.h>
import "C"

func isBufferWritable(buf *C.AVBufferRef) bool {
	return C.av_buffer_is_writable(buf) != 0
}
