package goengine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jejikeh/invaders/pkg/goecs"
	"github.com/jejikeh/invaders/pkg/gomemory"
)

type Config struct {
	Window *WindowConfig
}

type Engine struct {
	config Config
	window *EbitenWindow
	ECS    *goecs.Layer

	time int

	spriteBuf *gomemory.UnsafePool[int, ebiten.Image]
	shaders   *Shaders
}

func NewEngine(c *Config) (*Engine, error) {
	engine := &Engine{
		config:    *c,
		spriteBuf: gomemory.NewUnsafePool[int, ebiten.Image](1024),
		shaders:   NewShaders(8),
	}

	var err error

	engine.window, err = NewEbitenWindow(engine.config.Window, engine)
	if err != nil {
		return nil, err
	}

	engine.ECS = goecs.NewLayer()

	engine.ECS.AddSystems(engine.shaders.loadShaders)
	engine.ECS.AddSystems(updateDebugInfo)

	spawnDebugInfo(engine.ECS)

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
	e.time++
	e.ECS.Update()

	return nil
}

// @Incomplete: Abstract away ebiten in engine layer
func (e *Engine) Draw(screen *ebiten.Image) {
	e.shaders.updateShaderUniforms(e)

	drawEbitenSprites(e, screen, e.ECS)
	drawDebugInfo(screen, e.ECS)
}
