package astiav

//#include <libavutil/avutil.h>
//#include <errno.h>
import "C"
import (
	"fmt"
	"strings"
)

// https://ffmpeg.org/doxygen/7.0/group__lavu__error.html#ga586e134e9dad8f57a218b2cd8734b601
type Error int

const (
	ErrBsfNotFound      = Error(C.AVERROR_BSF_NOT_FOUND)
	ErrBufferTooSmall   = Error(C.AVERROR_BUFFER_TOO_SMALL)
	ErrBug              = Error(C.AVERROR_BUG)
	ErrBug2             = Error(C.AVERROR_BUG2)
	ErrDecoderNotFound  = Error(C.AVERROR_DECODER_NOT_FOUND)
	ErrDemuxerNotFound  = Error(C.AVERROR_DEMUXER_NOT_FOUND)
	ErrEagain           = Error(-(C.EAGAIN))
	ErrEio              = Error(-(C.EIO))
	ErrEncoderNotFound  = Error(C.AVERROR_ENCODER_NOT_FOUND)
	ErrEof              = Error(C.AVERROR_EOF)
	ErrEperm            = Error(-(C.EPERM))
	ErrEpipe            = Error(-(C.EPIPE))
	ErrEtimedout        = Error(-(C.ETIMEDOUT))
	ErrExit             = Error(C.AVERROR_EXIT)
	ErrExperimental     = Error(C.AVERROR_EXPERIMENTAL)
	ErrExternal         = Error(C.AVERROR_EXTERNAL)
	ErrFilterNotFound   = Error(C.AVERROR_FILTER_NOT_FOUND)
	ErrHttpBadRequest   = Error(C.AVERROR_HTTP_BAD_REQUEST)
	ErrHttpForbidden    = Error(C.AVERROR_HTTP_FORBIDDEN)
	ErrHttpNotFound     = Error(C.AVERROR_HTTP_NOT_FOUND)
	ErrHttpOther4Xx     = Error(C.AVERROR_HTTP_OTHER_4XX)
	ErrHttpServerError  = Error(C.AVERROR_HTTP_SERVER_ERROR)
	ErrHttpUnauthorized = Error(C.AVERROR_HTTP_UNAUTHORIZED)
	ErrInputChanged     = Error(C.AVERROR_INPUT_CHANGED)
	ErrInvaliddata      = Error(C.AVERROR_INVALIDDATA)
	ErrMaxStringSize    = Error(C.AV_ERROR_MAX_STRING_SIZE)
	ErrMuxerNotFound    = Error(C.AVERROR_MUXER_NOT_FOUND)
	ErrOptionNotFound   = Error(C.AVERROR_OPTION_NOT_FOUND)
	ErrOutputChanged    = Error(C.AVERROR_OUTPUT_CHANGED)
	ErrPatchwelcome     = Error(C.AVERROR_PATCHWELCOME)
	ErrProtocolNotFound = Error(C.AVERROR_PROTOCOL_NOT_FOUND)
	ErrStreamNotFound   = Error(C.AVERROR_STREAM_NOT_FOUND)
	ErrUnknown          = Error(C.AVERROR_UNKNOWN)
)

func newError(ret C.int) error {
	i := int(ret)
	if i >= 0 {
		return nil
	}
	return Error(i)
}

// https://ffmpeg.org/doxygen/7.0/group__lavu__error.html#ga5792b4a2d18d7d9cb0efbcfc335dce24
func (e Error) Error() string {
	s, _ := stringFromC(255, func(buf *C.char, size C.size_t) error {
		return newError(C.av_strerror(C.int(e), buf, size))
	})
	return s
}

func (e Error) Is(err error) bool {
	a, ok := err.(Error)
	if !ok {
		return false
	}
	return int(a) == int(e)
}

type loggedError struct {
	e   Error
	msg []string
}

func (e *loggedError) Error() string {
	return fmt.Sprintf("%s: %s", e.e, strings.Join(e.msg, ": "))
}

func (e *loggedError) Is(err error) bool {
	switch a := err.(type) {
	case *loggedError:
		return int(a.e) == int(e.e)
	case Error:
		return int(a) == int(e.e)
	default:
		return false
	}
}
