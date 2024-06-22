package arena

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

var ErrBufOverflow = errors.New("buffer overflow")

type BufItem[T any] struct {
	value *T
}

func (w *BufItem[T]) get() *T {
	return w.value
}

type Buf[T any] struct {
	// buffer     []byte
	// byteOffset int
	ptr    unsafe.Pointer
	cursor unsafe.Pointer

	typeSize  uintptr
	typeAlign uintptr

	refs   []BufItem[T]
	length int
}

func NewBuf[T any](count int, ts ...T) *Buf[T] {
	var t T
	if len(ts) > 0 {
		t = ts[0]
	} else {
		t = *new(T)
	}

	buffer := make([]byte, int(alignedSizeof(t))*count)
	ptr := unsafe.Pointer(unsafe.SliceData(buffer))

	return &Buf[T]{
		refs:   make([]BufItem[T], count),
		ptr:    ptr,
		cursor: ptr,
		// buffer: buffer,
		// byteOffset: (int(alignedSizeof(t)) * count) - 1,
		// byteOffset: 0,
		// @Incomplete: This leads to wrong pointer calculations...
		// Maybe GC moves buffer, and cursor dosen't get updated?
		// cursor:    unsafe.Add(unsafe.Pointer(&buffer[0]), unsafe.Sizeof(t)*uintptr(count)),
		typeSize:  indirectSizeof(t),
		typeAlign: unsafe.Alignof(t),
	}
}

func (b *Buf[T]) Store(construct ...func(*T)) *T {
	// start := uintptr(unsafe.Pointer(&b.buffer[b.byteOffset]))
	// @Cleanup: The alignment is still not right.
	// cursor := unsafe.Pointer((uintptr(unsafe.Pointer(&b.buffer[b.byteOffset])) - b.typeSize) &^ (b.typeAlign - 1))
	// cursor := unsafe.Pointer(&b.buffer[b.byteOffset])
	// if b.byteOffset < 0 {
	// panic(ErrBufOverflow)
	// }

	// b.byteOffset -= int(start - uintptr(cursor))
	// @Incomplete: Fix alignment, reverse allocations.
	// ptr := unsafe.Pointer(&b.buffer[b.byteOffset])
	// cursor := unsafe.Pointer((uintptr(ptr)))

	// b.byteOffset += int(b.typeSize)

	buf := unsafe.Slice((*byte)(b.cursor), b.typeSize)
	for i := range buf {
		buf[i] = 0
	}

	b.refs[b.length] = BufItem[T]{
		value: (*T)(b.cursor),
	}
	b.length++
	b.cursor = unsafe.Add(b.cursor, b.typeSize)

	for _, c := range construct {
		c(b.refs[b.length-1].get())
	}

	return b.refs[b.length-1].get()
}

func (b *Buf[T]) Load(idx int) *T {
	return b.refs[idx].get()
}

func (b *Buf[T]) Length() int {
	return b.length
}

func (b *Buf[T]) write(w io.Writer) (int, error) {
	n, err := w.Write(unsafe.Slice((*byte)(b.ptr), int(b.typeSize)*b.length))
	if err != nil {
		return n, fmt.Errorf("failed to write buffer to writed: %w", err)
	}

	return n, nil
}

func indirectSizeof(t any) uintptr {
	return reflect.Indirect(reflect.ValueOf(t)).Type().Size()
}

func alignedSizeof(t any) uintptr {
	return (indirectSizeof(t) + unsafe.Alignof(t)) &^ (unsafe.Alignof(t) - 1)
}
