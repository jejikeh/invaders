package goengine

import (
	"fmt"
	"github.com/jejikeh/invaders/pkg/goengine/internal"
)

type Engine struct {
	config  Config
	window  internal.Window
	running bool
}

func NewEngine(c *Config) (*Engine, error) {
	w, err := internal.NewWindow(c.Window)
	if err != nil {
		return nil, err
	}
	
	return &Engine{
		config: *c,
		window: w,
	}, nil
}

func (e *Engine) Running() bool {
	return e.running
}

func (e *Engine) UpdateConfig(update func(*Config)) {
	update(&e.config)
	fmt.Println(e.config.Window.Title)
	e.window.UpdateConfig(e.config.Window)
}
