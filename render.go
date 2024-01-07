package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Render struct {
	RenderTexture rl.RenderTexture2D
}

func NewRender() *Render {
	render := &Render{}

	SetWindowFlags()
	rl.InitWindow(
		int32(GameDisplay.Width),
		int32(GameDisplay.Height),
		GameDisplay.Title,
	)

	rl.SetWindowMinSize(
		int(GameDisplay.Width/WindowMinimalSizeDelimeter),
		GameDisplay.Height/WindowMinimalSizeDelimeter,
	)

	if GameDisplay.MaxFPS > 0 {
		rl.SetTargetFPS(int32(GameDisplay.MaxFPS))
	}

	render.RenderTexture = rl.LoadRenderTexture(int32(GameDisplay.Width), int32(GameDisplay.Height))
	rl.SetTextureFilter(render.RenderTexture.Texture, rl.TextureFilterLinear)

	return render
}

func (render *Render) Unload() {
	rl.UnloadRenderTexture(render.RenderTexture)
	rl.CloseWindow()
}

func SetWindowFlags() {
	var windowFlags uint32
	if GameDisplay.VSync {
		windowFlags |= rl.FlagVsyncHint
	}

	rl.SetConfigFlags(rl.FlagWindowResizable | windowFlags)
}

// HACK: Need scale in debug.go to figure out mouse delta...
var MouseScale rl.Vector2

func (r *Render) Draw(textureDraw, drawingDraw func()) {
	calculateDestinationRectangle := func() rl.Rectangle {
		scale := float32(rl.GetScreenHeight()) / float32(GameDisplay.Height)

		x0 := (float32(rl.GetScreenWidth()) - float32(GameDisplay.Width)*scale) / 2
		var y0 float32 = 0.0

		x1 := float32(GameDisplay.Width) * scale
		y1 := float32(rl.GetScreenHeight())

		// NOTE: Handle mouse offset and scaling when window resizes
		rl.SetMouseOffset(-int(x0), -int(y0))
		rl.SetMouseScale(1/scale, 1/scale)
		MouseScale = rl.NewVector2(scale, scale)

		return rl.NewRectangle(x0, y0, x1, y1)
	}

	rl.BeginTextureMode(r.RenderTexture)
	rl.ClearBackground(rl.Black)

	textureDraw()

	rl.EndTextureMode()
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawTexturePro(
		r.RenderTexture.Texture,
		rl.NewRectangle(0, 0, float32(r.RenderTexture.Texture.Width), -float32(r.RenderTexture.Texture.Height)),
		calculateDestinationRectangle(),
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)

	drawingDraw()

	rl.EndDrawing()
}

// FIXME: Origin seems to be wrong
func (r *Render) RenderTexture2D(
	texture *rl.Texture2D,
	positon rl.Vector2,
	size rl.Vector2,
	rotation float32,
	tint rl.Color,
) {
	textureSizeX := float32(texture.Width) * size.X
	textureSizeY := float32(texture.Height) * size.Y
	// origin := rl.NewVector2(
	// 	(float32(texture.Width)/2)*size.X,
	// 	(float32(texture.Height/2))*size.Y,
	// )

	rl.DrawTexturePro(
		*texture,
		rl.NewRectangle(0, 0, float32(texture.Width), -float32(texture.Height)),
		rl.NewRectangle(positon.X, positon.Y, textureSizeX, textureSizeY),
		rl.NewVector2(0, 0),
		rotation,
		tint,
	)
}

// @Hack: To call UpdateWindow in main thread. Because vars hot file reload works in goroutine
var ShouldUpdateWindowInMainThread = false

func (r *Render) Simulate() {
	if ShouldUpdateWindowInMainThread {
		rl.SetWindowSize(GameDisplay.Width, GameDisplay.Height)
		rl.SetWindowTitle(GameDisplay.Title)

		ShouldUpdateWindowInMainThread = false
	}
}
