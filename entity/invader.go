package entity

import (
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jejikeh/invaders/entity/components"
)

type Invader struct {
	Transform *components.Transform
	Sprite    *components.Sprite
}

func NewInvader(x, y, width, height, angle float64) *Invader {
	imageContent, err := os.Open("resources/alien.png")

	if err != nil {
		panic(err)
	}

	defer imageContent.Close()

	invaderImage, _, err := image.Decode(imageContent)

	if err != nil {
		panic(err)
	}

	t := components.NewTransform(x, y, width, height, angle)
	s := components.NewSprite(ebiten.NewImageFromImage(invaderImage), t)

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
