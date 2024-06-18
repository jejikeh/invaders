package gomemory

// @Cleanup: The pool memory can be corrupted, since it embed arena, but arena allocator can allocate any type of object, and pool cannot. This can be prevented to create check in pool.Allocate for correct item size passed to method
type Pool[T any] struct {
	*MallocArena

	itemSize int
	items    map[int]*T
}

func NewPool[T any](capacity int) *Pool[T] {
	return &Pool[T]{
		MallocArena: NewMallocArena(SizeOfAligned[T](capacity)),
		itemSize:    SizeOfAligned[T](1),
		items:       make(map[int]*T, capacity),
	}
}

func (p *Pool[T]) New(i int) *T {
	t := New[T](p)
	p.items[i] = t
	return t
}

func (p *Pool[T]) Get(i int) (*T, bool) {
	item, ok := p.items[i]
	return item, ok
}
