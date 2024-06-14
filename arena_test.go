package gomemory

import (
	"fmt"
	"testing"
)

func TestMallocArenaNewObject(t *testing.T) {
	count := 1000
	arena := NewMallocArena(AproximateSize[uint32](count))
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

func TestMallocArenaFree(t *testing.T) {
	arena := NewMallocArena(1024)
	defer arena.Free()

	type someStruct struct {
		A uint32
		S string
	}

	x := New[someStruct](arena)
	x.A = 1
	x.S = "foo"

	arena.Free()

	y := New[someStruct](arena)
	y.A = 2
	y.S = "bar"

	if x.A != y.A {
		t.Errorf("x.A expected %d got %d", y.A, x.A)
	}

	if x.S != y.S {
		t.Errorf("x.S expected %s got %s", y.S, x.S)
	}
}

func BenchmarkMallocArenaRuntimeNewObject(b *testing.B) {
	type noScanObject struct {
		a byte
		b int
		c uint64
		d complex128
	}

	for _, objectCount := range []int{100, 1000, 10000, 1000000} {
		b.Run(fmt.Sprintf("%d", objectCount), func(b *testing.B) {
			arena := NewMallocArena(AproximateSize[noScanObject](objectCount * b.N))
			defer arena.Free()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				for j := 0; j < objectCount; j++ {
					x := New[noScanObject](arena)
					x.b = j
					x.a = byte(1)
					x.d = complex(float64(j), float64(j))
					x.c = uint64(j)
				}
			}
		})
	}
}

// func BenchmarkRuntimeNewObject(b *testing.B) {
// 	for _, objectCount := range []int{100, 1_000, 10_000, 100_000, 1_000_000} {
// 		b.Run(fmt.Sprintf("%d", objectCount), func(b *testing.B) {
// 			a := newRuntimeAllocator[noScanObject]()
// 			b.ReportAllocs()
// 			for i := 0; i < b.N; i++ {
// 				for j := 0; j < objectCount; j++ {
// 					_ = a.new()
// 				}
// 			}
// 		})
// 	}
// }
