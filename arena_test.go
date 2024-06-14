package gomemory

import (
	"testing"
	"unsafe"
)

func TestArenaNew(t *testing.T) {
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

func TestArenaGrow(t *testing.T) {
	arena := NewArena(8)

	type Bytes8 struct {
		A [2]uint32
	}

	type Test [8]Bytes8

	_ = New[Test](arena)
}

func TestAllign(t *testing.T) {
	type T struct {
		A uint32
		B uint32
		C bool
	}

	x := &T{A: 1, B: 2, C: true}

	xPtrBefore := unsafe.Pointer(&x)

	y := Align(x)
	xPtrAfter := unsafe.Pointer(&y)
	assertEqual(t, uintptr(xPtrBefore), uintptr(xPtrAfter))

}

func assertEqual[T comparable](t *testing.T, a, b T) {
	if a != b {
		t.Errorf("expected %v got %v", b, a)
	}
}
