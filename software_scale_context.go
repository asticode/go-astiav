package astiav

//#include <libswscale/swscale.h>
import "C"
import (
	"errors"
	"unsafe"
)

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html
type SoftwareScaleContext struct {
	c *C.struct_SwsContext
}

type softwareScaleContextUpdate struct {
	dstFormat *PixelFormat
	dstH      *int
	dstW      *int
	flags     *SoftwareScaleContextFlags
	srcFormat *PixelFormat
	srcH      *int
	srcW      *int

	alphaBlend   *SoftwareScaleContextAlphaBlend
	dither       *SoftwareScaleContextDither
	gammaFlag    *int
	intent       *SoftwareScaleContextIntent
	opaque       *unsafe.Pointer
	scalerParam0 *float64
	scalerParam1 *float64
	threads      *int

	srcRange   *int
	dstRange   *int
	srcVChrPos *int
	srcHChrPos *int
	dstVChrPos *int
	dstHChrPos *int
}

// https://ffmpeg.org/doxygen/8.0/group__libsws.html#ga59cc19eff0434e7ec11676dc5e222ff3
func CreateSoftwareScaleContext(srcW, srcH int, srcFormat PixelFormat, dstW, dstH int, dstFormat PixelFormat, flags SoftwareScaleContextFlags) (*SoftwareScaleContext, error) {
	ssc := &SoftwareScaleContext{}
	ssc.c = C.sws_getContext(
		C.int(srcW),
		C.int(srcH),
		C.enum_AVPixelFormat(srcFormat),
		C.int(dstW),
		C.int(dstH),
		C.enum_AVPixelFormat(dstFormat),
		C.int(flags),
		nil, nil, nil,
	)
	if ssc.c == nil {
		return nil, errors.New("astiav: empty new context")
	}

	classers.set(ssc)
	return ssc, nil
}

// https://ffmpeg.org/doxygen/8.0/group__libsws.html#gad90b463ceeafdfd526994742f9954dbb
func (ssc *SoftwareScaleContext) Free() {
	if ssc.c != nil {
		// Make sure to clone the classer before freeing the object since
		// the C free method may reset the pointer
		c := newClonedClasser(ssc)
		C.sws_freeContext(ssc.c)
		ssc.c = nil
		// Make sure to remove from classers after freeing the object since
		// the C free method may use methods needing the classer
		if c != nil {
			classers.del(c)
		}

	}
}

var _ Classer = (*SoftwareScaleContext)(nil)

// https://ffmpeg.org/doxygen/7.1/structSwsContext.html#a6866f52574bc730833d2580abc806261
func (ssc *SoftwareScaleContext) Class() *Class {
	if ssc.c == nil {
		return nil
	}
	return newClassFromC(unsafe.Pointer(ssc.c))
}

// https://ffmpeg.org/doxygen/8.0/swscale-v2_8txt.html#a20ffff3ac1378332422b93ed70264f4c
func (ssc *SoftwareScaleContext) ScaleFrame(src, dst *Frame) error {
	return newError(C.sws_scale_frame(ssc.c, dst.c, src.c))
}

func (ssc *SoftwareScaleContext) update(u softwareScaleContextUpdate) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}

	dstW := ssc.c.dst_w
	if u.dstW != nil {
		dstW = C.int(*u.dstW)
	}

	dstH := ssc.c.dst_h
	if u.dstH != nil {
		dstH = C.int(*u.dstH)
	}

	dstFormat := C.enum_AVPixelFormat(ssc.c.dst_format)
	if u.dstFormat != nil {
		dstFormat = C.enum_AVPixelFormat(*u.dstFormat)
	}

	srcW := ssc.c.src_w
	if u.srcW != nil {
		srcW = C.int(*u.srcW)
	}

	srcH := ssc.c.src_h
	if u.srcH != nil {
		srcH = C.int(*u.srcH)
	}

	srcFormat := C.enum_AVPixelFormat(ssc.c.src_format)
	if u.srcFormat != nil {
		srcFormat = C.enum_AVPixelFormat(*u.srcFormat)
	}

	flags := C.int(ssc.c.flags)
	if u.flags != nil {
		flags = C.int(*u.flags)
	}

	c := C.sws_getCachedContext(
		ssc.c,
		srcW,
		srcH,
		srcFormat,
		dstW,
		dstH,
		dstFormat,
		flags,
		nil, nil, nil,
	)
	if c == nil {
		return errors.New("astiav: empty new context")
	}

	ssc.c = c

	// Apply additional publicly exposed SwsContext fields directly.
	if u.opaque != nil {
		ssc.c.opaque = *u.opaque
	}
	if u.scalerParam0 != nil {
		ssc.c.scaler_params[0] = C.double(*u.scalerParam0)
	}
	if u.scalerParam1 != nil {
		ssc.c.scaler_params[1] = C.double(*u.scalerParam1)
	}
	if u.threads != nil {
		ssc.c.threads = C.int(*u.threads)
	}
	if u.dither != nil {
		ssc.c.dither = C.SwsDither(*u.dither)
	}
	if u.alphaBlend != nil {
		ssc.c.alpha_blend = C.SwsAlphaBlend(*u.alphaBlend)
	}
	if u.gammaFlag != nil {
		ssc.c.gamma_flag = C.int(*u.gammaFlag)
	}
	if u.srcRange != nil {
		ssc.c.src_range = C.int(*u.srcRange)
	}
	if u.dstRange != nil {
		ssc.c.dst_range = C.int(*u.dstRange)
	}
	if u.srcVChrPos != nil {
		ssc.c.src_v_chr_pos = C.int(*u.srcVChrPos)
	}
	if u.srcHChrPos != nil {
		ssc.c.src_h_chr_pos = C.int(*u.srcHChrPos)
	}
	if u.dstVChrPos != nil {
		ssc.c.dst_v_chr_pos = C.int(*u.dstVChrPos)
	}
	if u.dstHChrPos != nil {
		ssc.c.dst_h_chr_pos = C.int(*u.dstHChrPos)
	}
	if u.intent != nil {
		ssc.c.intent = C.int(*u.intent)
	}

	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aef45de443b59978fd38ad1531c618574
func (ssc *SoftwareScaleContext) Flags() SoftwareScaleContextFlags {
	if ssc.c == nil {
		return 0
	}
	return SoftwareScaleContextFlags(ssc.c.flags)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aef45de443b59978fd38ad1531c618574
func (ssc *SoftwareScaleContext) SetFlags(swscf SoftwareScaleContextFlags) error {
	return ssc.update(softwareScaleContextUpdate{flags: &swscf})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a883a891c8a2d4ea7a15a3a7055f64386
func (ssc *SoftwareScaleContext) DestinationWidth() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.dst_w)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a883a891c8a2d4ea7a15a3a7055f64386
func (ssc *SoftwareScaleContext) SetDestinationWidth(i int) error {
	return ssc.update(softwareScaleContextUpdate{dstW: &i})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a7facd34608c9258dae8c2942e3dce78f
func (ssc *SoftwareScaleContext) DestinationHeight() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.dst_h)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a7facd34608c9258dae8c2942e3dce78f
func (ssc *SoftwareScaleContext) SetDestinationHeight(i int) error {
	return ssc.update(softwareScaleContextUpdate{dstH: &i})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0ff71c9ef5ab6dabf90378fa7bf836ec
func (ssc *SoftwareScaleContext) DestinationPixelFormat() PixelFormat {
	if ssc.c == nil {
		return 0
	}
	return PixelFormat(ssc.c.dst_format)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0ff71c9ef5ab6dabf90378fa7bf836ec
func (ssc *SoftwareScaleContext) SetDestinationPixelFormat(p PixelFormat) error {
	return ssc.update(softwareScaleContextUpdate{dstFormat: &p})
}

func (ssc *SoftwareScaleContext) DestinationResolution() (width int, height int) {
	if ssc.c == nil {
		return 0, 0
	}
	return int(ssc.c.dst_w), int(ssc.c.dst_h)
}

func (ssc *SoftwareScaleContext) SetDestinationResolution(w int, h int) error {
	return ssc.update(softwareScaleContextUpdate{dstW: &w, dstH: &h})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aa7dc7a4f9ec57a7c37957259a51cd920
func (ssc *SoftwareScaleContext) SourceWidth() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.src_w)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0ff71c9ef5ab6dabf90378fa7bf836ec
func (ssc *SoftwareScaleContext) SetSourceWidth(i int) error {
	return ssc.update(softwareScaleContextUpdate{srcW: &i})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0dbc8c02bd3b4cd472e07008009751ff
func (ssc *SoftwareScaleContext) SourceHeight() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.src_h)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0ff71c9ef5ab6dabf90378fa7bf836ec
func (ssc *SoftwareScaleContext) SetSourceHeight(i int) error {
	return ssc.update(softwareScaleContextUpdate{srcH: &i})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aab113373f157ee3b255ad97481af0cd9
func (ssc *SoftwareScaleContext) SourcePixelFormat() PixelFormat {
	if ssc.c == nil {
		return 0
	}
	return PixelFormat(ssc.c.src_format)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aab113373f157ee3b255ad97481af0cd9
func (ssc *SoftwareScaleContext) SetSourcePixelFormat(p PixelFormat) error {
	return ssc.update(softwareScaleContextUpdate{srcFormat: &p})
}

func (ssc *SoftwareScaleContext) SourceResolution() (int, int) {
	if ssc.c == nil {
		return 0, 0
	}
	return int(ssc.c.src_w), int(ssc.c.src_h)
}

func (ssc *SoftwareScaleContext) SetSourceResolution(w int, h int) error {
	return ssc.update(softwareScaleContextUpdate{srcW: &w, srcH: &h})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a90a448639c27486dc88b3ef4fa1252de
func (ssc *SoftwareScaleContext) Opaque() unsafe.Pointer {
	if ssc.c == nil {
		return nil
	}
	return ssc.c.opaque
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0dbc8c02bd3b4cd472e07008009751ff
func (ssc *SoftwareScaleContext) SetOpaque(p unsafe.Pointer) error {
	return ssc.update(softwareScaleContextUpdate{opaque: &p})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a3cc13e08b01c1152405b3a0e4313b255
func (ssc *SoftwareScaleContext) Threads() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.threads)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a3cc13e08b01c1152405b3a0e4313b255
func (ssc *SoftwareScaleContext) SetThreads(i int) error {
	return ssc.update(softwareScaleContextUpdate{threads: &i})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#abdb353dd741ba1ff0bfab6bcb133c682
func (ssc *SoftwareScaleContext) Dither() SoftwareScaleContextDither {
	if ssc.c == nil {
		return 0
	}
	return SoftwareScaleContextDither(ssc.c.dither)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#abdb353dd741ba1ff0bfab6bcb133c682
func (ssc *SoftwareScaleContext) SetDither(d SoftwareScaleContextDither) error {
	return ssc.update(softwareScaleContextUpdate{dither: &d})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a8179d7c46f6e0acf4541abe42d03d743
func (ssc *SoftwareScaleContext) AlphaBlend() SoftwareScaleContextAlphaBlend {
	if ssc.c == nil {
		return 0
	}
	return SoftwareScaleContextAlphaBlend(ssc.c.alpha_blend)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a8179d7c46f6e0acf4541abe42d03d743
func (ssc *SoftwareScaleContext) SetAlphaBlend(a SoftwareScaleContextAlphaBlend) error {
	return ssc.update(softwareScaleContextUpdate{alphaBlend: &a})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a28ec4f81a1e3dcce7c92f2e7be8e9bd1
func (ssc *SoftwareScaleContext) GammaFlag() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.gamma_flag)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a28ec4f81a1e3dcce7c92f2e7be8e9bd1
func (ssc *SoftwareScaleContext) SetGammaFlag(i int) error {
	return ssc.update(softwareScaleContextUpdate{gammaFlag: &i})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aebfe8c01f9ea0bb80c7ae4d433cc5062
func (ssc *SoftwareScaleContext) ScalerParams() (float64, float64) {
	if ssc.c == nil {
		return 0, 0
	}
	return float64(ssc.c.scaler_params[0]), float64(ssc.c.scaler_params[1])
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aebfe8c01f9ea0bb80c7ae4d433cc5062
func (ssc *SoftwareScaleContext) SetScalerParams(p0, p1 float64) error {
	return ssc.update(softwareScaleContextUpdate{scalerParam0: &p0, scalerParam1: &p1})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#af10043377f39ca25a7af39f214c3fdff
func (ssc *SoftwareScaleContext) SourceRange() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.src_range)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#af10043377f39ca25a7af39f214c3fdff
func (ssc *SoftwareScaleContext) SetSourceRange(i int) error {
	return ssc.update(softwareScaleContextUpdate{srcRange: &i})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#ac8b4b166f254ccb8def6ce21deead82b
func (ssc *SoftwareScaleContext) DestinationRange() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.dst_range)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#ac8b4b166f254ccb8def6ce21deead82b
func (ssc *SoftwareScaleContext) SetDestinationRange(i int) error {
	return ssc.update(softwareScaleContextUpdate{dstRange: &i})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#af1c9fa89159377a7949d9a1a9586f44f
func (ssc *SoftwareScaleContext) SourceChromaPosition() (v int, h int) {
	if ssc.c == nil {
		return 0, 0
	}
	return int(ssc.c.src_v_chr_pos), int(ssc.c.src_h_chr_pos)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#af1c9fa89159377a7949d9a1a9586f44f
func (ssc *SoftwareScaleContext) SetSourceChromaPosition(v int, h int) error {
	return ssc.update(softwareScaleContextUpdate{srcVChrPos: &v, srcHChrPos: &h})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a8d5fc2f5fa1e0f15b0a5b45b82b14c1a
func (ssc *SoftwareScaleContext) DestinationChromaPosition() (v int, h int) {
	if ssc.c == nil {
		return 0, 0
	}
	return int(ssc.c.dst_v_chr_pos), int(ssc.c.dst_h_chr_pos)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a8d5fc2f5fa1e0f15b0a5b45b82b14c1a
func (ssc *SoftwareScaleContext) SetDestinationChromaPosition(v int, h int) error {
	return ssc.update(softwareScaleContextUpdate{dstVChrPos: &v, dstHChrPos: &h})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a1c26a06608196ce7b73a34f7d20e2c13
func (ssc *SoftwareScaleContext) Intent() SoftwareScaleContextIntent {
	if ssc.c == nil {
		return 0
	}
	return SoftwareScaleContextIntent(ssc.c.intent)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a1c26a06608196ce7b73a34f7d20e2c13
func (ssc *SoftwareScaleContext) SetIntent(i SoftwareScaleContextIntent) error {
	return ssc.update(softwareScaleContextUpdate{intent: &i})
}
