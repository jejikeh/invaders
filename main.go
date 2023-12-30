package main

import (
	"errors"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// TODO: Create a okay easing functions

const Aspect = 224.0 / 288.0
const VerticalPixels = 720

const WindowWidth = Aspect * VerticalPixels
const WindowHeight = VerticalPixels
const WindowMinimalSizeDelimeter = 1

const ResourseFolder = "resources/"
const FontFolder = ResourseFolder + "fonts/"

const EntitiesBaseSize = 0.3

const BigFontSize = 192
const SmallFontSize = BigFontSize * 2 / 3

type GameTextures struct {
	Player *rl.Texture2D
	Dude   *rl.Texture2D
	Alien  *rl.Texture2D
}

var Textures *GameTextures

func (g *GameTextures) getInvaderTexture(invaderType InvaderType) (*rl.Texture2D, error) {
	switch invaderType {
	case Dude:
		return g.Dude, nil
	case Alien:
		return g.Alien, nil
	}

	return nil, errors.New("invader type not found.")
}

type InvadersManager struct {
	Invaders []*Invader
}

func (invadersManager *InvadersManager) Spawn(invaderType InvaderType, position rl.Vector2) *Invader {
	texture, err := Textures.getInvaderTexture(invaderType)
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

// TODO: Think about render pipeline. How to draw this stuff more 'okeyish?'
func (invadersManager *InvadersManager) DrawShadow() {
	for _, invader := range invadersManager.Invaders {
		// TODO: Rename renderEntityShadow to just renderShadow?
		invader.Entity.renderEntityShadow()
	}
}

func (invadersManager *InvadersManager) Update() {
	for _, invader := range invadersManager.Invaders {
		invader.update()
		invader.updateEffects()
	}
}

type EffectType int

const (
	None EffectType = iota
	Respawn
)

type EntityEffect interface {
	getEntity() *Entity
	getLifetime() float32
	update()
	unset()
	getType() EffectType
}

const RespawnEffectDefaultLifetime = 1

type RespanwEffect struct {
	Entity   *Entity
	Lifetime float32
}

func (ef *RespanwEffect) getEntity() *Entity {
	return ef.Entity
}

func (ef *RespanwEffect) getLifetime() float32 {
	return ef.Lifetime
}

func (ef *RespanwEffect) unset() {
	ef.Lifetime = 0
	ef.Entity.Tint = rl.RayWhite
	ef.Entity.removeEffect(ef)
}

func (ef *RespanwEffect) update() {
	ef.Lifetime -= rl.GetFrameTime()

	if ef.Lifetime <= 0 {
		ef.unset()
		return
	}

	// FIXME: Make it smooth and beautiful
	tintValue := float32(math.Sin(float64(p(1, RespawnEffectDefaultLifetime-ef.Lifetime) * 10)))

	// TODO: Maybe more interesting effect visualization?
	ef.Entity.Tint = rl.NewColor(c(tintValue), c(tintValue), c(tintValue), c(tintValue))

	shadowTintValue := float32(math.Abs(float64(tintValue / 4)))
	ef.Entity.ShadowTint = rl.NewColor(c(.1), c(.1), c(.1), c(shadowTintValue))
}

func (ef *RespanwEffect) getType() EffectType {
	return Respawn
}

func spike(t float32) float32 {
	if t <= .5 {
		return easeIn(t / .5)
	}

	return easeIn(flip(t) / .5)
}

func flip(t float32) float32 {
	return 1 - t
}

func easeIn(t float32) float32 {
	return t * t
}

type Entity struct {
	Texture      *rl.Texture2D
	Position     rl.Vector2
	Size         rl.Vector2
	Rotation     float32
	Tint         rl.Color
	ShadowHeight float32
	ShadowTint   rl.Color
	Visible      bool
	Effects      []EntityEffect
}

func (entity *Entity) addEffect(effect EntityEffect) {
	entity.Effects = append(entity.Effects, effect)
}

func (entity *Entity) removeEffect(effect EntityEffect) {
	for i, ef := range entity.Effects {
		if ef == effect {
			entity.Effects = append(entity.Effects[:i], entity.Effects[i+1:]...)
			break
		}
	}
}

func (entity *Entity) updateEffects() {
	for _, effect := range entity.Effects {
		effect.update()
	}
}

func (entity *Entity) containsEffectOfType(effectType EffectType) bool {
	for _, effect := range entity.Effects {
		if effect.getType() == effectType {
			return true
		}
	}

	return false
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
	entity.ShadowTint = rl.NewColor(c(.1), c(.1), c(.1), c(.1))
	entity.Visible = true
	entity.addEffect(&RespanwEffect{
		Entity:   entity,
		Lifetime: RespawnEffectDefaultLifetime,
	})

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

	// if rl.IsKeyPressed(rl.KeyA) {
	// invader.Visible = !invader.Visible
	// }

	// TODO: Make it smooth and beautiful
	// This is fly 'simulation'. The same logic copied in player.update()
	// NOTE: Maybe it will be beter to move this to render or something like that
	// It even can be a separate effect though
	invader.Position.X += float32(math.Sin(float64(float32(rl.GetTime()*2)))) / 10
	invader.Position.Y += float32(math.Sin(float64(float32(rl.GetTime()*2)))) / 10
}

type Player struct {
	*Entity
	Speed        float32
	Score        int
	Lives        int
	MaxSpeed     float32
	Acceleration float32
	Friction     float32
	Velocity     rl.Vector2
}

func newPlayer() *Player {
	return &Player{
		Entity: newEntity(
			*Textures.Player,
			rl.NewVector2(WindowWidth/2, WindowHeight/1.2),
			rl.NewVector2(EntitiesBaseSize, EntitiesBaseSize),
			.0,
			rl.RayWhite,
		),
		Speed:        5,
		Lives:        3,
		MaxSpeed:     400,
		Acceleration: 1500,
		Friction:     600 / 2,
	}
}

var textureCircle rl.Texture2D

func main() {
	initWindowAndOterStuff()
	defer rl.CloseWindow()

	renderTexture := initRenderTexture()
	defer rl.UnloadRenderTexture(renderTexture)

	bigFont := initFont(BigFontSize)
	defer rl.UnloadFont(bigFont)

	// TODO: Make font manager as global variable or just move small and big fonts to global scope
	smallFont := initFont(SmallFontSize)
	defer rl.UnloadFont(smallFont)

	Textures = initGameTextures()
	defer unloadGameTextures(Textures)

	invadersManager := &InvadersManager{}

	invadersManager.Spawn(Dude, rl.NewVector2(WindowWidth/2, WindowHeight/4))
	invadersManager.Spawn(Alien, rl.NewVector2(WindowWidth/2, WindowHeight/2))

	player := newPlayer()

	imageCircle := rl.GenImageGradientRadial(16, 16, 0.3, rl.White, rl.Black)
	textureCircle = rl.LoadTextureFromImage(imageCircle)
	defer rl.UnloadImage(imageCircle)
	defer rl.UnloadTexture(textureCircle)

	ps := initFlameParticleSystem()
	ps.Start()

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyEscape) {
			rl.CloseWindow()
			break
		}

		rl.BeginTextureMode(renderTexture)
		renderGradientBackground()

		player.renderEntityShadow()
		invadersManager.DrawShadow()

		ps.SetOrigin(player.Position)
		ps.Update()
		ps.Draw()

		invadersManager.Draw()
		invadersManager.Update()

		player.update()
		player.updateEffects()
		renderEntity(player.Entity)

		rl.DrawTextEx(smallFont, fmt.Sprintf("Frames: %d", rl.GetFPS()), rl.NewVector2(WindowWidth-130, 10), BigFontSize*0.2, 0, rl.RayWhite)
		rl.DrawTextEx(smallFont, fmt.Sprintf("Score: %d", player.Score), rl.NewVector2(10, 10), BigFontSize*0.2, 0, rl.RayWhite)

		rl.EndTextureMode()

		drawRenderTexture(renderTexture)
	}
}

func (p *Player) update() {
	if p == nil {
		return
	}

	p.ShadowHeight = float32(math.Sin(float64(rl.GetTime())/1.2)) * 8

	if rl.IsKeyPressed(rl.KeyP) {
		p.addEffect(&RespanwEffect{
			Entity:   p.Entity,
			Lifetime: RespawnEffectDefaultLifetime,
		})
	}

	p.updateMovement()

	p.Position.X += float32(math.Sin(float64(float32(rl.GetTime()*2)))) / 10
	p.Position.Y += float32(math.Sin(float64(float32(rl.GetTime()*2)))) / 10

	if p.Position.Y > WindowHeight {
		p.Position.Y = WindowHeight
	}

	if p.Position.Y < 0 {
		p.Position.Y = 0
	}

	if p.Position.X > WindowWidth {
		p.Position.X = 0
	}

	if p.Position.X < 0 {
		p.Position.X = WindowWidth
	}
}

func (p *Player) updateMovement() {
	input := getInputVector()

	if input.X == 0 && input.Y == 0 {
		if rl.Vector2Length(p.Velocity) > p.Friction*rl.GetFrameTime() {
			p.Velocity = rl.Vector2Subtract(
				p.Velocity,
				rl.Vector2Multiply(
					rl.Vector2Normalize(p.Velocity),
					rl.NewVector2(p.Friction*rl.GetFrameTime(), p.Friction*rl.GetFrameTime()),
				),
			)
		} else {
			p.Velocity = rl.NewVector2(0, 0)
		}
	} else {
		p.Velocity = rl.Vector2Add(
			p.Velocity,
			rl.Vector2Multiply(
				input,
				rl.NewVector2(p.Acceleration*rl.GetFrameTime(), p.Acceleration*rl.GetFrameTime()),
			),
		)

		p.Velocity.X = rl.Clamp(p.Velocity.X, -p.MaxSpeed, p.MaxSpeed)
		p.Velocity.Y = rl.Clamp(p.Velocity.Y, -p.MaxSpeed, p.MaxSpeed)
	}

	p.Position.X += p.Velocity.X * rl.GetFrameTime()
	p.Position.Y += p.Velocity.Y * rl.GetFrameTime()
}

func getInputVector() rl.Vector2 {
	var inputVector rl.Vector2
	inputVector.X = btof32(rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD)) - btof32(rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA))
	inputVector.Y = btof32(rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS)) - btof32(rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW))

	inputVector = rl.Vector2Normalize(inputVector)

	return inputVector
}

func btof32(b bool) float32 {
	if b {
		return 1
	}
	return 0
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

func initFont(size int32) rl.Font {
	return rl.LoadFontEx(FontFolder+"Martel-Regular.ttf", size, nil)
}

func initGameTextures() *GameTextures {
	player := rl.LoadTexture(ResourseFolder + "player.png")
	dude := rl.LoadTexture(ResourseFolder + "dude.png")
	alien := rl.LoadTexture(ResourseFolder + "alien.png")

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

func p(t float32, b float32) float32 {
	return t / (t + b)
}

func renderEntity(entity *Entity) {
	if entity == nil || !entity.Visible {
		return
	}

	entityHeight := float32(math.Abs(float64(entity.ShadowHeight))) + 2.5
	shadowColor := entity.ShadowTint

	if !entity.containsEffectOfType(Respawn) {
		shadowColor.A = c(1 / entityHeight)
	}

	// // render shadow
	// renderTexture(
	// 	entity.Texture,
	// 	rl.Vector2AddValue(entity.Position, entityHeight),
	// 	rl.Vector2AddValue(entity.Size, entityHeight/100),
	// 	entity.Rotation,
	// 	shadowColor,
	// )

	// render entity
	renderTexture(
		entity.Texture,
		rl.Vector2Subtract(entity.Position, rl.NewVector2(0, entityHeight)),
		rl.Vector2AddValue(entity.Size, entityHeight/100),
		entity.Rotation,
		entity.Tint,
	)
}

func (entity *Entity) renderEntityShadow() {
	if entity == nil || !entity.Visible {
		return
	}

	entityHeight := float32(math.Abs(float64(entity.ShadowHeight))) + 2.5
	shadowColor := entity.ShadowTint

	if !entity.containsEffectOfType(Respawn) {
		shadowColor.A = c(1 / entityHeight)
	}

	// render shadow
	renderTexture(
		entity.Texture,
		rl.Vector2AddValue(entity.Position, entityHeight),
		rl.Vector2AddValue(entity.Size, entityHeight/100),
		entity.Rotation,
		shadowColor,
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

func initFlameParticleSystem() *ParticleSystem {
	ps := &ParticleSystem{}

	configFlame1 := EmitterConfig{
		StartSize:    rl.NewVector2(2, 2),
		EndSize:      rl.NewVector2(1, 1),
		Capacity:     100,
		EmmisionRate: 500,
		Origin: rl.Vector2{
			X: 0,
			Y: 0,
		},
		OriginAcceleration: [2]float32{
			50,
			100,
		},
		Offset: [2]float32{
			0,
			10,
		},
		Direction: rl.NewVector2(0, -1),
		DirectionAngle: [2]float32{
			90,
			90,
		},
		Velocity: [2]float32{
			30,
			150,
		},
		VelocityAngle: [2]float32{
			90,
			90,
		},
		StartColor: rl.NewColor(255, 20, 0, 255),
		EndColor:   rl.NewColor(255, 20, 0, 0),
		Age: [2]float32{
			1.0,
			2.0,
		},
		Texture:   textureCircle,
		BlendMode: rl.BlendAdditive,
	}

	emitterFlame1 := NewEmitter(configFlame1)

	configFlame2 := configFlame1
	configFlame2.StartSize = rl.NewVector2(1, 1)
	configFlame2.EndSize = rl.NewVector2(0, 0)
	configFlame2.Capacity = 20
	configFlame2.EmmisionRate = 20
	configFlame2.StartColor = rl.NewColor(255, 255, 255, 10)
	configFlame2.EndColor = rl.NewColor(255, 255, 255, 0)
	configFlame2.Age = [2]float32{
		0.5,
		1.0,
	}

	emitterFlame2 := NewEmitter(configFlame2)

	configSmokeEmitter := configFlame2
	configSmokeEmitter.StartSize = rl.NewVector2(2, 2)
	configSmokeEmitter.EndSize = rl.NewVector2(1, 1)
	configSmokeEmitter.Capacity = 500
	configSmokeEmitter.EmmisionRate = 100
	configSmokeEmitter.StartColor = rl.NewColor(125, 125, 125, 30)
	configSmokeEmitter.EndColor = rl.NewColor(125, 125, 125, 10)
	configSmokeEmitter.Age = [2]float32{
		3.0,
		5.0,
	}

	smokeEmitter := NewEmitter(configSmokeEmitter)

	ps.Add(emitterFlame1)
	ps.Add(smokeEmitter)
	ps.Add(emitterFlame2)

	return ps
}
