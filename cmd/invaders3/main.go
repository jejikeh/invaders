package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"io/fs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// @Incomplete: Rotations are completely broken.

//go:embed assets/*
var assets embed.FS

func main() {
	// @Incomplete: Handle window config flags here. Maybe, at some point remove
	const width = 1024
	const height = 768

	rl.InitWindow(width, height, "invaders 3")
	defer rl.CloseWindow()

	camera := rl.Camera3D{
		// @Cleanup: Fix strange artifacts with the quad rendering. X: 0.1
		Position:   rl.Vector3{X: 0.1, Y: 20, Z: 10},
		Target:     rl.Vector3{X: 0, Y: 0, Z: 0},
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       60,
		Projection: rl.CameraPerspective,
	}

	tilesPacked, unloadAssets, err := NewTexture(assets, "assets/Tilemap/tiles_packed.png")
	if err != nil {
		panic(err)
	}
	defer unloadAssets()

	shipsPacked, unloadAssets, err := NewTexture(assets, "assets/Tilemap/ships_packed.png")
	if err != nil {
		panic(err)
	}
	defer unloadAssets()

	// shader, unloadShader, err := NewShader(assets, "", "assets/shaders/grayscale.fs")
	// if err != nil {
	// 	panic(err)
	// }
	// defer unloadShader()

	tilesAtlas := NewAtlas(tilesPacked, rl.Vector2{X: 12, Y: 10})
	shipAtlas := NewAtlas(shipsPacked, rl.Vector2{X: 4, Y: 6})

	renderTexture := rl.LoadRenderTexture(width, height)
	defer rl.UnloadRenderTexture(renderTexture)

	renderTexture1 := rl.LoadRenderTexture(width, height)
	defer rl.UnloadRenderTexture(renderTexture)

	renderTexture2 := rl.LoadRenderTexture(width, height)
	defer rl.UnloadRenderTexture(renderTexture2)

	playerPos := rl.Vector3{X: 4, Y: 1, Z: 4}
	playerVel := rl.Vector3{X: 0, Y: 0, Z: 0}
	var playerRot float32 = 0.0
	cameraPos := rl.Vector3{X: 0, Y: 0, Z: 0}

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		var velocity rl.Vector3
		if rl.IsKeyDown(rl.KeyW) {
			velocity.Z -= 1
		}

		if rl.IsKeyDown(rl.KeyS) {
			velocity.Z += 1
		}

		if rl.IsKeyDown(rl.KeyA) {
			velocity.X -= 1
		}

		if rl.IsKeyDown(rl.KeyD) {
			velocity.X += 1
		}

		if rl.IsKeyReleased(rl.KeyW) || rl.IsKeyReleased(rl.KeyS) {
			velocity.Z = 0
		}

		if rl.IsKeyReleased(rl.KeyA) || rl.IsKeyReleased(rl.KeyD) {
			velocity.X = 0
		}

		velocity = rl.Vector3Scale(rl.Vector3Normalize(velocity), 10*dt)

		playerVel = rl.Vector3Lerp(playerVel, velocity, 10*dt)
		playerPos = rl.Vector3Add(playerPos, playerVel)
		playerRot = rl.Lerp(playerRot, -velocity.X*5, 10*dt)

		cameraPos.X = rl.Lerp(cameraPos.X, -velocity.Z, 2*dt)

		camera.Target.Z = rl.Lerp(camera.Target.Z, playerPos.Z-10, 2*dt)

		// dt := rl.GetFrameTime()
		// rl.UpdateCamera(&camera, rl.CameraOrbital)
		// rl.UpdateCamera(&camera, rl.CameraFree)
		rl.UpdateCameraPro(&camera, cameraPos, rl.Vector3{}, 0)

		textureRender := func(r rl.RenderTexture2D, render func()) {
			rl.BeginTextureMode(r)
			rl.ClearBackground(rl.Color{})

			rl.BeginMode3D(camera)

			rl.Begin(rl.Triangles)

			render()

			rl.End()
			rl.EndMode3D()
			rl.EndTextureMode()
		}

		combineTextures := func(r rl.RenderTexture2D, ts ...rl.RenderTexture2D) {
			rl.BeginBlendMode(rl.BlendAlpha)
			rl.BeginTextureMode(r)

			rl.ClearBackground(rl.Black)

			for _, t := range ts {
				rl.DrawTexturePro(
					t.Texture,
					rl.NewRectangle(0, 0, float32(t.Texture.Width), float32(t.Texture.Height)),
					rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
					rl.Vector2{X: 0, Y: 0},
					0,
					rl.White,
				)
			}

			rl.EndTextureMode()
			rl.EndBlendMode()
		}

		present := func(t rl.RenderTexture2D) {
			rl.BeginDrawing()

			rl.ClearBackground(rl.Black)

			rl.DrawTexturePro(
				t.Texture,
				rl.NewRectangle(0, 0, float32(t.Texture.Width), float32(t.Texture.Height)),
				rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
				rl.NewVector2(0, 0),
				0,
				rl.White,
			)

			rl.EndBlendMode()
			rl.EndDrawing()
		}

		textureRender(renderTexture1, func() {
			const size = 100
			for z := range size {
				z := z - size/2
				for x := range size {
					x := x - size/2
					DrawQuad(
						tilesAtlas,
						rl.Vector3{X: float32(x), Y: 0.2, Z: float32(z)},
						rl.Vector3{X: 1, Y: 1, Z: 1},
						(12*10)-10,
					)
				}
			}
		})

		textureRender(renderTexture2, func() {
			for z := range 6 {
				DrawQuad(
					shipAtlas,
					rl.Vector3{Y: 3 + float32(z)/2, Z: float32(z)},
					rl.Vector3{X: 2, Y: 1, Z: 2},
					uint32(z),
				)
			}

			DrawQuad(
				shipAtlas,
				playerPos,
				rl.Vector3{X: 2, Y: 1, Z: 2},
				8,
				func(v rl.Vector3) rl.Vector3 {
					return rl.Vector3RotateByAxisAngle(v, rl.Vector3{X: playerPos.X, Y: playerPos.Y, Z: playerPos.Z}, playerRot)
				},
			)
		})

		combineTextures(renderTexture, renderTexture1, renderTexture2)

		present(renderTexture)
	}
}

func NewShader(f fs.FS, vertex, fragment string) (*rl.Shader, func(), error) {
	// @Incomplete: Handle errors
	readFile := func(name string) string {
		file, err := f.Open(name)
		if err != nil {
			return ""
		}

		r := bufio.NewReader(file)

		texBuf := make([]byte, 0, 1024)
		buf := make([]byte, 1024)
		for {
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				return ""
			}

			if n == 0 {
				break
			}

			texBuf = append(texBuf, buf[:n]...)
		}

		return string(texBuf)
	}

	shader := rl.LoadShaderFromMemory(readFile(vertex), readFile(fragment))

	return &shader, func() { rl.UnloadShader(shader) }, nil
}

func NewTexture(f fs.FS, name string) (*rl.Texture2D, func(), error) {
	file, err := f.Open(name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load texture: %w", err)
	}

	r := bufio.NewReader(file)

	texBuf := make([]byte, 0, 1024)
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return nil, nil, fmt.Errorf("failed to read texture: %w", err)
		}

		if n == 0 {
			break
		}

		texBuf = append(texBuf, buf[:n]...)
	}

	texture := rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", texBuf, int32(len(texBuf))))

	return &texture, func() { rl.UnloadTexture(texture) }, nil
}

type RegionAtlas struct {
	texture    *rl.Texture2D
	dimensions rl.Vector2
}

func NewAtlas(texture *rl.Texture2D, dimensions rl.Vector2) RegionAtlas {
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

func DrawQuad(atlas RegionAtlas, position, scale rl.Vector3, idx uint32, modifyVertex ...func(rl.Vector3) rl.Vector3) {
	rl.SetTexture(atlas.texture.ID)

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

	pos, size := atlas.getPositionAt(idx)

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
