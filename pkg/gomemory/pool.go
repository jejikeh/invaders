package gomemory

import (
	"github.com/jejikeh/invaders/pkg/gomemory/buf"
)

// @Cleanup: The pool memory can be corrupted, since it embed arena,
// but arena allocator can allocate any type of object, and pool cannot.
// This can be prevented to create check in pool.Allocate for correct item size passed to method

type Pool[K comparable, V any] struct {
	*buf.Buf[V]
	items map[K]int
}

func NewPool[K comparable, V any](count int, ts ...V) *Pool[K, V] {
	return &Pool[K, V]{
		Buf:   buf.New(count, ts...),
		items: make(map[K]int, count),
	}
}

func (p *Pool[K, V]) StoreAt(idx K, construct ...func(*V)) *V {
	t, _ := p.Store(construct...)
	p.items[idx] = p.Length() - 1

	return t
}

func (p *Pool[K, V]) LoadAt(idx K) (*V, bool) {
	bufIdx, ok := p.items[idx]

	return p.Load(bufIdx), ok
}
