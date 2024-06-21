package arena

import (
	"errors"
	"io"
	"reflect"
	"unsafe"
)

var ErrAlignmentIsNotPowerOfTwo = errors.New("alignment is not power of two")
var ErrBufferOverflow = errors.New("buffer overflow")

type MemoryBuffer[T any] struct {
	buffer  []byte
	objects []*T
	cursor  unsafe.Pointer
	end     unsafe.Pointer
	start   unsafe.Pointer
	tAlign  uintptr
	tSize   uintptr
}

func NewMemoryBufferAny(t any, count int) *MemoryBuffer[any] {
	return newMemoryBuffer(t, count)
}

func NewMemoryBuffer[T any](count int) *MemoryBuffer[T] {
	return newMemoryBuffer(*new(T), count)
}

func newMemoryBuffer[T any](t T, count int) *MemoryBuffer[T] {
	b := &MemoryBuffer[T]{}

	b.buffer = make([]byte, SizeOfAligned[T](count))
	b.objects = make([]*T, 0)
	b.start = unsafe.Pointer(&b.buffer[0])
	b.end = unsafe.Pointer(&b.buffer[len(b.buffer)-1])
	b.cursor = b.end
	b.tAlign = unsafe.Alignof(t)
	b.tSize = unsafe.Sizeof(t)

	return b
}

func (b *MemoryBuffer[T]) New() (*T, error) {
	b.cursor = unsafe.Pointer((uintptr(b.cursor) - b.tSize) &^ (b.tAlign - 1))
	// b.cursor = unsafe.Pointer((uintptr(b.cursor) - size))

	if uintptr(b.cursor) < uintptr(b.start) {
		return nil, ErrBufferOverflow
	}

	t := (*T)(b.cursor)
	b.objects = append(b.objects, t)

	return t, nil
}

func (b *MemoryBuffer[T]) At(i int) *T {
	return b.objects[i]
}

func (b *MemoryBuffer[T]) Dump(w io.Writer) (int, error) {
	return w.Write(unsafe.Slice((*byte)(b.cursor), b.size()))
}

func (b *MemoryBuffer[T]) size() uintptr {
	return uintptr(b.end) - uintptr(b.cursor)
}

func SizeOfAligned[T any](count int) int {
	t := new(T)
	size := int(indirectSize(t))
	align := int(unsafe.Alignof(t))
	alignedSize := size
	for range count {
		aligned := (alignedSize + align - 1) &^ (align - 1)
		alignedSize = aligned + size
	}

	return alignedSize
}

func indirectSize[T any](t T) uintptr {
	return reflect.Indirect(reflect.ValueOf(t)).Type().Size()
}
