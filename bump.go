package gomemory

import (
	"unsafe"
)

type Bump interface {
	Mark() unsafe.Pointer
	ResetTo(unsafe.Pointer)
}

var ErrAlreadyMarked = "already marked"
var ErrNotMarked = "not marked"

// @Incomplete: Remove SingleBump, make it normal bump
type SingleBump struct {
	*MallocArena

	mark bool
}

func NewSingleBump(size int) *SingleBump {
	return &SingleBump{
		MallocArena: NewMallocArena(size),
	}
}

func (b *SingleBump) Free() {
	b.mark = false
	b.MallocArena.Free()
}

func (b *SingleBump) Mark() unsafe.Pointer {
	if b.mark {
		panic(ErrAlreadyMarked)
	}
	b.mark = true

	return b.cursor
}

func (b *SingleBump) ResetTo(ptr unsafe.Pointer) {
	if !b.mark {
		panic(ErrNotMarked)
	}

	b.cursor = ptr
	b.mark = false
}
