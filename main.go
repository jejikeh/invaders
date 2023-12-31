package main

import (
	"math"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// TODO: Create a okay easing functions
// TODO: Delete emmiters?
// TODO: Optimize emmiters

const Aspect = 224.0 / 288.0
const VerticalPixels = 720

const WindowWidth = Aspect * VerticalPixels
const WindowHeight = VerticalPixels
const WindowMinimalSizeDelimeter = 1

const ResourseFolder = "resources/"
const FontFolder = ResourseFolder + "fonts/"

const EntitiesBaseSize = 0.3

const BigFontSize = 192
const SmallFontSize = 38

var Assets *AssetsManager
var Renderer *Render
var Entities *EntityManager
var Emitters *EmitterManager
var Debug *DebugHud

func main() {
	Debug = NewDebugHud()

	Renderer = NewRender()
	defer Renderer.Unload()

	Assets = NewAssetsManager()
	defer Assets.Unload()

	player := NewPlayer()
	Entities = NewEntityManager()
	Entities.Add(player)

	Emitters = NewEmitterManager()
	flameParticle := initFlameParticleSystem(rl.NewVector2(WindowWidth/2, WindowHeight/2))

	Emitters.AddHandlers(flameParticle,
		func(emitter *ParticleSystem) {
			emitter.Start()
		},
		func(emitter *ParticleSystem) {
			emitter.SetOrigin(player.Position)
		},
		func(emitter *ParticleSystem) {},
	)

	Entities.Add(Emitters)
	// Entities.Add(Debug)

	// invadersManager = &InvadersManager{}

	// invadersManager.Spawn(Dude, rl.NewVector2(WindowWidth/2, WindowHeight/4))
	// invadersManager.Spawn(Alien, rl.NewVector2(WindowWidth/2, WindowHeight/2))

	// player := newPlayer()

	// imageCircle := rl.GenImageGradientRadial(16, 16, 0.3, rl.White, rl.Black)
	// textureCircle = rl.LoadTextureFromImage(imageCircle)
	// defer rl.UnloadImage(imageCircle)
	// defer rl.UnloadTexture(textureCircle)

	// playerRockets = newPlayerRockets()

	// emitters = &Emitters{
	// Systems: []*ParticleSystem{},
	// }

	Entities.Start()
	Debug.Start()

	gui.LoadStyle("resources/cherry/style_cherry.rgs")

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyEscape) {
			rl.CloseWindow()
			break
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			flameParticle.SetLoop(!flameParticle.GetLoop())
		}

		Entities.Update()
		Debug.Update()

		Renderer.Draw(
			func() {
				renderGradientBackground()
				Entities.FirstDraw()
				Entities.Draw()

				// rl.DrawTextEx(*Assets.FontsManager.SmallFont, fmt.Sprintf("Frames: %d", rl.GetFPS()), rl.NewVector2(WindowWidth-130, 10), SmallFontSize, 0, rl.RayWhite)
				// rl.DrawTextEx(*Assets.FontsManager.SmallFont, fmt.Sprintf("Score: %d", player.Score), rl.NewVector2(10, 10), SmallFontSize, 0, rl.RayWhite)
				// rl.DrawTextEx(*Assets.FontsManager.SmallFont, fmt.Sprintf("Emitters: %d", Emitters.Count()), rl.NewVector2(10, 30), SmallFontSize, 0, rl.RayWhite)

				Debug.Draw()
			},
			func() {
				// NOTE: I thought use that for debug hud
			},
		)
		// invadersManager.DrawShadow()

		// invadersManager.Draw()
		// invadersManager.Update()

		// playerRockets.update()
		// playerRockets.draw()

		// emitters.Update()
		// emitters.Draw()

		// player.update()
		// player.updateEffects()
		// player.draw()

		// rl.EndTextureMode()
		// drawRenderTexture(renderTexture)
	}
}

// func (g *GameTextures) getInvaderTexture(invaderType InvaderType) (*rl.Texture2D, error) {
// 	switch invaderType {
// 	case Dude:
// 		return g.Dude, nil
// 	case Alien:
// 		return g.Alien, nil
// 	}

// 	return nil, errors.New("invader type not found")
// }

// type InvadersManager struct {
// 	Invaders []*Invader
// }

// func (invadersManager *InvadersManager) Spawn(invaderType InvaderType, position rl.Vector2) *Invader {
// 	texture, err := Textures.getInvaderTexture(invaderType)
// 	if err != nil {
// 		panic(err)
// 	}

// 	invader := &Invader{
// 		Entity: newEntity(
// 			*texture,
// 			position,
// 			rl.NewVector2(EntitiesBaseSize, EntitiesBaseSize),
// 			180,
// 			rl.RayWhite,
// 		),
// 		Type: invaderType,
// 	}

// 	invadersManager.Invaders = append(invadersManager.Invaders, invader)

// 	return invader
// }

// func (invadersManager *InvadersManager) Kill(invader *Invader) {
// 	for i, en := range invadersManager.Invaders {
// 		if en == invader {
// 			invadersManager.Invaders = append(invadersManager.Invaders[:i], invadersManager.Invaders[i+1:]...)
// 			break
// 		}
// 	}
// }

// func (invadersManager *InvadersManager) Draw() {
// 	for _, invader := range invadersManager.Invaders {
// 		renderEntity(invader.Entity)
// 	}
// }

// TODO: Think about render pipeline. How to draw this stuff more 'okeyish?'
// func (invadersManager *InvadersManager) DrawShadow() {
// 	for _, invader := range invadersManager.Invaders {
// 		// TODO: Rename renderEntityShadow to just renderShadow?
// 		invader.Entity.renderEntityShadow()
// 	}
// }

// func (invadersManager *InvadersManager) Update() {
// 	for _, invader := range invadersManager.Invaders {
// 		invader.update()
// 		invader.updateEffects()
// 	}
// }

// type EffectType int

// const (
// 	None EffectType = iota
// 	Respawn
// )

// type EntityEffect interface {
// 	getEntity() *Entity
// 	getLifetime() float32
// 	update()
// 	unset()
// 	getType() EffectType
// }

// const RespawnEffectDefaultLifetime = 1

// type RespanwEffect struct {
// 	Entity   *Entity
// 	Lifetime float32
// }

// func (ef *RespanwEffect) getEntity() *Entity {
// 	return ef.Entity
// }

// func (ef *RespanwEffect) getLifetime() float32 {
// 	return ef.Lifetime
// }

// func (ef *RespanwEffect) unset() {
// 	ef.Lifetime = 0
// 	ef.Entity.Tint = rl.RayWhite
// 	ef.Entity.removeEffect(ef)
// }

// func (ef *RespanwEffect) update() {
// 	ef.Lifetime -= rl.GetFrameTime()

// 	if ef.Lifetime <= 0 {
// 		ef.unset()
// 		return
// 	}

// 	// FIXME: Make it smooth and beautiful
// 	tintValue := float32(math.Sin(float64(p(1, RespawnEffectDefaultLifetime-ef.Lifetime) * 10)))

// 	// TODO: Maybe more interesting effect visualization?
// 	ef.Entity.Tint = rl.NewColor(c(tintValue), c(tintValue), c(tintValue), c(tintValue))

// 	shadowTintValue := float32(math.Abs(float64(tintValue / 4)))
// 	ef.Entity.ShadowTint = rl.NewColor(c(.1), c(.1), c(.1), c(shadowTintValue))
// }

// func (ef *RespanwEffect) getType() EffectType {
// 	return Respawn
// }

// func (entity *Entity) addEffect(effect EntityEffect) {
// 	entity.Effects = append(entity.Effects, effect)
// }

// func (entity *Entity) removeEffect(effect EntityEffect) {
// 	for i, ef := range entity.Effects {
// 		if ef == effect {
// 			entity.Effects = append(entity.Effects[:i], entity.Effects[i+1:]...)
// 			break
// 		}
// 	}
// }

// func (entity *Entity) updateEffects() {
// 	for _, effect := range entity.Effects {
// 		effect.update()
// 	}
// }

// func (entity *Entity) containsEffectOfType(effectType EffectType) bool {
// 	for _, effect := range entity.Effects {
// 		if effect.getType() == effectType {
// 			return true
// 		}
// 	}

// 	return false
// }

// type InvaderType int

// const (
// 	Dude InvaderType = iota
// 	Alien
// )

// type Invader struct {
// 	*Entity
// 	Type InvaderType
// }

// func (invader *Invader) update() {
// 	invader.ShadowHeight = float32(math.Sin(float64(rl.GetTime())/1.2)) * 8

// 	// if rl.IsKeyPressed(rl.KeyA) {
// 	// invader.Visible = !invader.Visible
// 	// }

// 	// TODO: Make it smooth and beautiful
// 	// This is fly 'simulation'. The same logic copied in player.update()
// 	// NOTE: Maybe it will be beter to move this to render or something like that
// 	// It even can be a separate effect though
// 	invader.Position.X += float32(math.Sin(float64(float32(rl.GetTime()*2)))) / 10
// 	invader.Position.Y += float32(math.Sin(float64(float32(rl.GetTime()*2)))) / 10
// }

// type Rocket struct {
// 	*Entity
// 	Speed         float32
// 	TrailEmitters *ParticleSystem
// }

// func newRocket(player *Player) *Rocket {
// 	rocket := &Rocket{
// 		Entity: newEntity(
// 			*Textures.Rocket,
// 			rl.NewVector2(player.Position.X, player.Position.Y-EntitiesBaseSize),
// 			rl.NewVector2(EntitiesBaseSize*1.5, EntitiesBaseSize*1.5),
// 			.0,
// 			rl.RayWhite,
// 		),
// 		Speed:         300,
// 		TrailEmitters: initFlameParticleSystem(player.Position),
// 	}

// 	rocket.TrailEmitters.Start()

// 	return rocket
// }

// type PlayerRockets struct {
// 	Rockets []*Rocket
// }

// func newPlayerRockets() *PlayerRockets {
// 	return &PlayerRockets{
// 		Rockets: []*Rocket{},
// 	}
// }

// func (p *PlayerRockets) addRocket(player *Player) {
// 	p.Rockets = append(p.Rockets, newRocket(player))
// }

// func (p *PlayerRockets) update() {
// 	for i, rocket := range p.Rockets {
// 		rocket.update()

// 		if rocket.Position.Y < -1000 {
// 			// TODO: Refactor this. Maybe refactor entire PlayerRockets thing
// 			rocket.TrailEmitters.Stop()
// 			rocket.TrailEmitters = nil

// 			if p.Rockets[i] != nil {
// 				p.Rockets[i] = nil
// 			}

// 			p.Rockets = append(p.Rockets[:i], p.Rockets[i+1:]...)
// 		}

// 		// TODO: Maybe crate something like event system with "Rocket Collided" and so on
// 		// TODO: Or look how it`s was done in jai invaders
// 		for _, ii := range invadersManager.Invaders {
// 			if ii.CollidesWith(rocket) {
// 				invadersManager.Kill(ii)
// 				rocket.TrailEmitters.Stop()
// 				rocket.TrailEmitters.Stop()
// 				rocket.TrailEmitters = nil

// 				if p.Rockets[i] != nil {
// 					p.Rockets[i] = nil
// 				}

// 				p.Rockets = append(p.Rockets[:i], p.Rockets[i+1:]...)

// 				emitters.Add(initExplosionParticleSystem(rl.NewVector2(rocket.Position.X, rocket.Position.Y)))
// 				emitters.Burst()

// 				return
// 			}
// 		}
// 	}
// }

// func (p *PlayerRockets) draw() {
// 	for _, rocket := range p.Rockets {
// 		rocket.draw()
// 	}
// }

// func (rocket *Rocket) update() {
// 	rocket.Position.Y -= rocket.Speed * float32(rl.GetFrameTime())

// 	trailEffectPosition := rl.NewVector2(rocket.Position.X+5, rocket.Position.Y+10)

// 	rocket.TrailEmitters.SetOrigin(trailEffectPosition)
// 	rocket.TrailEmitters.Update()
// }

// func (i *Invader) CollidesWith(rocket *Rocket) bool {
// 	return rl.CheckCollisionRecs(
// 		rl.NewRectangle(rocket.Position.X, rocket.Position.Y, float32(rocket.Entity.Texture.Width)*rocket.Size.X, float32(rocket.Entity.Texture.Height)*rocket.Size.Y),
// 		rl.NewRectangle(i.Position.X, i.Position.Y, float32(i.Texture.Width)*i.Size.X, float32(i.Texture.Height)*i.Size.Y),
// 	)
// }

// func (rocket *Rocket) draw() {
// 	rocket.Entity.renderEntityShadow()
// 	rocket.TrailEmitters.Draw()
// 	renderEntity(rocket.Entity)
// }

// TODO: remove that from global scope apparently
// var textureCircle rl.Texture2D

// var playerRockets *PlayerRockets
// var invadersManager *InvadersManager
// var emitters *Emitters

// type Emitters struct {
// 	Systems []*ParticleSystem
// }

// func (e *Emitters) Add(emitter *ParticleSystem) {
// 	e.Systems = append(e.Systems, emitter)
// }

// func (e *Emitters) Remove(emitter *ParticleSystem) {
// 	for i, s := range e.Systems {
// 		if s == emitter {
// 			e.Systems = append(e.Systems[:i], e.Systems[i+1:]...)
// 		}
// 	}
// }

// func (e *Emitters) Update() {
// 	for _, s := range e.Systems {
// 		s.Update()
// 	}
// }

// func (e *Emitters) Draw() {
// 	for _, s := range e.Systems {
// 		s.Draw()
// 	}
// }

// func (e *Emitters) Burst() {
// 	for _, s := range e.Systems {
// 		s.Burst()
// 	}
// }

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
	// shadowColor := entity.ShadowTint

	// if !entity.containsEffectOfType(Respawn) {
	// 	shadowColor.A = c(1 / entityHeight)
	// }

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

// func (entity *Entity) renderEntityShadow() {
// 	if entity == nil || !entity.Visible {
// 		return
// 	}

// 	entityHeight := float32(math.Abs(float64(entity.ShadowHeight))) + 2.5
// 	shadowColor := entity.ShadowTint

// 	if !entity.containsEffectOfType(Respawn) {
// 		shadowColor.A = c(1 / entityHeight)
// 	}

// 	// render shadow
// 	renderTexture(
// 		entity.Texture,
// 		rl.Vector2AddValue(entity.Position, entityHeight),
// 		rl.Vector2AddValue(entity.Size, entityHeight/100),
// 		entity.Rotation,
// 		shadowColor,
// 	)
// }

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

const FlameSize = 0.5

func initFlameParticleSystem(origin rl.Vector2) *ParticleSystem {
	ps := &ParticleSystem{}

	configFlame1 := EmitterConfig{
		Loop:         true,
		StartSize:    rl.NewVector2(2*FlameSize, 2*FlameSize),
		EndSize:      rl.NewVector2(1*FlameSize, 1*FlameSize),
		Capacity:     100,
		EmmisionRate: 500,
		Origin:       origin,
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
			0.0,
			1.2,
		},
		Texture:   Assets.TexturesManager.Circle,
		BlendMode: rl.BlendAdditive,
	}

	emitterFlame1 := NewEmitter(configFlame1)

	configFlame2 := configFlame1
	configFlame2.StartSize = rl.NewVector2(2*FlameSize, 2*FlameSize)
	configFlame2.EndSize = rl.NewVector2(0*FlameSize, 0*FlameSize)
	configFlame2.Capacity = 20
	configFlame2.EmmisionRate = 20
	configFlame2.StartColor = rl.NewColor(255, 255, 255, 10)
	configFlame2.EndColor = rl.NewColor(255, 255, 255, 0)
	configFlame2.Age = [2]float32{
		0.0,
		1.0,
	}

	emitterFlame2 := NewEmitter(configFlame2)

	configSmokeEmitter := configFlame2
	configSmokeEmitter.StartSize = rl.NewVector2(2*FlameSize, 2*FlameSize)
	configSmokeEmitter.EndSize = rl.NewVector2(1*FlameSize, 1*FlameSize)
	configSmokeEmitter.Capacity = 100
	configSmokeEmitter.EmmisionRate = 100
	configSmokeEmitter.StartColor = rl.NewColor(125, 125, 125, 30)
	configSmokeEmitter.EndColor = rl.NewColor(125, 125, 125, 10)
	configSmokeEmitter.Age = [2]float32{
		0.0,
		1.5,
	}

	smokeEmitter := NewEmitter(configSmokeEmitter)

	ps.Add(emitterFlame1)
	ps.Add(smokeEmitter)
	ps.Add(emitterFlame2)

	return ps
}
