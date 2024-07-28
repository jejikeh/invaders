package main

import (
	"embed"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets/*
var assets embed.FS

type AssetPaths map[string]string

type Assets struct {
	Textures map[string]rl.Texture2D
	Shaders  map[string]rl.Shader
}

func NewAssets(assetsMap AssetPaths) (*Assets, error) {
	a := &Assets{
		Textures: make(map[string]rl.Texture2D, len(assetsMap)),
	}

	for name, path := range assetsMap {
		content, err := assets.ReadFile(path)
		if err != nil {
			a.Unload()

			return nil, err
		}

		switch filepath.Ext(path) {
		case ".png":
			a.loadTexture(content, name)
		case ".fs":
			// @Cleanup: Load fragment and vertex files as one shader
			a.loadShader(nil, content, name)
		case ".vs":
			// @Cleanup: Load fragment and vertex files as one shader
			a.loadShader(content, nil, name)
		}
	}

	a.Textures["core/circle"] = rl.LoadTextureFromImage(rl.GenImageGradientRadial(16, 16, 0.3, rl.White, rl.Black))

	return a, nil
}

func (a *Assets) Unload() {
	for name, texture := range a.Textures {
		logger.engine.Printf("Unloading texture %s", name)
		rl.UnloadTexture(texture)
	}

	for name, shader := range a.Shaders {
		logger.engine.Printf("Unloading shader %s", name)
		rl.UnloadShader(shader)
	}
}

func (a *Assets) loadTexture(content []byte, name string) {
	logger.engine.Printf("Loading texture %s", name)

	image := rl.LoadImageFromMemory(".png", content, int32(len(content)))
	texture := rl.LoadTextureFromImage(image)

	a.Textures[name] = texture
}

func (a *Assets) loadShader(vertex, fragment []byte, name string) {
	logger.engine.Printf("Loading shader %s", name)

	shader := rl.LoadShaderFromMemory(string(vertex), string(fragment))
	a.Shaders[name] = shader
}
