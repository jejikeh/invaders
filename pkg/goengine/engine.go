package goengine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jejikeh/invaders/pkg/goecs"
)

type Config struct {
	Window *WindowConfig
}

type Engine struct {
	config Config
	window *EbitenWindow
	ECS    *goecs.Layer
}

func NewEngine(c *Config) (*Engine, error) {
	engine := &Engine{
		config: *c,
	}

	var err error

	engine.window, err = NewEbitenWindow(engine.config.Window, engine)
	if err != nil {
		return nil, err
	}

	engine.ECS = goecs.NewLayer()

	return engine, err
}

func (e *Engine) UpdateConfig(update func(*Config)) {
	update(&e.config)

	e.window.UpdateConfig(e.config.Window)
}

func (e *Engine) Run() error {
	return e.window.Open()
}

func (e *Engine) Update() error {
	e.ECS.Update()

	return nil
}

// @Incomplete: Abstract away ebiten in engine layer
func (e *Engine) Draw(screen *ebiten.Image) {
	drawEbitenSprites(screen, e.ECS)
}
