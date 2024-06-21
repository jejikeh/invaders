package arena

import (
	"testing"
	"unsafe"
)

func TestNewMemoryBuffer(t *testing.T) {
	ints := NewMemoryBuffer[int](1025)
	if ints == nil {
		t.Error("ints is nil")
	}

	x, err := ints.New()
	if err != nil {
		t.Error(err)
	}

	if x == nil {
		t.Error("x is nil")

		return
	}

	*x = 123
	x = nil

	atX := ints.At(0)
	if atX == nil {
		t.Error("atX is nil")

		return
	}

	if *atX != 123 {
		t.Errorf("expected %d, but got %d", 123, *atX)
	}
}

func TestMemoryBufferAllocOneObject(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "int",
			data: int(123),
		},
		{
			name: "byte",
			data: byte(123),
		},
		{
			name: "uint32",
			data: uint32(123),
		},
		{
			name: "uint64",
			data: uint64(123),
		},
		{
			name: "float32",
			data: new(float32),
		},
		{
			name: "float64",
			data: new(float64),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			mem := NewMemoryBufferAny(tt.data, 1024)
			v, err := mem.New()

			if err != nil {
				t.Error(err)

				return
			}

			ptr := unsafe.Pointer(v)

			if ptr == nil {
				t.Error("ptr is nil")

				return
			}

			if uintptr(ptr) < uintptr(mem.start) || uintptr(ptr) > uintptr(mem.end) {
				t.Errorf("ptr is not in arena")

				return
			}

			*(*any)(ptr) = tt.data
			if *(*any)(ptr) != tt.data {
				t.Errorf("expected %d, but got %d", tt.data.(int), *(*int)(ptr))

				return
			}
		})
	}
}

// func TestMemoryBufferAllocManyIntObjects(t *testing.T) {
// 	mem := NewMemoryBuffer[int](1025)
// 	var align uintptr = unsafe.Alignof(int(0))

// 	for i := 0; i < 1024; i++ {
// 		v, err := mem.Make()
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		ptr := unsafe.Pointer(v)

// 		if ptr == nil {
// 			t.Error("ptr is nil")
// 		}

// 		if uintptr(ptr)%align != 0 {
// 			t.Errorf("expected %d, but got %d", 0, uintptr(ptr)%align)
// 		}

// 		if uintptr(ptr) < uintptr(mem.start) || uintptr(ptr) > uintptr(mem.end) {
// 			t.Errorf("ptr is not in arena")
// 		}

// 		*(*int)(ptr) = i

// 		if *(*int)(ptr) != i {
// 			t.Errorf("expected %d, but got %d", i, *(*int)(ptr))
// 		}

// 		if mem.size() != uintptr(SizeOfAligned[int](i+1)) {
// 			t.Errorf("expected %d, but got %d", uintptr(SizeOfAligned[int](i+1)), mem.size())
// 		}
// 	}

// 	buf := new(bytes.Buffer)

// 	dump, err := mem.Dump(buf)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if dump != int(mem.size()) {
// 		t.Errorf("expected %d, but got %d", mem.size(), dump)
// 	}

// 	var nums [1024]uint64
// 	if err := binary.Read(buf, binary.LittleEndian, &nums); err != nil {
// 		t.Error(err)
// 	}

// 	for i := 0; i < 1024; i++ {
// 		if nums[1023-i] != uint64(i) {
// 			t.Errorf("expected %d, but got %d", i, nums[i])
// 		}
// 	}
// }
