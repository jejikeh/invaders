package gomemory

import (
	"github.com/jejikeh/invaders/pkg/gomemory/arena"
)

// @Cleanup: The pool memory can be corrupted, since it embed arena,
// but arena allocator can allocate any type of object, and pool cannot.
// This can be prevented to create check in pool.Allocate for correct item size passed to method

type Pool[T any] struct {
	*arena.Buf[T]
	items map[int]int
}

func NewPool[T any](count int, ts ...T) *Pool[T] {
	return &Pool[T]{
		Buf:   arena.NewBuf(count, ts...),
		items: make(map[int]int),
	}
}

func (p *Pool[T]) StoreAt(idx int, construct ...func(*T)) *T {
	t := p.Store(construct...)
	p.items[idx] = p.Length() - 1

	return t
}

func (p *Pool[T]) LoadAt(idx int) (*T, bool) {
	bufIdx, ok := p.items[idx]

	return p.Load(bufIdx), ok
}
