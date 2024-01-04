package main

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	MajorVersion string
	MinorVersion string
)

const Aspect = 4 / 3
const VerticalPixels = 720

const WindowWidth = Aspect * VerticalPixels
const WindowHeight = VerticalPixels
const WindowMinimalSizeDelimeter = 1

const ResourseFolder = "resources/"
const FontFolder = ResourseFolder + "fonts/"

const AudioFolder = ResourseFolder + "audio/"
const SoundFolder = AudioFolder + "sound/"
const MusicFolder = AudioFolder + "music/"

const SavesFolder = ResourseFolder + "saves/"

const EntitiesBaseSize = 0.3

const BigFontSize = 192
const SmallFontSize = 38

var ShouldClose bool

var Assets *AssetsManager
var Renderer *Render
var Entities *EntityManager
var Emitters *EmitterManager
var AudioManager *Mixer
var Users *UserManager

var Debug *Hud

var Mode GameMode = Game

var CurrentPage PageElement

func main() {
	// @Refactor: Create global hud manager or something like that
	// And handle there menu, editor, game hub maybe??
	// @Cleanup: Create global state manager what will be manages the game state (Menu, Editor, Game)
	// @Cleanup: Make possible to reset any game state to initial state

	InitVariables(ResourseFolder + "invaders.variables")
	return
	Debug = NewDebugHud()

	CurrentPage = &MenuPage{}

	Renderer = NewRender()
	defer Renderer.Unload()

	gui.LoadStyle("resources/cherry/style_cherry.rgs")
	rl.SetExitKey(rl.KeyMinus)

	Assets = NewAssetsManager()
	defer Assets.Unload()

	AudioManager = NewMixer()
	defer AudioManager.Unload()

	// @Cleanup: New user manager is such a mess.
	Users = NewUserManager()
	Users.AddUser()

	player := NewPlayer()
	Entities = NewEntityManager()
	Entities.Add(player)

	Emitters = NewEmitterManager()
	flameParticle := initFlameParticleSystem(rl.NewVector2(player.Position.X, player.Position.Y))

	Emitters.AddHandlers(flameParticle,
		func(emitter *ParticleSystem) {
			emitter.ReStart()
		},
		func(emitter *ParticleSystem) {
			rect := player.GetRectangle()
			emitter.SetOrigin(rl.NewVector2((rect.X+rect.Width/2)-2, rect.Y+rect.Height))
		},
		func(emitter *ParticleSystem) {},
	)

	Emitters.Add(initFlameParticleSystem(rl.NewVector2(player.Position.X, player.Position.Y)))

	Entities.Start()
	Emitters.Start()

	// @Cleanup: Maybe instead of hardcoded music file names we can get it from somewhere else
	AudioManager.SetMusicLoop("music", true)
	AudioManager.PlayMusic("music")

	// @Hack: for some freaking reason IsKeyPressed invokes two times...
	wasPressedPrevFrame := false

	for !rl.WindowShouldClose() && !ShouldClose {
		// @Hack: for some freaking reason IsKeyPressed invokes two times...
		{
			if rl.IsKeyPressed(rl.KeyEscape) && !wasPressedPrevFrame {
				ToggleMenu()
				wasPressedPrevFrame = true
			}

			if rl.IsKeyUp(rl.KeyEscape) {
				wasPressedPrevFrame = false
			}
		}

		AudioManager.UpdateMusic()

		if Mode == Game {
			SimulateInvaders()
		} else if Mode == Menu {
			CurrentPage.Simulate()
		}

		Renderer.Draw(
			func() {
				if Mode == Game {
					DrawInvaders()
					Debug.Draw()
				} else if Mode == Menu {
					// So the game will be visivle in the menu
					DrawInvaders()
					CurrentPage.Draw()
				}
			},
			func() {
				// NOTE: I thought use that for debug hud
				// @Note: This is rendered not in the Render Texture, but in window viewport
				// @Note: Maybe, in the future all game will be resized to window size?
			},
		)
	}
}

func SimulateInvaders() {
	Entities.Update()
	Emitters.Update()

	// @Note: Maybe change from debug to editor mode?
	Debug.Update()
}

func DrawInvaders() {
	renderGradientBackground()

	Entities.FirstDraw()

	Emitters.Draw()
	Entities.Draw()
}

// @Cleanup: Move this to player scope
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

// @Cleanup: Move this to player scope
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
		rl.NewColor(c(.1), c(.1), c(.2), 0),
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
