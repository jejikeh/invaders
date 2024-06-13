package goarena

import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

type Arena struct {
	data     unsafe.Pointer
	size     uintptr
	capacity uint32
}

func NewArena(capacity uint32) *Arena {
	a := &Arena{
		data:     unsafe.Pointer(&make([]byte, capacity)[0]),
		size:     0,
		capacity: capacity,
	}

	runtime.SetFinalizer(a, func(a *Arena) {
		a.Free()
	})

	return a
}

func (a *Arena) Alloc(size uintptr) (unsafe.Pointer, error) {
	if a.size+size > uintptr(a.capacity) {
		panic("arena overflow")
	}

	a.data = unsafe.Pointer(uintptr(a.data) + a.size)
	buf := unsafe.Slice((*byte)(a.data), a.size)
	for i := range buf {
		buf[i] = 0x00000000
	}
	a.size += size

	return unsafe.Pointer(unsafe.SliceData(buf)), nil
}

func (a *Arena) Free() {
	a.size = 0
	a.data = unsafe.Pointer(&make([]byte, a.capacity)[0])
}

func (a *Arena) Print() {
	fmt.Printf("Arena{data: %p, size: %d, capacity: %d}\n", a.data, a.size, a.capacity)
}

func New[T any](arena *Arena) *T {
	t := new(T)
	buf, err := arena.Alloc(reflect.Indirect(reflect.ValueOf(t)).Type().Size())
	if err != nil {
		panic(err)
	}
	t = (*T)(buf)

	return t
}

var ErrAlignmentIsNotPowerOfTwo = fmt.Errorf("alignment size is not power of two")

func alignPointer(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	if size != 0 && (size&(size-1)) != 0 {
		panic(ErrAlignmentIsNotPowerOfTwo)
	}

	offset := uintptr(ptr) & (size - 1)

	if offset != 0 {
		ptr = unsafe.Add(ptr, size-offset)
	}

	return ptr
}

func Align[T any](t T) T {
	align := uintptr(unsafe.Alignof(t))
	newPtr := alignPointer(unsafe.Pointer(&t), align)
	return *(*T)(newPtr)
}
