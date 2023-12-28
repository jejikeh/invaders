package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const Aspect = 224.0 / 288.0
const VerticalPixels = 720

const WindowWidth = Aspect * VerticalPixels
const WindowHeight = VerticalPixels
const WindowMinimalSizeDelimeter = 1

const ResoursePath = "resources/"

type GameTextures struct {
	Player rl.Texture2D
}

type Entity struct {
	Texture      *rl.Texture2D
	Position     rl.Vector2
	Size         rl.Vector2
	Rotation     float32
	Tint         rl.Color
	ShadowHeight float32
	Visible      bool
}

func newEntity(
	texture rl.Texture2D,
	position rl.Vector2,
	size rl.Vector2,
	rotation float32,
	tint rl.Color,
) *Entity {
	entity := new(Entity)
	entity.Texture = &texture
	entity.Position = position
	entity.Size = size
	entity.Rotation = rotation
	entity.Tint = tint
	entity.ShadowHeight = 5
	entity.Visible = true

	return entity
}

func main() {
	initWindowAndOterStuff()
	defer rl.CloseWindow()

	renderTexture := initRenderTexture()
	defer rl.UnloadRenderTexture(renderTexture)

	font := initFont()
	defer rl.UnloadFont(font)

	gameTextures := initGameTextures()
	defer unloadGameTextures(gameTextures)

	player := newEntity(
		gameTextures.Player,
		rl.NewVector2(WindowWidth/2, WindowHeight/2),
		rl.NewVector2(.5, .5),
		.0,
		rl.RayWhite,
	)

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyEscape) {
			rl.CloseWindow()
			break
		}

		rl.BeginTextureMode(renderTexture)
		renderGradientBackground()

		if rl.IsKeyPressed(rl.KeySpace) {
			player.Visible = !player.Visible
		}

		if rl.IsKeyPressed(rl.KeyR) {
			player = nil
		}

		if player != nil {
			player.Rotation += 1
		}

		renderEntity(player)

		rl.EndTextureMode()

		drawRenderTexture(renderTexture)
	}
}

func initWindowAndOterStuff() {
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagVsyncHint)
	rl.InitWindow(
		WindowWidth,
		WindowHeight,
		"Invaders",
	)

	rl.SetWindowMinSize(
		WindowWidth/WindowMinimalSizeDelimeter,
		WindowHeight/WindowMinimalSizeDelimeter,
	)

	rl.SetTargetFPS(60)
}

func initRenderTexture() rl.RenderTexture2D {
	renderTexture := rl.LoadRenderTexture(WindowWidth, WindowHeight)
	rl.SetTextureFilter(renderTexture.Texture, rl.FilterPoint)

	return renderTexture
}

func drawRenderTexture(renderTexture rl.RenderTexture2D) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	rl.DrawTexturePro(
		renderTexture.Texture,
		rl.NewRectangle(0, 0, float32(renderTexture.Texture.Width), -float32(renderTexture.Texture.Height)),
		calculateDestinationRectangle(),
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)

	rl.EndDrawing()
}

func calculateDestinationRectangle() rl.Rectangle {
	scale := float32(rl.GetScreenHeight()) / WindowHeight

	x0 := (float32(rl.GetScreenWidth()) - float32(WindowWidth)*scale) / 2
	var y0 float32 = 0.0

	x1 := float32(WindowWidth) * scale
	y1 := float32(rl.GetScreenHeight())

	return rl.NewRectangle(x0, y0, x1, y1)
}

func initFont() rl.Font {
	return rl.LoadFont(ResoursePath + "font.ttf")
}

func initGameTextures() GameTextures {
	return GameTextures{
		Player: rl.LoadTexture(ResoursePath + "player.png"),
	}
}

func unloadGameTextures(gameTextures GameTextures) {
	rl.UnloadTexture(gameTextures.Player)
}

func renderGradientBackground() {
	rl.ClearBackground(rl.Black)

	rl.DrawRectangleGradientV(
		0,
		WindowHeight/1.5,
		WindowWidth,
		WindowHeight/3,
		rl.NewColor(c(.1), c(.1), c(.9), 255),
		rl.NewColor(c(.2), c(.5), c(.7), 255),
	)

	rl.DrawRectangleGradientV(
		0,
		0,
		WindowWidth,
		WindowHeight/1.5,
		rl.NewColor(c(.1), c(.1), c(.2), 255),
		rl.NewColor(c(.1), c(.1), c(.9), 255),
	)
}

func c(v float32) uint8 {
	return uint8(v * 255)
}

func renderEntity(entity *Entity) {
	if entity == nil || !entity.Visible {
		return
	}

	// render shadow
	renderTexture(
		entity.Texture,
		rl.Vector2Add(
			entity.Position,
			rl.NewVector2(entity.ShadowHeight, entity.ShadowHeight)),
		entity.Size,
		entity.Rotation,
		rl.Black,
	)

	renderTexture(
		entity.Texture,
		entity.Position,
		entity.Size,
		entity.Rotation,
		entity.Tint,
	)
}

func renderTexture(
	texture *rl.Texture2D,
	positon rl.Vector2,
	size rl.Vector2,
	rotation float32,
	tint rl.Color,
) {
	textureSizeX := float32(texture.Width) * size.X
	textureSizeY := float32(texture.Height) * size.Y
	origin := rl.NewVector2(
		(float32(texture.Width)/2)*size.X,
		(float32(texture.Height/2))*size.Y,
	)

	rl.DrawTexturePro(
		*texture,
		rl.NewRectangle(0, 0, float32(texture.Width), -float32(texture.Height)),
		rl.NewRectangle(positon.X, positon.Y, textureSizeX, textureSizeY),
		origin,
		rotation,
		tint,
	)
}
