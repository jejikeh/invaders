package gomemory

import (
	"testing"
)

func TestBumpResetToMark(t *testing.T) {
	t.Parallel()

	b := NewBump(1024)
	defer b.Free()

	x := New[uint32](b)
	*x = 1
	xMark := b.Mark()

	y := New[uint32](b)
	*y = 2

	if *x != 1 {
		t.Errorf("x = %d, want 1", *x)
	}

	if *y != 2 {
		t.Errorf("y = %d, want 2", *y)
	}

	b.ResetTo(xMark)

	z := New[uint32](b)
	*z = 3

	if *x != 1 {
		t.Errorf("x = %d, want 1", *x)
	}

	if *z != 3 {
		t.Errorf("z = %d, want 3", *z)
	}

	if *y != 3 {
		t.Errorf("y = %d, want 3", *y)
	}
}
