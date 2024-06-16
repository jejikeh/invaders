package gomemory

import (
	"fmt"
	"testing"
)

func TestBump(t *testing.T) {
	b := NewSingleBump(1024)
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

func BenchmarkSingleBumpRuntimeNewObject(b *testing.B) {
	type noScanObject struct {
		a byte
		b int
		c uint64
		d complex128
	}

	for _, objectCount := range []int{100, 1000, 10000, 1000000} {
		b.Run(fmt.Sprintf("%d", objectCount), func(b *testing.B) {
			arena := NewSingleBump(AlignedSizeOf[noScanObject](objectCount * b.N))
			defer arena.Free()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				iMark := arena.Mark()
				for j := 0; j < objectCount; j++ {
					x := New[noScanObject](arena)
					x.b = j
					x.a = byte(1)
					x.d = complex(float64(j), float64(j))
					x.c = uint64(j)
				}
				arena.ResetTo(iMark)
			}
		})
	}
}
