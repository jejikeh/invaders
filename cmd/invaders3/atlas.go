package main

import rl "github.com/gen2brain/raylib-go/raylib"

type RegionAtlas struct {
	texture    rl.Texture2D
	dimensions rl.Vector2
}

func NewAtlas(texture rl.Texture2D, dimensions rl.Vector2) RegionAtlas {
	return RegionAtlas{
		texture:    texture,
		dimensions: dimensions,
	}
}

func (a RegionAtlas) getPositionAt(idx uint32) (pos, size rl.Vector2) {
	size.X = 1.0 / a.dimensions.X
	size.Y = 1.0 / a.dimensions.Y

	return rl.Vector2{
		X: float32(idx%uint32(a.dimensions.X)) * size.X,
		Y: float32(idx/uint32(a.dimensions.X)) * size.Y,
	}, size
}
