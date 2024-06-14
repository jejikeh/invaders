package main

import "github.com/jejikeh/goarena"

func main() {
	arena := goarena.NewArena(1)
	defer arena.Free()

	type T struct {
		A uint32
	}

	_ = goarena.New[T](arena)

	arena.Print()

	_ = goarena.New[T](arena)

	arena.Print()
}
