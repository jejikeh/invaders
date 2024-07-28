package assets

import (
	"fmt"
	"io/fs"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jejikeh/invaders/pkg/engine/log"
	"github.com/jejikeh/invaders/pkg/gomemory/buf"
)

type Importer interface {
	Import(fsys fs.ReadFileFS) error
}

type ShaderExport struct {
	Name         string
	VertexPath   string
	FragmentPath string

	shader rl.Shader
}

func (s *ShaderExport) Import(fsys fs.ReadFileFS) error {
	if filepath.Ext(s.FragmentPath) != ".fs" && s.FragmentPath != "" {
		return fmt.Errorf("invalid fragment path, should be .fs: %s", s.FragmentPath)
	} else if filepath.Ext(s.VertexPath) != ".vs" && s.VertexPath != "" {
		return fmt.Errorf("invalid vertex path, should be .vs: %s", s.VertexPath)
	}

	var vertex, fragment []byte

	if s.VertexPath != "" {
		var err error
		vertex, err = fsys.ReadFile(s.VertexPath)

		if err != nil {
			return err
		}
	}

	if s.FragmentPath != "" {
		var err error
		fragment, err = fsys.ReadFile(s.FragmentPath)

		if err != nil {
			return err
		}
	}

	s.shader = rl.LoadShaderFromMemory(string(vertex), string(fragment))

	log.Engine.Printf("Loaded shader %s", s.Name)

	return nil
}

type AtlasExport struct {
	Name string
	Path string
	// @Cleanup: Replace rl.Vector2 with our type
	SpriteCount rl.Vector2

	atlas Atlas
}

func (s *AtlasExport) Import(fsys fs.ReadFileFS) error {
	content, err := fsys.ReadFile(s.Path)
	if err != nil {
		return err
	}

	image := rl.LoadImageFromMemory(".png", content, int32(len(content)))

	s.atlas = Atlas{
		Texture:     rl.LoadTextureFromImage(image),
		SpriteCount: s.SpriteCount,
	}

	log.Engine.Printf("Loaded atlas %s", s.Name)

	return nil
}

type AssetPaths map[string]string

type Assets struct {
	Atlases *buf.Map[string, Atlas]
	Shaders *buf.Map[string, rl.Shader]

	Unknown []Importer
}

func Import(fsys fs.ReadFileFS, importers ...Importer) (*Assets, error) {
	a := &Assets{
		Atlases: buf.NewMap[string, Atlas](len(importers)),
		Shaders: buf.NewMap[string, rl.Shader](len(importers)),
	}

	for _, importer := range importers {
		if err := importer.Import(fsys); err != nil {
			log.Engine.Printf("Failed to load asset: %s", err)

			return nil, err
		}

		switch importer.(type) {
		case *ShaderExport:
			shader, _ := importer.(*ShaderExport)
			*a.Shaders.Get(shader.Name) = shader.shader
		case *AtlasExport:
			texture, _ := importer.(*AtlasExport)
			*a.Atlases.Get(texture.Name) = texture.atlas
		default:
			a.Unknown = append(a.Unknown, importer)
		}
	}

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

type Atlas struct {
	Texture     rl.Texture2D
	SpriteCount rl.Vector2
}

func (a Atlas) UVCoords(idx uint32) (pos, size rl.Vector2) {
	size.X = 1.0 / a.SpriteCount.X
	size.Y = 1.0 / a.SpriteCount.Y

	return rl.Vector2{
		X: float32(idx%uint32(a.SpriteCount.X)) * size.X,
		Y: float32(idx/uint32(a.SpriteCount.X)) * size.Y,
	}, size
}
