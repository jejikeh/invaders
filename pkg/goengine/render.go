package goengine

import (
	"bytes"
	"image"
	"image/color"
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

func NewScreenImage(engine *Engine, sprite *EbitenSprite) {
	sprite.imageID = engine.spriteBuf.Length()
	ebitenImage := engine.spriteBuf.StoreAt(sprite.imageID, func(i *ebiten.Image) {
		img := ebiten.NewImage(engine.window.Layout(0, 0))
		img.Fill(color.Black)
		*i = *img
	})

	sprite.height = float64(ebitenImage.Bounds().Dy())
	sprite.width = float64(ebitenImage.Bounds().Dx())
	sprite.loaded = true
}

func drawEbitenSprites(engine *Engine, screen *ebiten.Image, layer *goecs.Layer) {
	entities := layer.Request(
		goecs.GetComponentID[EbitenSprite](layer),
		goecs.GetComponentID[Transfrom](layer),
	)

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

		if ebitenImage, imageLoaded := engine.spriteBuf.LoadAt(sprite.imageID); imageLoaded {
			if shader, attachedShader := goecs.GetComponent[Shader](layer, entity); !attachedShader {
				op := &ebiten.DrawImageOptions{}
				applyTransform(&op.GeoM, sprite, transform)
				screen.DrawImage(ebitenImage, op)
			} else {
				if !shader.loaded {
					continue
				}

				ebitenShader, ok := engine.shaders.LoadAt(shader.Name)
				if !ok {
					continue
				}

				op := &ebiten.DrawRectShaderOptions{}
				applyTransform(&op.GeoM, sprite, transform)
				op.Uniforms = engine.shaders.uniforms

				if ebitenImage, ok := engine.spriteBuf.LoadAt(sprite.imageID); ok {
					op.Images[0] = ebitenImage
					screen.DrawRectShader(ebitenImage.Bounds().Dx(), ebitenImage.Bounds().Dy(), ebitenShader, op)
				}
			}
		}
	}
}

func applyTransform(geom *ebiten.GeoM, sprite *EbitenSprite, transform *Transfrom) {
	geom.Translate(-sprite.width/2, -sprite.height/2)
	geom.Rotate(transform.Rotation)
	geom.Translate(sprite.width/2, sprite.height/2)
	geom.Translate(transform.Position.X, transform.Position.Y)
	geom.Scale(transform.Scale.X, transform.Scale.Y)
}
