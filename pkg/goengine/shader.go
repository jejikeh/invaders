package goengine

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jejikeh/invaders/pkg/goecs"
	"github.com/jejikeh/invaders/pkg/gomemory"
)

// @Incomplete: This is also can be in separate system.
// Maybe, we can separate asset loading and declare loading system for each asset differently?

type Shaders struct {
	*gomemory.UnsafePool[string, ebiten.Shader]
	idx      map[string]int
	uniforms map[string]any
}

func NewShaders(count int) *Shaders {
	return &Shaders{
		UnsafePool: gomemory.NewUnsafePool[string, ebiten.Shader](count),
		idx:        make(map[string]int),
		uniforms: map[string]any{
			"Time": 0.0,
		},
	}
}

type Shader struct {
	Name   string // @Cleanup: Do i need this? I can get name of the file...
	Path   string
	loaded bool
	// @Incomplete: Add more uniforms?
}

func (s *Shaders) loadShaders(layer *goecs.Layer) {
	for _, entity := range layer.Request(goecs.GetComponentID[Shader](layer)) {
		shader, _ := goecs.GetComponent[Shader](layer, entity)

		if shader.loaded {
			continue
		}

		s.StoreAt(shader.Name, func(s *ebiten.Shader) {
			file, err := os.ReadFile(shader.Path)
			if err != nil {
				panic(err)
			}

			shader, err := ebiten.NewShader(file)
			if err != nil {
				panic(err)
			}

			*s = *shader
		})

		shader.loaded = true
	}
}

func (s *Shaders) updateShaderUniforms(engine *Engine) {
	s.uniforms["Time"] = float32(engine.time / 60)
}
