package goengine

import "github.com/jejikeh/invaders/pkg/goengine/internal"

type Config struct {
	Window *WindowConfig
}

type WindowConfig = internal.WindowConfig
type WindowProvider = internal.WindowProvider

const (
	Ebiten WindowProvider = iota
)
