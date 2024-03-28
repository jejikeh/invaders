package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jejikeh/invaders/entity"
)

type Game struct {
	Invader *entity.Invader
}

func NewGame() *Game {
	g := &Game{}

	g.Invader = entity.NewInvader(10, 10, 2, 2, 0)

	return g
}

func (g *Game) LoadScene() {

}

func (g *Game) Update() error {
	if err := g.Invader.Update(); err != nil {
		return err
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Invader.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}
