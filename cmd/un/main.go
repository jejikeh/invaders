package main

import (
	"fmt"
	"unsafe"
)

type T struct {
	x bool
	y []int16
}

const N = unsafe.Offsetof(T{}.y)
const M = unsafe.Sizeof(T{}.y[0])

func main() {
	t := T{y: []int16{123, 456, 789}}
	// "uintptr(p)+N+M+M" is the address of t.y[2].
	ty2 := (*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&t)) + unsafe.Offsetof(t.y) + unsafe.Sizeof(t.y[0])))
	fmt.Println(*ty2) // 789
}
