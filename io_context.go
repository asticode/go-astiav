package astiav

/*
#cgo pkg-config: libavformat
#include <libavformat/avformat.h>

extern int go_ioctx_proxy_read(void*, uint8_t*, int);
extern int go_ioctx_proxy_write(void*, uint8_t*, int);
extern int64_t go_ioctx_proxy_seek(void*, int64_t, int);
*/
import "C"
import (
	"io"
	"sync"
	"unsafe"
)

var (
	ioContextLock      = sync.Mutex{}
	ioContextCallbacks = make([]*ioContextCbs, 0)
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avio.h#L161
type IOContext struct {
	c *C.struct_AVIOContext
}

const (
	// Seek constants
	IOContextSeekSize  = C.AVSEEK_SIZE
	IOContextSeekForce = C.AVSEEK_FORCE
	// internal seek flags
	seekableNormalFlag = C.AVIO_SEEKABLE_NORMAL
	seekableTimeFlag   = C.AVIO_SEEKABLE_TIME
)

func NewIOContext() *IOContext {
	return &IOContext{}
}

type IOContextReadFunc func(buf []byte) int
type IOContextWriteFunc func(buf []byte) int
type IOContextSeekFunc func(offset int64, whence int) int64

// AllocIOContextCallback - create IOContext with custom callbacks
func AllocIOContextCallback(
	readCb IOContextReadFunc,
	writeCb IOContextWriteFunc,
	seekCb IOContextSeekFunc,
) *IOContext {
	wf := C.int(0)
	if writeCb != nil {
		wf = C.int(1)
	}

	id := addIOCallback(
		&ioContextCbs{
			readCb, writeCb, seekCb,
		},
	)
	id_c := C.int(id)

	ctx := C.avio_alloc_context(
		nil, C.int(0), wf, unsafe.Pointer(&id_c),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_read)),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_write)),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_seek)),
	)

	ctx.direct = 1

	return &IOContext{
		c: ctx,
	}
}

// AllocIOContextReader - create IOContext for reading
func AllocIOContextReader(
	rdr io.Reader,
) *IOContext {
	return AllocIOContextReaderAndSeeker(rdr, nil)
}

// AllocIOContextReadSeeker - create IOContext for reading
func AllocIOContextReadSeeker(
	rskr io.ReadSeeker,
) *IOContext {
	return AllocIOContextReaderAndSeeker(rskr, rskr)
}

// AllocIOContextReaderAndSeeker - create IOContext for reading and seeking
func AllocIOContextReaderAndSeeker(
	rdr io.Reader, skr io.Seeker,
) *IOContext {
	id := addIOCallback(
		&ioContextCbs{
			readCb: func(inputBuf []byte) int {
				n, err := rdr.Read(inputBuf)
				if err != nil {
					if err == io.EOF {
						return int(ErrEof)
					}
					Logf(LogLevelError, "[astiav] read error: %v", err)
					return int(ErrEio)
				}
				return n
			},
			seekCb: func(offset int64, whence int) int64 {
				if whence == IOContextSeekSize {
					whence = io.SeekEnd
				}

				if whence == IOContextSeekForce {
					whence = io.SeekStart
				}
				n, err := skr.Seek(offset, whence)
				if err != nil {
					Logf(LogLevelError, "[astiav] seek error: %s\n", err)
					return -1
				}
				// Logf(LogLevelInfo, "seek %d %d > %d\n", offset, whence, n)
				return n
			},
		},
	)
	id_c := C.int(id)

	ctx := C.avio_alloc_context(
		(*C.uchar)(C.av_malloc(1)), C.int(0), C.int(0), unsafe.Pointer(&id_c),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_read)),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_write)),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_seek)),
	)
	ctx.direct = 1

	return &IOContext{
		c: ctx,
	}
}

// AllocIOContextWriter - create IOContext for writing
func AllocIOContextWriter(wtr io.Writer) *IOContext {
	return AllocIOContextWriterAndSeeker(wtr, nil)
}

// AllocIOContextWriteSeeker - create IOContext for writing
func AllocIOContextWriteSeeker(wrskr io.WriteSeeker) *IOContext {
	return AllocIOContextWriterAndSeeker(wrskr, wrskr)
}

// AllocIOContextWriterAndSeeker - create IOContext for writing and seeking
func AllocIOContextWriterAndSeeker(wtr io.Writer, skr io.Seeker) *IOContext {
	wf := C.int(0)
	if wtr != nil {
		wf = C.int(1)
	}

	id := addIOCallback(
		&ioContextCbs{
			writeCb: func(inputBuf []byte) int {
				n, err := wtr.Write(inputBuf)
				if err != nil {
					return int(ErrEio)
				}
				return n
			},
			seekCb: func(offset int64, whence int) int64 {
				if whence == IOContextSeekSize {
					whence = io.SeekEnd
					offset = 0
				}

				if whence == IOContextSeekForce {
					whence = io.SeekStart
					offset = 0
				}
				n, err := skr.Seek(offset, whence)
				if err != nil {
					return -1
				}
				return n
			},
		},
	)
	id_c := C.int(id)
	ctx := C.avio_alloc_context(
		nil, C.int(0), wf, unsafe.Pointer(&id_c),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_read)),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_write)),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_seek)),
	)
	ctx.direct = 1

	return &IOContext{
		c: ctx,
	}
}

// AllocIOContextBufferReader - create IOContext for reading from provided buffer
func AllocIOContextBufferReader(buf []byte) *IOContext {
	var pos = 0

	id := addIOCallback(
		&ioContextCbs{
			readCb: func(inputBuf []byte) int {
				inputSize := len(inputBuf)
				if pos >= len(buf) {
					return int(ErrEof)
				}
				if (pos + len(inputBuf)) > len(buf) {
					inputSize = len(buf) - pos
				}
				copy(inputBuf[:inputSize], buf[pos:])
				return inputSize
			},
			seekCb: func(offset int64, whence int) int64 {
				switch whence {
				case io.SeekCurrent:
					pos += int(offset)
				case io.SeekEnd:
					pos = len(buf) + int(offset)
				case IOContextSeekSize:
					pos = len(buf)
					break
				default:
					pos = int(offset)
				}
				return int64(pos)
			},
		},
	)

	ctx := C.avio_alloc_context(
		nil, 0, C.int(0), unsafe.Pointer(&id),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_read)),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_write)),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_seek)),
	)

	ctx.direct = 1

	return &IOContext{
		c: ctx,
	}
}

// AllocIOContextBufferWriter - create IOContext for writing to provided buffer
func AllocIOContextBufferWriter(buf []byte) *IOContext {
	var pos = 0

	id := addIOCallback(
		&ioContextCbs{
			writeCb: func(inputBuf []byte) int {
				inputSize := len(inputBuf)
				if pos >= len(buf) {
					return int(ErrEof)
				}
				if (pos + len(inputBuf)) > len(buf) {
					return int(ErrBufferTooSmall)
				}
				copy(buf[pos:], inputBuf[:inputSize])
				pos += inputSize
				return inputSize
			},
			seekCb: func(offset int64, whence int) int64 {
				switch whence {
				case IOContextSeekSize:
					pos = len(buf) + int(offset)
					break
				default:
					pos = int(offset)
				}
				return int64(pos)
			},
		},
	)

	id_c := C.int(id)

	ctx := C.avio_alloc_context(
		nil, C.int(0), C.int(1), unsafe.Pointer(&id_c),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_read)),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_write)),
		(*[0]byte)(unsafe.Pointer(C.go_ioctx_proxy_seek)),
	)

	ctx.direct = 1

	return &IOContext{
		c: ctx,
	}
}

func newIOContextFromC(c *C.struct_AVIOContext) *IOContext {
	if c == nil {
		return nil
	}
	return &IOContext{
		c: c,
	}
}

func (ic *IOContext) Free() {
	ic.freeCbs()
	C.avio_context_free(&ic.c)
}

func (ic *IOContext) Close() error {
	if ic.c != nil {
		ic.freeCbs()
		return newError(C.avio_close(ic.c))
	}
	return nil
}

func (ic *IOContext) Closep() error {
	return ic.Close()
	if ic.c != nil {
		ic.freeCbs()
		return newError(C.avio_closep(&ic.c))
	}
	return nil
}

func (ic *IOContext) Open(filename string, flags IOContextFlags) error {
	cfi := C.CString(filename)
	defer C.free(unsafe.Pointer(cfi))
	return newError(C.avio_open(&ic.c, cfi, C.int(flags)))
}

func (ic *IOContext) Accept(client *IOContext) error {
	return newError(C.avio_accept(ic.c, &client.c))
}

func (ic *IOContext) Handshake() error {
	return newError(C.avio_handshake(ic.c))
}

func (ic *IOContext) OpenWith(filename string, flags IOContextFlags, opts *Dictionary) error {
	cfi := C.CString(filename)
	defer C.free(unsafe.Pointer(cfi))
	var copts **C.struct_AVDictionary
	if opts != nil {
		copts = &opts.c
	}
	return newError(C.avio_open2(&ic.c, cfi, C.int(flags), nil, copts))
}

func (ic *IOContext) EofReached() bool {
	return int(ic.c.eof_reached) != 0
}

func (ic *IOContext) Error() error {
	return newError(ic.c.error)
}

func (ic *IOContext) Write(b []byte) error {
	if b == nil {
		return nil
	}
	C.avio_write(ic.c, (*C.uchar)(&b[0]), C.int(len(b)))
	return ic.Error()
}

func (ic *IOContext) Flush() {
	C.avio_flush(ic.c)
}

func (ic *IOContext) Read(b []byte) (int, error) {
	sliceSize := len(b)
	ret := C.avio_read(ic.c, (*C.uchar)(&b[0]), C.int(sliceSize))
	if ret < 0 {
		return int(ret), newError(ret)
	}

	return int(ret), nil
}

func (ic *IOContext) Seekable() bool {
	return flags(ic.c.seekable).has(seekableNormalFlag)
}

func (ic *IOContext) TimeSeekable() bool {
	return flags(ic.c.seekable).has(seekableTimeFlag)
}

func (ic *IOContext) Direct() bool {
	return int(ic.c.direct) != 0
}

func (ic *IOContext) Writable() bool {
	return int(ic.c.write_flag) != 0
}

func (ic *IOContext) BytesWritten() int64 {
	return int64(ic.c.written)
}

func (ic *IOContext) BytesRead() int64 {
	return int64(ic.c.bytes_read)
}

func (ic *IOContext) Seek(offset int64, whence int) (int64, error) {
	if !ic.Seekable() {
		return 0, newError(C.int(ErrEio))
	}
	ret := C.avio_seek(ic.c, C.int64_t(offset), C.int(whence))
	if ret < 0 {
		return 0, newError(C.int(ret))
	}
	return int64(ret), nil
}

func (ic *IOContext) CurrentPosition() int64 {
	return int64(ic.c.pos)
}

func (ic *IOContext) Size() int64 {
	return int64(C.avio_size(ic.c))
}

func (ic *IOContext) freeCbs() {
	ioContextLock.Lock()
	defer ioContextLock.Unlock()
	id_c := (*C.int)(ic.c.opaque)
	if id_c != nil && len(ioContextCallbacks) > int(*id_c) {
		ioContextCallbacks[int(*id_c)] = nil
	}
}

func fetchIOCallback(id int) (*ioContextCbs, bool) {
	ioContextLock.Lock()
	defer ioContextLock.Unlock()
	if len(ioContextCallbacks) <= id {
		return nil, false
	}
	ctx := ioContextCallbacks[id]
	return ctx, ctx != nil
}

func addIOCallback(ctx *ioContextCbs) int {
	ioContextLock.Lock()
	defer ioContextLock.Unlock()
	ioContextCallbacks = append(ioContextCallbacks, ctx)
	return len(ioContextCallbacks) - 1
}
