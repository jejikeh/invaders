package gomemory

/*
#cgo CFLAGS: -g -Wall
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"unsafe"
)

// @Incomplete: Do not panic, return error.
// @Cleanup: Create new structure like MemoryBuffer and move start and end to it.
// Because i use the same structure in pool. And the MemoryBuffer need to be finalized.
var (
	ErrAlignmentIsNotPowerOfTwo = errors.New("alignment size is not power of two")
	ErrArenaOverflow            = errors.New("arena overflow")
)

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
	// @Cleanup: Better error handling here.
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

func SizeOfAligned[T any](count int) int {
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

func (b *MallocArena) WriteRawMemory(w io.Writer) (int, error) {
	buf := unsafe.Slice((*byte)(b.cursor), b.size())

	n, err := w.Write(buf)
	if err != nil {
		return n, fmt.Errorf("failed to write buffer to writed: %w", err)
	}

	return n, nil
}

func indirectSize[T any](t T) uintptr {
	return reflect.Indirect(reflect.ValueOf(t)).Type().Size()
}

func isPowerOfTwo[T int | uint | uintptr](x T) bool {
	return x != 0 && x&(x-1) == 0
}
