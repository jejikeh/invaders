package gomemory

import "unsafe"

type Bump interface {
	Mark() unsafe.Pointer
	ResetTo(unsafe.Pointer)
}

var ErrAlreadyBumped = "already bumped"
var ErrNotBumped = "not bumped"

type SingleBump struct {
	MallocArena

	bump bool
}

func NewSingleBump(size int) *SingleBump {
	return &SingleBump{
		MallocArena: *NewMallocArena(size),
	}
}

func (b *SingleBump) Mark() unsafe.Pointer {
	if b.bump {
		panic(ErrAlreadyBumped)
	}
	b.bump = true

	return b.cursor
}

func (b *SingleBump) ResetTo(ptr unsafe.Pointer) {
	if !b.bump {
		panic(ErrNotBumped)
	}

	b.cursor = ptr
}
