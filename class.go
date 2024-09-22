package astiav

//#include "class.h"
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/log.h#L66
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

func (c *Class) Category() ClassCategory {
	return ClassCategory(C.astiavClassCategory(c.c, c.ptr))
}

func (c *Class) ItemName() string {
	return C.GoString(C.astiavClassItemName(c.c, c.ptr))
}

func (c *Class) Name() string {
	return C.GoString(c.c.class_name)
}

func (c *Class) Parent() *Class {
	return newClassFromC(unsafe.Pointer(C.astiavClassParent(c.c, c.ptr)))
}

func (c *Class) String() string {
	return fmt.Sprintf("%s [%s] @ %p", c.ItemName(), c.Name(), c.ptr)
}

type Classer interface {
	Class() *Class
}

type UnknownClasser struct {
	c *Class
}

func newUnknownClasser(ptr unsafe.Pointer) *UnknownClasser {
	return &UnknownClasser{c: newClassFromC(ptr)}
}

func (c *UnknownClasser) Class() *Class {
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
