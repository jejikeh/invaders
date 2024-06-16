package gomemory

/*
#cgo CFLAGS: -g -Wall
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"unsafe"
)

var ErrAlignmentIsNotPowerOfTwo = fmt.Errorf("alignment size is not power of two")
var ErrArenaOverflow = fmt.Errorf("arena overflow")

type Arena interface {
	Alloc(size uintptr, align uintptr) unsafe.Pointer
	Free()
}

func New[T any](arena Arena) *T {
	t := new(T)
	buf := arena.Alloc(indirectSize(t), unsafe.Alignof(t))
	return (*T)(buf)
}

type MallocArena struct {
	start  unsafe.Pointer
	end    unsafe.Pointer
	cursor unsafe.Pointer
}

func NewMallocArena(size int) *MallocArena {
	start := C.malloc(C.size_t(size))
	end := unsafe.Add(start, size)

	m := &MallocArena{
		start:  start,
		end:    end,
		cursor: end,
	}

	runtime.SetFinalizer(m, func(m *MallocArena) {
		C.free(m.start)
	})

	return m
}

func (b *MallocArena) Alloc(size uintptr, align uintptr) unsafe.Pointer {
	if !isPowerOfTwo(align) {
		panic(ErrAlignmentIsNotPowerOfTwo)
	}

	newCursorPos := (uintptr(b.cursor) - size) & ^(align - 1)
	if newCursorPos < uintptr(b.start) {
		panic(ErrArenaOverflow)
	}
	b.cursor = unsafe.Pointer(newCursorPos)

	return b.cursor
}

func (b *MallocArena) Free() {
	b.cursor = b.end
}

func (b *MallocArena) size() uintptr {
	return uintptr(b.end) - uintptr(b.cursor)
}

func AlignedSizeOf[T any](count int) int {
	t := new(T)
	size := int(indirectSize(t))
	align := int(unsafe.Alignof(t))
	alignedSize := size
	for range count {
		aligned := (alignedSize + align - 1) & ^(align - 1)
		alignedSize = aligned + size
	}
	return alignedSize - size
}

func SizeOf[T any](count int) int {
	t := new(T)
	size := int(indirectSize(t))
	return size * count
}

func (b *MallocArena) DumpBuffer(w io.Writer) (int, error) {
	buf := unsafe.Slice((*byte)(b.cursor), b.size())

	return w.Write(buf)
}

func indirectSize[T any](t T) uintptr {
	return reflect.Indirect(reflect.ValueOf(t)).Type().Size()
}

func isPowerOfTwo[T int | uint | uintptr](x T) bool {
	return x != 0 && x&(x-1) == 0
}
