package entity

import "github.com/hajimehoshi/ebiten/v2"

type EntityInterface interface {
	Draw(screen *ebiten.Image)
	Update() error
}

type EntityType int

const (
	InvaderType EntityType = 1 << iota
)
