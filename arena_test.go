package goarena

import (
	"testing"
	"unsafe"
)

func TestArena(t *testing.T) {
	arena := NewArena(1024)

	type T struct {
		A uint32
	}

	t1 := New[T](arena)
	assertEqual(t, arena.size, unsafe.Sizeof(t1))

	t2 := New[T](arena)
	t2.A = 1
	assertEqual(t, t1.A, 0)
	assertEqual(t, t2.A, 1)
	assertEqual(t, arena.size, unsafe.Sizeof(t1)+unsafe.Sizeof(t2))

	t3 := New[T](arena)
	t3.A = 2
	assertEqual(t, t1.A, 0)
	assertEqual(t, t2.A, 1)
	assertEqual(t, t3.A, 2)
	assertEqual(t, arena.size, unsafe.Sizeof(t1)+unsafe.Sizeof(t2)+unsafe.Sizeof(t3))

	arena.Free()

	assertEqual(t, arena.size, 0)

	t4 := New[T](arena)
	t4.A = 3

	assertEqual(t, t4.A, 3)
	assertEqual(t, t1.A, 0)
	assertEqual(t, arena.size, unsafe.Sizeof(t4))
}

func assertEqual[T comparable](t *testing.T, a, b T) {
	if a != b {
		t.Errorf("expected %v got %v", b, a)
	}
}
