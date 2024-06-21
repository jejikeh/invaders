package goengine

import (
	"bytes"
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jejikeh/invaders/pkg/goecs"
)

type EbitenSprite struct {
	Image *ebiten.Image
	Path  string
}

func drawEbitenSprites(screen *ebiten.Image, layer *goecs.Layer) {
	entities := layer.Request(goecs.GetComponentID[EbitenSprite](layer), goecs.GetComponentID[Transfrom](layer))

	for _, entity := range entities {
		sprite, _ := goecs.GetComponent[EbitenSprite](layer, entity)
		transform, _ := goecs.GetComponent[Transfrom](layer, entity)

		// @Cleanup: Move this to separate system
		// @Incomplete: Create like asset loader system. It might require you to create like a smartpointers inside gomemory.
		if sprite.Image == nil {
			file, err := os.ReadFile(sprite.Path)
			if err != nil {
				log.Printf("Failed to read %s sprite file: %s", sprite.Path, err)
				continue
			}

			img, _, err := image.Decode(bytes.NewReader(file))
			if err != nil {
				log.Printf("Failed to decode %s sprite file: %s", sprite.Path, err)
				continue
			}

			sprite.Image = ebiten.NewImageFromImage(img)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(transform.Position.X, transform.Position.Y)
		op.GeoM.Scale(transform.Scale.X, transform.Scale.Y)
		op.GeoM.Rotate(transform.Rotation)

		screen.DrawImage(sprite.Image, op)
	}
}
