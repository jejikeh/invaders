package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jejikeh/invaders/entity/components"
)

type Invader struct {
	Transform *components.Transform
	Sprite    *components.Sprite
}

func NewInvader(x, y, width, height, angle float64) *Invader {
	t := components.NewTransform(x, y, width, height, angle)
	s := components.NewSprite("resources/alien.png", t)

	return &Invader{
		Transform: t,
		Sprite:    s,
	}
}

func (i *Invader) Draw(screen *ebiten.Image) {
	i.Sprite.Draw(screen)
}

func (i *Invader) Update() error {
	return nil
}
