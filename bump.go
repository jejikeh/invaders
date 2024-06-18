package gomemory

import (
	"unsafe"
)

var ErrAlreadyMarked = "already marked"
var ErrNotMarked = "not marked"

type Bump struct {
	*MallocArena
}

func NewBump(size int) *Bump {
	return &Bump{
		MallocArena: NewMallocArena(size),
	}
}

func (b *Bump) Free() {
	b.MallocArena.Free()
}

func (b *Bump) Mark() unsafe.Pointer {
	return b.cursor
}

func (b *Bump) ResetTo(ptr unsafe.Pointer) {
	b.cursor = ptr
}
