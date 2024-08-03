package buf

type Map[K comparable, V any] struct {
	*Buf[V]
	items map[K]int
}

func NewMap[K comparable, V any](count int) *Map[K, V] {
	return &Map[K, V]{
		Buf:   New[V](count),
		items: make(map[K]int, count),
	}
}

func (b *Map[K, V]) Get(at K) *V {
	bufIdx, ok := b.items[at]
	if ok {
		return b.Buf.Get(bufIdx)
	}

	t := b.Buf.New()
	b.items[at] = b.Buf.Len() - 1

	return t
}

func (b *Map[K, V]) Keys() []K {
	keys := make([]K, 0, len(b.items))

	for k := range b.items {
		keys = append(keys, k)
	}

	return keys
}
