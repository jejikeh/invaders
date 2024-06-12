package goarena

import (
	"fmt"
	"unsafe"
)

var ErrOutOfMemory = fmt.Errorf("out of memory")

type Arena struct {
	data     unsafe.Pointer
	size     uintptr
	capacity uintptr
}

func NewArena(capacity uintptr) *Arena {
	return &Arena{
		data:     unsafe.Pointer(&make([]byte, capacity)[0]),
		size:     0,
		capacity: capacity,
	}
}

func (a *Arena) New(size uintptr) (unsafe.Pointer, error) {
	if a.size+size > a.capacity {
		return nil, ErrOutOfMemory
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

func New[T any](arena *Arena) *T {
	t := new(T)
	buf, err := arena.New(unsafe.Sizeof(t))
	if err != nil {
		panic(err)
	}
	t = (*T)(buf)

	return t
}
