package gomemory

import (
	"errors"
	"unsafe"

	"github.com/jejikeh/invaders/pkg/gomath"
)

const bitsInByte = 4

var ErrBitSetOverflow = errors.New("bitset overflow")

type BitSet[T gomath.Int] struct {
	bits T
}

func NewBitSet[T gomath.Int](values ...T) *BitSet[T] {
	bitset := &BitSet[T]{}

	for _, value := range values {
		bitset.Set(value)
	}

	return bitset
}

func (b *BitSet[T]) Set(v T) {
	checkBitOverflow[T](v)
	b.bits |= 1 << v
}

func (b *BitSet[T]) Has(v T) bool {
	checkBitOverflow[T](v)

	return b.bits&(1<<v) != 0
}

func (b *BitSet[T]) Check(mask *BitSet[T]) bool {
	return b.bits&mask.bits == mask.bits
}

func (b *BitSet[T]) Unset(v T) {
	checkBitOverflow[T](v)

	b.bits &= ^(1 << v)
}

func Sizeof[T gomath.Int](v T) int {
	return int(unsafe.Sizeof(v)) * bitsInByte
}

func checkBitOverflow[T gomath.Int](v T) {
	if unsafe.Sizeof(v)*4 < uintptr(v) {
		panic(ErrBitSetOverflow)
	}
}
