package assets

import (
	"io/fs"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jejikeh/invaders/pkg/engine/log"
)

type Atlas struct {
	Texture     rl.Texture2D
	SpriteCount rl.Vector2
}

type AtlasExport struct {
	Name string
	Path string
	// @Cleanup: Replace rl.Vector2 with custom type
	SpriteCount rl.Vector2

	atlas Atlas
}

func (a Atlas) UVCoords(idx uint32) (pos, size rl.Vector2) {
	size.X = 1.0 / a.SpriteCount.X
	size.Y = 1.0 / a.SpriteCount.Y

	return rl.Vector2{
		X: float32(idx%uint32(a.SpriteCount.X)) * size.X,
		Y: float32(idx/uint32(a.SpriteCount.X)) * size.Y,
	}, size
}

func (s *AtlasExport) Import(fsys fs.ReadFileFS, assets *Assets) error {
	content, err := fsys.ReadFile(s.Path)
	if err != nil {
		return err
	}

	image := rl.LoadImageFromMemory(".png", content, int32(len(content)))

	s.atlas = Atlas{
		Texture:     rl.LoadTextureFromImage(image),
		SpriteCount: s.SpriteCount,
	}

	*assets.Atlases.Get(s.Name) = s.atlas

	log.Engine.Printf("Loaded atlas %s", s.Name)

	return nil
}
