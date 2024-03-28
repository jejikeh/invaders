package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WindowWidth = 640
	WindowHeight = 480
	
	GameWidth = 320
	GameHeight = 240
)

func main() {
	g := NewGame()

	ebiten.SetWindowSize(WindowWidth, WindowHeight)

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
