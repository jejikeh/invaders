package goengine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jejikeh/invaders/pkg/gomath"
)

type WindowConfig struct {
	Title string
	Size  gomath.Vec2
	Scale float64
}

type EbitenWindow struct {
	config  WindowConfig
	layout  gomath.Vec2
	handler *Engine
}

func NewEbitenWindow(config *WindowConfig, engine *Engine) (*EbitenWindow, error) {
	window := &EbitenWindow{
		config:  *config,
		handler: engine,
	}

	window.UpdateConfig(&window.config)

	return window, nil
}

// @Cleanup.
func (e *EbitenWindow) Open() error {
	return ebiten.RunGame(e)
}

func (e *EbitenWindow) UpdateConfig(newConfig *WindowConfig) {
	if &e.config != newConfig {
		e.config = *newConfig
	}

	e.layout = gomath.NewVec2FromVec(e.config.Size).Scale(e.config.Scale)

	ebiten.SetWindowSize(int(e.config.Size.X), int(e.config.Size.Y))
	ebiten.SetWindowTitle(e.config.Title)
}

func (e *EbitenWindow) Update() error {
	return e.handler.Update()
}

func (e *EbitenWindow) Draw(screen *ebiten.Image) {
	e.handler.Draw(screen)
}

func (e *EbitenWindow) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(e.layout.X), int(e.layout.Y)
}

func (e *EbitenWindow) LayoutF(outsideWidth, outsideHeight float64) (float64, float64) {
	return e.layout.X, e.layout.Y
}
