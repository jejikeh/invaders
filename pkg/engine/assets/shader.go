package assets

import (
	"fmt"
	"io/fs"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jejikeh/invaders/pkg/engine/log"
)

type ShaderExport struct {
	Name         string
	VertexPath   string
	FragmentPath string

	shader rl.Shader
}

func (s *ShaderExport) Import(fsys fs.ReadFileFS, assets *Assets) error {
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

	*assets.Shaders.Get(s.Name) = s.shader

	log.Engine.Printf("Loaded shader %s", s.Name)

	return nil
}
