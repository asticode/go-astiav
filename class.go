package astiav

//#include "class.h"
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

// https://ffmpeg.org/doxygen/8.0/structAVClass.html
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

// https://ffmpeg.org/doxygen/8.0/structAVClass.html#a5fc161d93a0d65a608819da20b7203ba
func (c *Class) Category() ClassCategory {
	return ClassCategory(C.astiavClassCategory(c.c, c.ptr))
}

// https://ffmpeg.org/doxygen/8.0/structAVClass.html#ad763b2e6a0846234a165e74574a550bd
func (c *Class) ItemName() string {
	return C.GoString(C.astiavClassItemName(c.c, c.ptr))
}

// https://ffmpeg.org/doxygen/8.0/structAVClass.html#aa8883e113a3f2965abd008f7667db7eb
func (c *Class) Name() string {
	return C.GoString(c.c.class_name)
}

// https://ffmpeg.org/doxygen/8.0/structAVClass.html#a88948c8a7c6515181771615a54a808bf
func (c *Class) Parent() *Class {
	return newClassFromC(unsafe.Pointer(C.astiavClassParent(c.c, c.ptr)))
}

func (c *Class) String() string {
	return fmt.Sprintf("%s [%s] @ %p", c.ItemName(), c.Name(), c.ptr)
}

type Classer interface {
	Class() *Class
}

var _ Classer = (*UnknownClasser)(nil)

type UnknownClasser struct {
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
	c *Class
}

func newClonedClasser(c Classer) *ClonedClasser {
	cl := c.Class()
	if cl == nil {
		return nil
	}
	return &ClonedClasser{c: newClassFromC(cl.ptr)}
}

func (c *ClonedClasser) Class() *Class {
	return c.c
}

var classers = newClasserPool()

type classerPool struct {
	m sync.Mutex
	p map[unsafe.Pointer]Classer
}

func newClasserPool() *classerPool {
	return &classerPool{p: make(map[unsafe.Pointer]Classer)}
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
	p.m.Lock()
	defer p.m.Unlock()
	if ptr := p.unsafePointer(c); ptr != nil {
		p.p[ptr] = c
	}
}

func (p *classerPool) del(c Classer) {
	p.m.Lock()
	defer p.m.Unlock()
	if ptr := p.unsafePointer(c); ptr != nil {
		delete(p.p, ptr)
	}
}

func (p *classerPool) get(ptr unsafe.Pointer) (Classer, bool) {
	p.m.Lock()
	defer p.m.Unlock()
	c, ok := p.p[ptr]
	return c, ok
}
