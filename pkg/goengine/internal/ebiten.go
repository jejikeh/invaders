package internal

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jejikeh/invaders/pkg/gomath"
)

type EbitenWindow struct {
	config WindowConfig
	layout *gomath.Vec2
}

func NewEbitenWindow(config *WindowConfig) (*EbitenWindow, error) {
	window := &EbitenWindow{
		config: *config,
	}

	window.UpdateConfig(&window.config)

	return window, window.open()
}

func (e *EbitenWindow) open() error {	
	return ebiten.RunGame(e)
}

func (e *EbitenWindow) UpdateConfig(newConfig *WindowConfig) {
	if &e.config != newConfig {
		e.config = *newConfig
	}
	
	e.layout = e.config.Size.DivFloat(e.config.Scale)
	
	ebiten.SetWindowSize(int(e.config.Size.X), int(e.config.Size.Y))
	ebiten.SetWindowTitle(e.config.Title)
}

func (e *EbitenWindow) Update() error {
	return nil
}

func (e *EbitenWindow) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (e *EbitenWindow) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(e.layout.X), int(e.layout.Y)
}

func (e *EbitenWindow) LayoutF(outsideWidth, outsideHeight float64) (float64, float64) {
	return e.layout.X, e.layout.Y
}
