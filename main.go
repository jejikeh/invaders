package main

import (
	"errors"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const Aspect = 224.0 / 288.0
const VerticalPixels = 720

const WindowWidth = Aspect * VerticalPixels
const WindowHeight = VerticalPixels
const WindowMinimalSizeDelimeter = 1

const ResoursePath = "resources/"

const EntitiesBaseSize = 0.3

type GameTextures struct {
	Player *rl.Texture2D
	Dude   *rl.Texture2D
	Alien  *rl.Texture2D
}

func (g *GameTextures) getInvaderTexture(invaderType InvaderType) (*rl.Texture2D, error) {
	switch invaderType {
	case Dude:
		return g.Dude, nil
	case Alien:
		return g.Alien, nil
	}

	return nil, errors.New("Invader type not found")
}

type InvadersManager struct {
	Invaders []*Invader
	Textures *GameTextures
}

func (invadersManager *InvadersManager) Spawn(invaderType InvaderType, position rl.Vector2) *Invader {
	texture, err := invadersManager.Textures.getInvaderTexture(invaderType)
	if err != nil {
		panic(err)
	}

	invader := &Invader{
		Entity: newEntity(
			*texture,
			position,
			rl.NewVector2(EntitiesBaseSize, EntitiesBaseSize),
			180,
			rl.RayWhite,
		),
		Type: invaderType,
	}

	invadersManager.Invaders = append(invadersManager.Invaders, invader)

	return invader
}

func (invadersManager *InvadersManager) Kill(invader *Invader) {
	for i, en := range invadersManager.Invaders {
		if en == invader {
			invadersManager.Invaders = append(invadersManager.Invaders[:i], invadersManager.Invaders[i+1:]...)
			break
		}
	}
}

func (invadersManager *InvadersManager) Draw() {
	for _, invader := range invadersManager.Invaders {
		renderEntity(invader.Entity)
	}
}

func (invadersManager *InvadersManager) Update() {
	for _, invader := range invadersManager.Invaders {
		invader.update()
	}
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

type InvaderType int

const (
	Dude InvaderType = iota
	Alien
)

type Invader struct {
	*Entity
	Type InvaderType
}

func (invader *Invader) update() {
	invader.ShadowHeight = float32(math.Sin(float64(rl.GetTime())/1.2)) * 8
}

type Player struct {
	*Entity
	Speed float32
	Score int
	Lives int
}

func newPlayer(texture rl.Texture2D) *Player {
	return &Player{
		Entity: newEntity(
			texture,
			rl.NewVector2(WindowWidth/2, WindowHeight/2),
			rl.NewVector2(1, 1),
			.0,
			rl.RayWhite,
		),
		Speed: 10,
	}
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

	invadersManager := &InvadersManager{
		Textures: gameTextures,
	}

	invadersManager.Spawn(Dude, rl.NewVector2(WindowWidth/2, WindowHeight/4))
	invadersManager.Spawn(Alien, rl.NewVector2(WindowWidth/2, WindowHeight/2))

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyEscape) {
			rl.CloseWindow()
			break
		}

		rl.BeginTextureMode(renderTexture)
		renderGradientBackground()

		invadersManager.Draw()
		invadersManager.Update()

		rl.EndTextureMode()

		drawRenderTexture(renderTexture)
	}
}

func updatePlayer(player *Player) {
	if player == nil {
		return
	}

	player.ShadowHeight += float32(math.Sin(float64(rl.GetTime())))
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
	rl.SetTextureFilter(renderTexture.Texture, rl.TextureFilterLinear)

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

func initGameTextures() *GameTextures {
	player := rl.LoadTexture(ResoursePath + "player.png")
	dude := rl.LoadTexture(ResoursePath + "dude.png")
	alien := rl.LoadTexture(ResoursePath + "alien.png")

	gameTexture := new(GameTextures)
	gameTexture.Player = &player
	gameTexture.Dude = &dude
	gameTexture.Alien = &alien

	return gameTexture
}

func unloadGameTextures(gameTextures *GameTextures) {
	rl.UnloadTexture(*gameTextures.Player)
	rl.UnloadTexture(*gameTextures.Dude)
	rl.UnloadTexture(*gameTextures.Alien)
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

// Create a color from 0 to 255 using percentage
// TODO: remove this and move to 255 values
func c(v float32) uint8 {
	return uint8(v * 255)
}

func renderEntity(entity *Entity) {
	if entity == nil || !entity.Visible {
		return
	}

	entityHeight := float32(math.Abs(float64(entity.ShadowHeight))) + 2.5

	// render shadow
	renderTexture(
		entity.Texture,
		rl.Vector2AddValue(entity.Position, entityHeight),
		rl.Vector2AddValue(entity.Size, entityHeight/100),
		entity.Rotation,
		rl.NewColor(c(.1), c(.1), c(.1), c(1/entityHeight)),
	)

	// render entity
	renderTexture(
		entity.Texture,
		rl.Vector2Subtract(entity.Position, rl.NewVector2(0, entityHeight)),
		rl.Vector2AddValue(entity.Size, entityHeight/100),
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
