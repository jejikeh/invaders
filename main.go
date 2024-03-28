package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := NewGame()

	ebiten.SetWindowSize(640, 480)

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
