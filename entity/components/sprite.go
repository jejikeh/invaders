package components

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Image *ebiten.Image

	Transform      *Transform
	DisplayOptions *ebiten.DrawImageOptions
}

func NewSprite(image *ebiten.Image, t *Transform) *Sprite {
	s := &Sprite{
		Image:          image,
		Transform:      t,
		DisplayOptions: &ebiten.DrawImageOptions{},
	}

	s.updateImageDisplayOptions()

	return s
}

func (s *Sprite) Draw(screen *ebiten.Image) {
	s.updateImageDisplayOptions()

	screen.DrawImage(s.Image, s.DisplayOptions)
}

func (s *Sprite) updateImageDisplayOptions() {
	w, h := s.Image.Bounds().Dx(), s.Image.Bounds().Dy()

	s.DisplayOptions.GeoM.Reset()

	s.DisplayOptions.GeoM.Scale(s.Transform.Width, s.Transform.Height)

	s.DisplayOptions.GeoM.Translate(-float64(w)/2*s.Transform.Width, -float64(h)/2*s.Transform.Height)

	s.DisplayOptions.GeoM.Rotate(2 * math.Pi * float64(s.Transform.Angle) / 360)

	s.DisplayOptions.GeoM.Translate(s.Transform.X, s.Transform.Y)
}
