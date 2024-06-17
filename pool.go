package gomemory

/*
#cgo CFLAGS: -g -Wall
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
	"runtime"
)

type Pool[T any] struct {
	start  unsafe.Pointer
	end    unsafe.Pointer
	elementSize int
}

func NewPool[T any](size int) *Pool[T] {
	start := C.malloc(C.size_t(size))
	end := unsafe.Add(start, size)

	p := &Pool[T]{
		start: start,
		end: end,
		elementSize: AlignedSizeOf[T](1),
	}
	
	runtime.SetFinalizer(p, func(p *Pool[T]) {
		C.free(p.start)
	})
	
	return p
}

func (p *Pool[T]) Alloc(size uintptr, align uintptr, index int) unsafe.Pointer {
	if !isPowerOfTwo(align) {
		panic(ErrAlignmentIsNotPowerOfTwo)
	}
	
	memoryPtr := ((uintptr(p.end) - (uintptr(index+1) * size))) & ^(align - 1)
	if memoryPtr < uintptr(p.start) {
		panic(ErrArenaOverflow)
	}
	
	return unsafe.Pointer(memoryPtr)
}

func (p *Pool[T]) New(index int) *T {
	t := new(T)
	memoryPtr := p.Alloc(indirectSize(t), unsafe.Alignof(t), index)
	return (*T)(memoryPtr)
}

func (p *Pool[T]) Get(index int) *T {
	t := new(T)
	memoryPtr := unsafe.Pointer(((uintptr(p.end) - (uintptr(index+1) * indirectSize(t)))) & ^(unsafe.Alignof(t) - 1))
	return (*T)(memoryPtr)
}