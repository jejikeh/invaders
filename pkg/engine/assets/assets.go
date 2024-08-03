package assets

import (
	"io/fs"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jejikeh/invaders/pkg/engine/log"
	"github.com/jejikeh/invaders/pkg/gomemory/buf"
)

type Importer interface {
	Import(fsys fs.ReadFileFS, assets *Assets) error
}

type Assets struct {
	muAtlases sync.Mutex
	Atlases   *buf.Map[string, Atlas]

	muShaders sync.Mutex
	Shaders   *buf.Map[string, rl.Shader]
}

func Import(fsys fs.ReadFileFS, importers ...Importer) (*Assets, error) {
	a := &Assets{
		Atlases: buf.NewMap[string, Atlas](len(importers)),
		Shaders: buf.NewMap[string, rl.Shader](len(importers)),
	}

	for _, importer := range importers {
		if err := importer.Import(fsys, a); err != nil {
			log.Engine.Printf("Failed to load asset: %s", err)
		}
	}

	log.Engine.Printf("Loaded %d shaders and %d atlases", a.Shaders.Len(), a.Atlases.Len())

	return a, nil
}

func (a *Assets) Unload() {
	for _, name := range a.Atlases.Keys() {
		log.Engine.Printf("Unloading atlas %s", name)
		rl.UnloadTexture(a.Atlases.Get(name).Texture)
	}

	for _, name := range a.Shaders.Keys() {
		log.Engine.Printf("Unloading shader %s", name)
		rl.UnloadShader(*a.Shaders.Get(name))
	}
}
