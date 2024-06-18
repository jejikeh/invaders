package gomemory

import (
	"errors"
	"unsafe"
)

var ErrBitSetOverflow = errors.New("bitset overflow")

type Int interface {
	~int | ~uint | ~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64
}

type BitSet[T Int] struct {
	bits T
}

func NewBitSet[T Int]() *BitSet[T] {
	return &BitSet[T]{}
}

func (b *BitSet[T]) Set(v T) {
	checkBitOverflow[T](v)
	b.bits |= 1 << v
}

func (b *BitSet[T]) Has(v T) bool {
	checkBitOverflow[T](v)
	return b.bits&(1<<v) != 0
}

func (b *BitSet[T]) Clear(v T) {
	checkBitOverflow[T](v)
	b.bits &= ^(1 << v)
}

func checkBitOverflow[T Int](v T) {
	if unsafe.Sizeof(v)*4 < uintptr(v) {
		panic(ErrBitSetOverflow)
	}
}
