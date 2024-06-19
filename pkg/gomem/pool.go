package gomem

import (
	"unsafe"
)

// @Cleanup: The pool memory can be corrupted, since it embed arena,
// but arena allocator can allocate any type of object, and pool cannot.
// This can be prevented to create check in pool.Allocate for correct item size passed to method

type Pool struct {
	*MallocArena

	itemSize     int
	items        map[int]unsafe.Pointer
	indirectSize uintptr
	alignof      uintptr
}

func NewPool[T any](capacity int) *Pool {
	return &Pool{
		MallocArena:  NewMallocArena(SizeOfAligned[T](capacity)),
		itemSize:     SizeOfAligned[T](1),
		items:        make(map[int]unsafe.Pointer, capacity),
		indirectSize: indirectSize(new(T)),
		alignof:      unsafe.Alignof(new(T)),
	}
}

func (p *Pool) NewAt(i int) unsafe.Pointer {
	buf := p.Alloc(p.indirectSize, p.alignof)
	p.items[i] = p.cursor

	return buf
}

func (p *Pool) GetAt(i int) (unsafe.Pointer, bool) {
	item, ok := p.items[i]

	return item, ok
}

func (p *Pool) Length() int {
	return len(p.items)
}

func ToTypedPool[T any](pool *Pool) *TypedPool[T] {
	if SizeOfAligned[T](1) != pool.itemSize {
		return nil
	}

	return &TypedPool[T]{
		Pool: pool,
	}
}

type TypedPool[T any] struct {
	*Pool
}

func NewTypedPool[T any](capacity int) *TypedPool[T] {
	return &TypedPool[T]{
		Pool: NewPool[T](capacity),
	}
}

func (p *TypedPool[T]) NewAt(i int) *T {
	return (*T)(p.Pool.NewAt(i))
}

func (p *TypedPool[T]) GetAt(i int) (*T, bool) {
	t, ok := p.Pool.GetAt(i)

	return (*T)(t), ok
}
