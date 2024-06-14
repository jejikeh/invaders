package gomemory

import (
	"testing"
	"unsafe"
)

func TestArenaNew(t *testing.T) {
	intSize := int(unsafe.Sizeof(uint32(0)))
	align := int(unsafe.Alignof(uint32(0)))
	count := 1000
	arena := NewMallocArena((intSize + align) * count)
	defer arena.Free()

	var ints []*uint32

	for i := range count {
		x := New[uint32](arena)
		*x = uint32(i)
		ints = append(ints, x)
	}

	for i := range count {
		v := ints[i]

		if v == nil {
			t.Errorf("expected non-nil value at index %d", i)
		} else if *v != uint32(i) {
			t.Errorf("expected %d got %d", i, *v)
		}
	}
}

func TestArenaGrow(t *testing.T) {
	// 	arena := NewBumpDown(8)

	// 	type Bytes8 struct {
	// 		A [2]uint32
	// 	}

	// 	type Test [8]Bytes8

	// _ = New[Test](arena)
}

func TestAllign(t *testing.T) {
	// type T struct {
	// 	A uint32
	// 	B uint32
	// 	C bool
	// }

	// x := &T{A: 1, B: 2, C: true}

	// xPtrBefore := unsafe.Pointer(&x)

	// y := Align(x)
	// xPtrAfter := unsafe.Pointer(&y)
	// assertEqual(t, uintptr(xPtrBefore), uintptr(xPtrAfter))

}
