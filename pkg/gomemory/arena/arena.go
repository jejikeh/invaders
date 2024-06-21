package arena

// type Arena struct {
// 	*MemoryBuffer
// 	itemSize     int
// 	items        map[int]any
// 	indirectSize uintptr
// 	alignof      uintptr
// }

// func NewArena[T any](size int) *Arena {
// 	return &Arena{
// 		MemoryBuffer: NewMemoryBuffer(size),
// 		items:        make(map[int]any),
// 		alignof:      unsafe.Alignof(new(T)),
// 		indirectSize: indirectSize(new(T)),
// 		itemSize:     SizeOfAligned[T](1),
// 	}
// }

// func (a *Arena) NewAt(i int) any {
// 	buf, err := a.Alloc(a.indirectSize, a.alignof)
// 	if err != nil {
// 		return nil
// 	}

// 	a.items[i] = buf

// 	return buf
// }

// func (a *Arena) GetAt(i int) (any, bool) {
// 	buf, ok := a.items[i]

// 	return buf, ok
// }

// func (a *Arena) Length() int {
// 	return len(a.items)
// }
