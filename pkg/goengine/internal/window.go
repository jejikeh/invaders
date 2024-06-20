package internal

import (
	"fmt"
	"github.com/jejikeh/invaders/pkg/gomath"
)

type WindowProvider int

const (
	Ebiten WindowProvider = iota
)

type WindowConfig struct {
	Title    string
	Size     *gomath.Vec2
	Scale 	 float64
	Provider WindowProvider
}

type Window interface {
	UpdateConfig(*WindowConfig)
}

func NewWindow(c *WindowConfig) (Window, error) {
	switch c.Provider {
	case Ebiten:
		return NewEbitenWindow(c)
	default:
		return nil, fmt.Errorf("unknow window provider %v", c.Provider)
	}
}
