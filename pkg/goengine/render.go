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
	Path    string
	loaded  bool
	imageID int
	width   float64
	height  float64
}

func drawEbitenSprites(engine *Engine, screen *ebiten.Image, layer *goecs.Layer) {
	entities := layer.Request(goecs.GetComponentID[EbitenSprite](layer), goecs.GetComponentID[Transfrom](layer))

	for _, entity := range entities {
		sprite, _ := goecs.GetComponent[EbitenSprite](layer, entity)
		transform, _ := goecs.GetComponent[Transfrom](layer, entity)

		// @Cleanup: Move this to separate system
		// @Incomplete: Create like asset loader system. It might require you to create like a smartpointers inside gomemory.
		if !sprite.loaded {
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

			sprite.imageID = engine.spriteBuf.Length()
			ebitenImage := engine.spriteBuf.StoreAt(sprite.imageID, func(i *ebiten.Image) {
				*i = *ebiten.NewImageFromImage(img)
			})

			sprite.height = float64(ebitenImage.Bounds().Dy())
			sprite.width = float64(ebitenImage.Bounds().Dx())

			sprite.loaded = true
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-sprite.width/2, -sprite.height/2)
		op.GeoM.Rotate(transform.Rotation)
		op.GeoM.Translate(sprite.width/2, sprite.height/2)
		op.GeoM.Translate(transform.Position.X, transform.Position.Y)
		op.GeoM.Scale(transform.Scale.X, transform.Scale.Y)

		if ebitenImage, ok := engine.spriteBuf.LoadAt(sprite.imageID); ok {
			screen.DrawImage(ebitenImage, op)
		}
	}
}
