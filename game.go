package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jejikeh/invaders/entity"
)

type GameState int

const (
	Menu GameState = iota
	Play
	Pause
)

type Game struct {
	State   GameState
	Invader *entity.Invader
}

func NewGame() *Game {
	g := &Game{
		State: Menu,
	}

	g.Invader = entity.NewInvader(GameWidth/2, GameHeight/2, 1, 1, 0)

	return g
}

func (g *Game) LoadScene() {

}

func (g *Game) Update() error {
	switch g.State {
	case Play:
		return g.updateGame()

	case Menu:
		return g.updateMenu()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.State {
	case Play:
		g.drawGame(screen)

	case Pause:
		g.drawPause(screen)

	case Menu:
		g.drawMenu(screen)
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, Menu!")
}

func (g *Game) updateMenu() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.State = Play
	}

	return nil
}

func (g *Game) drawGame(screen *ebiten.Image) {
	g.Invader.Draw(screen)
}

func (g *Game) updateGame() error {
	if err := g.Invader.Update(); err != nil {
		return err
	}

	return nil
}

func (g *Game) drawPause(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, Pause!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}
