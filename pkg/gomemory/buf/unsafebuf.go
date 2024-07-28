package buf

import (
	"reflect"
	"unsafe"

	"github.com/jejikeh/invaders/pkg/gomemory/rtype"
)

//go:linkname runtime_mallocgc runtime.mallocgc
func runtime_mallocgc(size uintptr, typ uintptr, needzero bool) unsafe.Pointer

type UnsafeBuf[T any] struct {
	mem   []T
	tType uintptr
	index int
}

func NewUnsafe[T any](count int, ts ...T) *UnsafeBuf[T] {
	var t T
	if len(ts) > 0 {
		t = ts[0]
	} else {
		t = *new(T)
	}

	ttype := rtype.GetITab(t)
	tSize := indirectSizeof(t)
	size := tSize * uintptr(count)

	return &UnsafeBuf[T]{
		mem: ptrToSlice[T](
			runtime_mallocgc(
				size,
				ttype,
				true,
			),
			int(size/tSize),
		),
		tType: ttype,
		index: 0,
	}
}

func (b *UnsafeBuf[T]) Store(construct ...func(*T)) (*T, error) {
	if b.index >= len(b.mem) {
		return nil, ErrBufOverflow
	}

	t := &b.mem[b.index]
	b.index++

	for _, c := range construct {
		c(t)
	}

	return t, nil
}

func (b *UnsafeBuf[T]) Length() int {
	return b.index
}

func (b *UnsafeBuf[T]) Load(idx int) *T {
	if idx >= b.index {
		return nil
	}

	return &b.mem[idx]
}

func indirectSizeof(t any) uintptr {
	return reflect.Indirect(reflect.ValueOf(t)).Type().Size()
}

func ptrToSlice[T any](ptr unsafe.Pointer, count int) []T {
	var ret []T
	s := (*struct {
		ptr unsafe.Pointer
		len int
		cap int
	})(unsafe.Pointer(&ret))

	s.ptr = ptr
	s.cap = count
	s.len = s.cap

	return ret
}
