package gomemory

import (
	"fmt"
	"bytes"
	"testing"
	"encoding/binary"
)

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

func TestMallocArenaMemoryLayout(t *testing.T) {
	count := 2
	arena := NewMallocArena(AlignedSizeOf[uint32](count))
	defer arena.Free()
	
	x := New[uint32](arena)
	*x = 123
	
	y := New[uint32](arena)
	*y = 456
	
	buf := new(bytes.Buffer)
	bufLen, err := arena.DumpBuffer(buf)
	if err != nil {
		t.Error(err)
	} else if bufLen != AlignedSizeOf[uint32](count) {
		t.Errorf("expected %d written bytes, but got %d", AlignedSizeOf[uint32](count), bufLen)
	}
	
	// @Incomplete: Endians.
	var num uint32
	if err := binary.Read(buf, binary.LittleEndian, &num); err != nil {
		t.Error(err)
	}
	
	if num != *x {
		t.Errorf("expected %d in buffer, but got %d", *x, num)
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
		a byte
		bufLen int
		c uint64
		d complex128
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
