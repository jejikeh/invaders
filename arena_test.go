package gomemory

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"
)

// @Cleanup: Make sure to name properly tests, also if`s looks kinda ugly a the moment.

func TestMallocArenaNewObject(t *testing.T) {
	count := 1000
	arena := NewMallocArena(AlignedSizeOf[uint32](count))
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

func TestAlignedSizeOfType(t *testing.T) {
	testAlignedSizeTimes[bool](t, 1000)
	testAlignedSizeTimes[int8](t, 1000)
	testAlignedSizeTimes[uint8](t, 1000)
	testAlignedSizeTimes[int16](t, 1000)
	testAlignedSizeTimes[uint16](t, 1000)
	testAlignedSizeTimes[int32](t, 1000)
	testAlignedSizeTimes[uint32](t, 1000)
	testAlignedSizeTimes[int64](t, 1000)
	testAlignedSizeTimes[uint64](t, 1000)
	testAlignedSizeTimes[int](t, 1000)
	testAlignedSizeTimes[uint](t, 1000)
	testAlignedSizeTimes[uintptr](t, 1000)
	testAlignedSizeTimes[float32](t, 1000)
	testAlignedSizeTimes[float64](t, 1000)
	testAlignedSizeTimes[complex64](t, 1000)
	testAlignedSizeTimes[complex128](t, 1000)
	testAlignedSizeTimes[string](t, 1000)

	testAlignedSizeTimes[struct {
		A int
		B uintptr
		C struct {
			CA string
			CB func(string) int
			CC float32
		}
		D complex64
		E unsafe.Pointer
	}](t, 1000)
}

func testAlignedSizeTimes[T any](t *testing.T, count int) {
	t.Helper()

	for n := 1; n <= count; n++ {
		alignedSize := AlignedSizeOf[T](n)
		arena := NewMallocArena(AlignedSizeOf[T](n))

		for range n {
			_ = New[T](arena)
		}

		if alignedSize != int(arena.size()) {
			t.Errorf("calculated aligned size is %d, but arena size %d for %d %T`s", alignedSize, int(arena.size()), n, *new(T))
		}

		buf := new(bytes.Buffer)
		if bufLen, err := arena.WriteRawMemory(buf); err != nil {
			t.Error(err)
		} else if bufLen != alignedSize {
			t.Errorf("calculated aligned size is %d, but WriteRawMemory size %d for %d %T`s", alignedSize, bufLen, n, *new(T))
		}

		if arena.Free(); arena.size() != 0 {
			t.Errorf("arena size is not 0 after Free")
		}
		
		buf = new(bytes.Buffer)
		if bufLen, err := arena.WriteRawMemory(buf); err != nil {
			t.Error(err)
		} else if bufLen != 0 {
			t.Errorf("Dump buffer expected to be 0, but got %d", bufLen)
		}
	}
}

func TestMallocArenaMemoryLayout(t *testing.T) {
	arena := NewMallocArena(AlignedSizeOf[uint32](2))
	defer arena.Free()

	x := New[uint32](arena)
	*x = 1

	y := New[uint32](arena)
	*y = 2

	buf := new(bytes.Buffer)
	bufLen, err := arena.WriteRawMemory(buf)
	if err != nil {
		t.Error(err)
	} else if bufLen != AlignedSizeOf[uint32](2) {
		t.Errorf("expected %d dumped bytes, but arena reported size is %d", AlignedSizeOf[uint32](2), bufLen)
	}

	// @Incomplete: Endians.
	var num [4]uint32
	if err := binary.Read(buf, binary.LittleEndian, &num); err != nil {
		t.Error(err)
	}

	if num[2] != *x {
		t.Errorf("expected %d in buffer, but got %d", *x, num[2])
	}
	
	if num[0] != *y {
		t.Errorf("expected %d in buffer, but got %d", *y, num[0])
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

func BenchmarkMallocArenaRuntimeNewObject(bufLen *testing.B) {
	type noScanObject struct {
		a      byte
		bufLen int
		c      uint64
		d      complex128
	}

	for _, objectCount := range []int{100, 1000, 10000, 1000000} {
		bufLen.Run(fmt.Sprintf("%d", objectCount), func(bufLen *testing.B) {
			arena := NewMallocArena(AlignedSizeOf[noScanObject](objectCount * bufLen.N))
			defer arena.Free()
			bufLen.ReportAllocs()
			for i := 0; i < bufLen.N; i++ {
				for j := 0; j < objectCount; j++ {
					x := New[noScanObject](arena)
					x.bufLen = j
					x.a = byte(1)
					x.d = complex(float64(j), float64(j))
					x.c = uint64(j)
				}
			}
		})
	}
}

// func BenchmarkRuntimeNewObject(bufLen *testing.B) {
// 	for _, objectCount := range []int{100, 1_000, 10_000, 100_000, 1_000_000} {
// 		bufLen.Run(fmt.Sprintf("%d", objectCount), func(bufLen *testing.B) {
// 			a := newRuntimeAllocator[noScanObject]()
// 			bufLen.ReportAllocs()
// 			for i := 0; i < bufLen.N; i++ {
// 				for j := 0; j < objectCount; j++ {
// 					_ = a.new()
// 				}
// 			}
// 		})
// 	}
// }
