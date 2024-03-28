package components

import (
	"math"
	"os"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Image *ebiten.Image

	Transform      *Transform
	DisplayOptions *ebiten.DrawImageOptions
}

func NewSprite(path string, t *Transform) *Sprite {
	s := &Sprite{
		Image:          loadImageFromPath(path),
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

func loadImageFromPath(path string) *ebiten.Image {
	// @Cleanup: Make something like AssetLoader to load this type of things?
	// Or make just simple file with constants to path of assets.
	// Also, this can be done in .vars file to manage it more nicely.
	
	imageContent, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer imageContent.Close()

	img, _, err := image.Decode(imageContent)

	if err != nil {
		panic(err)
	}
	
	return ebiten.NewImageFromImage(img)
}

func (s *Sprite) updateImageDisplayOptions() {
	w, h := s.Image.Bounds().Dx(), s.Image.Bounds().Dy()

	s.DisplayOptions.GeoM.Reset()

	s.DisplayOptions.GeoM.Scale(s.Transform.Width, s.Transform.Height)

	s.DisplayOptions.GeoM.Translate(-float64(w)/2*s.Transform.Width, -float64(h)/2*s.Transform.Height)

	s.DisplayOptions.GeoM.Rotate(2 * math.Pi * float64(s.Transform.Angle) / 360)

	s.DisplayOptions.GeoM.Translate(s.Transform.X, s.Transform.Y)
}
