package main

import (
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
var Debug *Hud

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
	flameParticle := initFlameParticleSystem(rl.NewVector2(player.Position.X, player.Position.Y))

	Emitters.AddHandlers(flameParticle,
		func(emitter *ParticleSystem) {
			// emitter.Stop()
		},
		func(emitter *ParticleSystem) {
			rect := player.GetRectangle()
			emitter.SetOrigin(rl.NewVector2((rect.X + rect.Width/2), rect.Y+rect.Height))
		},
		func(emitter *ParticleSystem) {},
	)

	Entities.Start()
	Emitters.Start()

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
		Emitters.Update()

		Debug.Update()

		Renderer.Draw(
			func() {
				// TODO: Move as separate entity?
				renderGradientBackground()

				Entities.FirstDraw()

				Emitters.Draw()
				Entities.Draw()

				Debug.Draw()
			},
			func() {
				// NOTE: I thought use that for debug hud
			},
		)
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

const FlameSize = 0.5

func initFlameParticleSystem(origin rl.Vector2) *ParticleSystem {
	ps := &ParticleSystem{}

	configFlame1 := EmitterConfig{
		Loop:         false,
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
