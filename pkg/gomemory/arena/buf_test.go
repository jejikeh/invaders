package arena

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"runtime"
	"runtime/debug"
	"testing"
	"unsafe"
)

func TestBufNew(t *testing.T) {
	debug.SetGCPercent(1)

	tests := []struct {
		name  string
		count int
		data  any
	}{
		{
			name:  "single int",
			count: 1,
			data:  int(1),
		},
		{
			name:  "two int array",
			count: 2,
			data:  [2]int{1, 2},
		},
		{
			name:  "two int slice",
			count: 2,
			data:  []int{1, 2},
		},
		{
			name:  "two int struct",
			count: 2,
			data: struct {
				x int
				y int
			}{1, 2},
		},
		{
			name:  "two int struct pointer",
			count: 2,
			data: &struct {
				x int
				y int
			}{1, 2},
		},
		{
			name:  "two int struct pointer array",
			count: 2,
			data: [2]*struct {
				x int
				y int
			}{{1, 2}, {3, 4}},
		},
		{
			name:  "two int struct pointer slice",
			count: 2,
			data: []*struct {
				x int
				y int
			}{{1, 2}, {3, 4}},
		},
		{
			name:  "two pointer stuct",
			count: 2,
			data: struct {
				x *int
				y *int
			}{new(int), new(int)},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			Buf := NewBuf(tt.count, tt.data)

			if len(Buf.refs) != tt.count {
				t.Errorf("len(Buf.refs) = %d, want %d", len(Buf.refs), tt.count)
			}

			d := Buf.Store(func(a *any) {
				*a = tt.data
			})

			if !reflect.DeepEqual(tt.data, *d) {
				t.Errorf("data = %v, want %v", *d, tt.data)
			}
		})
	}
}

func TestBufStorePointerOutsideBuf(t *testing.T) {
	debug.SetGCPercent(1)

	tests := []struct {
		name     string
		setup    func() *Buf[any]
		store    func(*any)
		validate func(*testing.T, *any)
	}{
		{
			name: "store pointer outside of buffer",
			setup: func() *Buf[any] {
				return NewBuf[any](1, struct {
					x *int
					y *int
				}{new(int), new(int)})
			},
			store: func(a *any) {
				s := struct {
					x *int
					y *int
				}{
					new(int),
					new(int),
				}

				*s.x = 2
				*s.y = 3

				*a = s
			},
			validate: func(t *testing.T, a *any) {
				t.Helper()

				s := (*a).(struct {
					x *int
					y *int
				})

				if *s.x != 2 || *s.y != 3 {
					t.Errorf("s.x = %d, s.y = %d, want s.x = 2, s.y = 3", *s.x, *s.y)
				}
			},
		},
		{
			name: "store pointer array outside of buffer",
			setup: func() *Buf[any] {
				return NewBuf[any](1, [2]*struct {
					x *int
					y *int
				}{{new(int), new(int)}, {new(int), new(int)}})
			},
			store: func(a *any) {
				s := [2]*struct {
					x *int
					y *int
				}{
					{new(int), new(int)},
					{new(int), new(int)},
				}

				*s[0].x = 2
				*s[0].y = 3
				*s[1].x = 4
				*s[1].y = 5

				*a = s
			},
			validate: func(t *testing.T, a *any) {
				t.Helper()

				s := (*a).([2]*struct {
					x *int
					y *int
				})

				if *s[0].x != 2 || *s[0].y != 3 || *s[1].x != 4 || *s[1].y != 5 {
					t.Errorf("s[0].x = %d, s[0].y = %d, s[1].x = %d, s[1].y = %d, want s[0].x = 2, s[0].y = 3, s[1].x = 4, s[1].y = 5", *s[0].x, *s[0].y, *s[1].x, *s[1].y)
				}
			},
		},
		{
			name: "store struct with array outside of buffer",
			setup: func() *Buf[any] {
				return NewBuf[any](1, struct {
					x *int
					y *int
					z [2]*struct {
						a *int
						b *int
					}
				}{new(int), new(int), [2]*struct {
					a *int
					b *int
				}{{new(int), new(int)}, {new(int), new(int)}}})
			},
			store: func(a *any) {
				s := struct {
					x *int
					y *int
					z [2]*struct {
						a *int
						b *int
					}
				}{
					new(int),
					new(int),
					[2]*struct {
						a *int
						b *int
					}{{new(int), new(int)}, {new(int), new(int)}},
				}

				*s.x = 2
				*s.y = 3
				*s.z[0].a = 4
				*s.z[0].b = 5
				*s.z[1].a = 6
				*s.z[1].b = 7

				*a = s
			},
			validate: func(t *testing.T, a *any) {
				t.Helper()

				s := (*a).(struct {
					x *int
					y *int
					z [2]*struct {
						a *int
						b *int
					}
				})

				if *s.x != 2 || *s.y != 3 || *s.z[0].a != 4 || *s.z[0].b != 5 || *s.z[1].a != 6 || *s.z[1].b != 7 {
					t.Errorf("s.x = %d, s.y = %d, s.z[0].a = %d, s.z[0].b = %d, *s.z[1].a = %d, *s.z[1].b = %d, want s.x = 2, s.y = 3, s.z[0].a = 4, s.z[0].b = 5, s.z[1].a = 6, s.z[1].b = 7", *s.x, *s.y, *s.z[0].a, *s.z[0].b, *s.z[1].a, *s.z[1].b)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			Buf := tt.setup()
			obj := Buf.Store(tt.store)
			tt.validate(t, obj)
		})
	}
}

func TestBufPointerOutsideBuf(t *testing.T) {
	t.Parallel()

	debug.SetGCPercent(1)

	type B struct {
		A int
	}

	type C struct {
		A int
		B []*B
	}

	type Test struct {
		A int
		B int
		C *C
	}

	Buf := NewBuf(1, Test{})
	test := Buf.Store(func(t *Test) {
		*t = Test{A: 1, B: 2, C: &C{A: 3, B: []*B{{A: 4}, {A: 5}}}}
	})

	ptr := unsafe.Pointer(test)

	runtime.GC()

	testCB1 := *(**C)(unsafe.Add(ptr, unsafe.Offsetof(test.C)))
	if testCB1.A != 3 {
		t.Errorf("test.C.B.A = %d, want 3", testCB1.A)
	}

	if len(testCB1.B) != 2 {
		t.Errorf("len(test.C.B) = %d, want 2", len(testCB1.B))
	}

	if testCB1.B[0].A != 4 {
		t.Errorf("test.C.B[0].A = %d, want 4", testCB1.B[0].A)
	}

	if testCB1.B[1].A != 5 {
		t.Errorf("test.C.B[1].A = %d, want 5", testCB1.B[1].A)
	}
}

func TestBufWrite(t *testing.T) {
	t.Parallel()

	Buf := NewBuf(2, uint32(0))

	Buf.Store(func(i *uint32) {
		*i = 1
	})

	Buf.Store(func(i *uint32) {
		*i = 2
	})

	buffer := new(bytes.Buffer)
	n, err := Buf.write(buffer)
	if err != nil {
		t.Fatalf("failed to write buffer: %v", err)
	}

	if n != 8 {
		t.Errorf("n = %d, want 8", n)
	}

	// @Cleanup: The alignment is still not right.

	var nums [2]uint32
	if err := binary.Read(buffer, binary.LittleEndian, &nums); err != nil {
		t.Fatalf("failed to read buffer: %v", err)
	}

	if nums[0] != 1 || nums[1] != 2 {
		t.Errorf("nums = %v, want [1, 2]", nums)
	}
}
