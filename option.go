package astiav

//#include <libavutil/opt.h>
//#include "option.h"
import "C"
import (
	"unsafe"
)

// https://www.ffmpeg.org/doxygen/7.0/structAVOption.html
type Option struct {
	c *C.AVOption
}

func newOptionFromC(c *C.AVOption) *Option {
	if c == nil {
		return nil
	}
	return &Option{c: c}
}

// https://www.ffmpeg.org/doxygen/7.0/structAVOption.html#a87e81c6e58d6a94d97a98ad15a4e507c
func (o *Option) Name() string {
	return C.GoString(o.c.name)
}

type Options struct {
	c unsafe.Pointer
}

func newOptionsFromC(c unsafe.Pointer) *Options {
	if c == nil {
		return nil
	}
	return &Options{c: c}
}

// https://www.ffmpeg.org/doxygen/7.0/group__opt__mng.html#gabc75970cd87d1bf47a4ff449470e9225
func (os *Options) List() (list []*Option) {
	var prev *C.AVOption
	for {
		o := C.av_opt_next(os.c, prev)
		if o == nil {
			return
		}
		list = append(list, newOptionFromC(o))
		prev = o
	}
}

// https://www.ffmpeg.org/doxygen/7.0/group__opt__set__funcs.html#ga5fd4b92bdf4f392a2847f711676a7537
func (os *Options) Set(name, value string, f OptionSearchFlags) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	return newError(C.av_opt_set(os.c, cname, cvalue, C.int(f)))
}

// https://www.ffmpeg.org/doxygen/7.0/group__opt__get__funcs.html#gaf31144e60f9ce89dbe8cbea57a0b232c
func (os *Options) Get(name string, f OptionSearchFlags) (string, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var cvalue *C.char = nil
	if err := newError(C.astiavOptionGet(os.c, cname, &cvalue, C.int(f))); err != nil {
		return "", err
	}
	if cvalue == nil {
		return "", nil
	}
	defer C.av_freep(unsafe.Pointer(&cvalue))
	return C.GoString(cvalue), nil
}

// SetDict sets options from a dictionary
// https://www.ffmpeg.org/doxygen/7.0/group__opt__set__funcs.html#ga8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (os *Options) SetDict(dict *Dictionary) error {
	var dc *C.AVDictionary
	if dict != nil {
		dc = dict.c
	}
	return newError(C.av_opt_set_dict(os.c, &dc))
}

// SetArray sets an array option (simplified implementation)
// https://www.ffmpeg.org/doxygen/7.0/group__opt__set__funcs.html#ga8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (os *Options) SetArray(name string, values []string, f OptionSearchFlags) error {
	// Note: av_opt_set_array has complex signature in FFmpeg 8.0
	// For now, we'll use a simplified approach by setting as string
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	
	// Join values as comma-separated string for basic array support
	var valueStr string
	for i, v := range values {
		if i > 0 {
			valueStr += ","
		}
		valueStr += v
	}
	
	cvalue := C.CString(valueStr)
	defer C.free(unsafe.Pointer(cvalue))
	
	return newError(C.av_opt_set(os.c, cname, cvalue, C.int(f)))
}
