package main

import "github.com/jejikeh/gomemory"

func main() {
	arena := gomemory.NewArena(1)
	defer arena.Free()

	type T struct {
		A uint32
	}

	_ = gomemory.New[T](arena)

	arena.Print()

	_ = gomemory.New[T](arena)

	arena.Print()
}
