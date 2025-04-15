package astiav

//#include "class.h"
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVClass.html
type Class struct {
	c   *C.AVClass
	ptr unsafe.Pointer
}

func newClassFromC(ptr unsafe.Pointer) *Class {
	if ptr == nil {
		return nil
	}
	c := (**C.AVClass)(ptr)
	if c == nil {
		return nil
	}
	return &Class{
		c:   *c,
		ptr: ptr,
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVClass.html#a5fc161d93a0d65a608819da20b7203ba
func (c *Class) Category() ClassCategory {
	return ClassCategory(C.astiavClassCategory(c.c, c.ptr))
}

// https://ffmpeg.org/doxygen/7.0/structAVClass.html#ad763b2e6a0846234a165e74574a550bd
func (c *Class) ItemName() string {
	return C.GoString(C.astiavClassItemName(c.c, c.ptr))
}

// https://ffmpeg.org/doxygen/7.0/structAVClass.html#aa8883e113a3f2965abd008f7667db7eb
func (c *Class) Name() string {
	return C.GoString(c.c.class_name)
}

// https://ffmpeg.org/doxygen/7.0/structAVClass.html#a88948c8a7c6515181771615a54a808bf
func (c *Class) Parent() *Class {
	return newClassFromC(unsafe.Pointer(C.astiavClassParent(c.c, c.ptr)))
}

func (c *Class) String() string {
	return fmt.Sprintf("%s [%s] @ %p", c.ItemName(), c.Name(), c.ptr)
}

type classerHandler struct {
	messages []string
}

func (h *classerHandler) handleLog(l LogLevel, msg string) {
	if 0 <= l && l <= LogLevelError {
		h.messages = append(h.messages, msg)
	}
}

func (h *classerHandler) newError(ret C.int) error {
	i := int(ret)
	if i >= 0 {
		return nil
	}
	msg := h.messages
	h.messages = nil
	return &loggedError{Error(ret), msg}
}

type Classer interface {
	Class() *Class
	handleLog(l LogLevel, msg string)
	newError(C.int) error
}

var _ Classer = (*UnknownClasser)(nil)

type UnknownClasser struct {
	classerHandler
	c *Class
}

func newUnknownClasser(ptr unsafe.Pointer) *UnknownClasser {
	return &UnknownClasser{c: newClassFromC(ptr)}
}

func (c *UnknownClasser) Class() *Class {
	return c.c
}

var _ Classer = (*ClonedClasser)(nil)

type ClonedClasser struct {
	classerHandler
	c *Class
}

func newClonedClasser(c Classer) *ClonedClasser {
	cl := c.Class()
	if cl == nil {
		return nil
	}
	return &ClonedClasser{c: newClassFromC(cl.ptr)}
}

func (c *ClonedClasser) Class() *Class { return c.c }

var classers = newClasserPool()

type classerPool struct {
	pm sync.Map
}

func newClasserPool() *classerPool {
	return &classerPool{}
}

func (p *classerPool) unsafePointer(c Classer) unsafe.Pointer {
	if c == nil {
		return nil
	}
	cl := c.Class()
	if cl == nil {
		return nil
	}
	return cl.ptr
}

func (p *classerPool) set(c Classer) {
	if ptr := p.unsafePointer(c); ptr != nil {
		p.pm.Store(ptr, c)
	}
}

func (p *classerPool) del(c Classer) {
	if ptr := p.unsafePointer(c); ptr != nil {
		p.pm.Delete(ptr)
	}
}

func (p *classerPool) get(ptr unsafe.Pointer) (Classer, bool) {
	if c, ok := p.pm.Load(ptr); ok {
		return c.(Classer), ok
	}
	return nil, false
}

func (p *classerPool) size() int {
	var i int
	p.pm.Range(func(key, value interface{}) bool {
		i++
		return true
	})
	return i
}
