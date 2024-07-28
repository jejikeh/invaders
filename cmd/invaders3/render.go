package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RenderTexture struct {
	rl.RenderTexture2D

	atlas *Atlas
}

func NewRenderTexture(width, height int32, atlas *Atlas) RenderTexture {
	return RenderTexture{
		RenderTexture2D: rl.LoadRenderTexture(width, height),
		atlas:           atlas,
	}
}

func (r RenderTexture) Unload() {
	rl.UnloadRenderTexture(r.RenderTexture2D)
}

func (r *RenderTexture) Draw() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.Black)
	rl.DrawTexturePro(
		r.Texture,
		rl.NewRectangle(0, 0, float32(r.Texture.Width), float32(r.Texture.Height)),
		rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)
}

func (r RenderTexture) BlendTextures(srcs ...RenderTexture) {
	rl.BeginBlendMode(rl.BlendAlpha)
	defer rl.EndBlendMode()

	rl.BeginTextureMode(r.RenderTexture2D)
	defer rl.EndTextureMode()

	rl.ClearBackground(rl.Black)

	for _, src := range srcs {
		rl.DrawTexturePro(
			src.Texture,
			rl.NewRectangle(0, 0, float32(src.Texture.Width), float32(src.Texture.Height)),
			rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
			rl.Vector2{X: 0, Y: 0},
			0,
			rl.White,
		)
	}
}

func (r RenderTexture) Func2D(camera rl.Camera3D, render func()) {
	rl.BeginTextureMode(r.RenderTexture2D)
	defer rl.EndTextureMode()

	rl.ClearBackground(rl.Color{})

	render()
}

func (r RenderTexture) Func3D(camera rl.Camera3D, render func()) {
	rl.BeginTextureMode(r.RenderTexture2D)
	defer rl.EndTextureMode()

	rl.BeginMode3D(camera)
	defer rl.EndMode3D()

	rl.Begin(rl.Triangles)
	defer rl.End()

	rl.ClearBackground(rl.Color{})

	render()
}

// @Incomplete: Rotations are completely broken.
func (r RenderTexture) QuadFunc(position, scale rl.Vector3, idx uint32, modifyVertex ...func(rl.Vector3) rl.Vector3) {
	rl.SetTexture(r.atlas.texture.ID)

	rl.Color4ub(255, 255, 255, 255)

	a := rl.Vector3{X: position.X - scale.X/2, Y: position.Y, Z: position.Z - scale.Z/2}
	b := rl.Vector3{X: position.X - scale.X/2, Y: position.Y, Z: position.Z + scale.Z/2}
	c := rl.Vector3{X: position.X + scale.X/2, Y: position.Y, Z: position.Z + scale.Z/2}
	d := rl.Vector3{X: position.X + scale.X/2, Y: position.Y, Z: position.Z - scale.Z/2}

	for _, f := range modifyVertex {
		a = f(a)
		b = f(b)
		c = f(c)
		d = f(d)
	}

	pos, size := r.atlas.getPositionAt(idx)

	rl.TexCoord2f(pos.X, pos.Y)
	rl.Vertex3f(a.X, a.Y, a.Z)

	rl.TexCoord2f(pos.X, pos.Y+size.Y)
	rl.Vertex3f(b.X, b.Y, b.Z)

	rl.TexCoord2f(pos.X+size.X, pos.Y+size.Y)
	rl.Vertex3f(c.X, c.Y, c.Z)

	rl.TexCoord2f(pos.X+size.X, pos.Y+size.Y)
	rl.Vertex3f(c.X, c.Y, c.Z)

	rl.TexCoord2f(pos.X+size.X, pos.Y)
	rl.Vertex3f(d.X, d.Y, d.Z)

	rl.TexCoord2f(pos.X, pos.Y)
	rl.Vertex3f(a.X, a.Y, a.Z)
}
