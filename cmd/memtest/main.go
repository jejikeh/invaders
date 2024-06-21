package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/jejikeh/invaders/pkg/gomemory"
	"github.com/ortuman/nuke"
)

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

func fibonacci() func() int {
	a := 0
	b := 1
	return func() int {
		a, b = b, a+b
		return b - a
	}
}

func main() {
	arena()

	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Print(f())
	}
}

func arena() {
	arena := nuke.NewMonotonicArena(256*1024, 80)

	// Allocate a new object of type Foo.
	// fooRef := nuke.New[Foo](arena)
	// arena := gomemory.NewMallocArena(gomemory.SizeOfAligned[Test](10))
	// test := gomemory.New[Test](arena)
	test := gomemory.New[Test](arena)
	test.A = 1
	test.B = 2
	test.C = &C{
		A: 3,
		B: []*B{{A: 4}, {A: 5}},
	}

	// b := &test.C.B

	// test := &Test{A: 1, B: 2, C: &C{A: 3, B: []*B{{A: 4}, {A: 5}}}}

	ptr := unsafe.Pointer(test)

	runtime.GC()

	testCB1 := *(**C)(unsafe.Add(ptr, unsafe.Offsetof(test.C)))
	if testCB1.A != 3 {
		panic("testCB1.A != 3")
	}

	if len(testCB1.B) != 2 {
		panic("testCB1.B != b")
	}

	// if &testCB1.B != b {
	// panic("testCB1.B != b")
	// }

	if testCB1.B[0].A != 4 {
		panic("testCB1.B[0].A != 4")
	}

	if testCB1.B[1].A != 5 {
		panic("testCB1.B[1].A != 5")
	}
}
