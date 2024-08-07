package main

import (
	"embed"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jejikeh/invaders/pkg/engine/assets"
	"github.com/jejikeh/invaders/pkg/engine/log"
)

//go:embed resources/*
var resources embed.FS

var assetImportScheme = []assets.Importer{
	&assets.AtlasExport{
		Name: "ships",
		Path: "resources/tilemaps/ships.png",
		SpriteCount: rl.Vector2{
			X: 4,
			Y: 6,
		},
	},
	&assets.AtlasExport{
		Name: "tiles",
		Path: "resources/tilemaps/tiles.png",
		SpriteCount: rl.Vector2{
			X: 10,
			Y: 12,
		},
	},
}

// @Incomplete: Move to game struct or ecs system.
var invadersLog = log.New(log.LogWriter, "[invaders]")

func main() {
	// @Incomplete: Handle window config flags here. Maybe, at some point remove
	const (
		width  = 1024
		height = 768
	)

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

	assets, err := assets.Import(resources, assetImportScheme...)
	if err != nil {
		invadersLog.Fatal(err)
	}
	defer assets.Unload()

	gameRenderTexture := NewRenderTexture(width, height, nil)
	defer gameRenderTexture.Unload()

	bgRenderTexture := NewRenderTexture(width, height, nil)
	defer bgRenderTexture.Unload()

	shipsRenderTexture := NewRenderTexture(width, height, assets.Atlases.Get("ships"))
	defer shipsRenderTexture.Unload()

	playerPos := rl.Vector3{X: 4, Y: 7, Z: 4}
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
		camera.Target.Z = rl.Lerp(camera.Target.Z, playerPos.Z-13, 2*dt)

		rl.UpdateCameraPro(&camera, cameraPos, rl.Vector3{}, 0)

		// rl.BeginShaderMode(*assets.Shaders.Get("grayscale"))

		bgRenderTexture.Func2D(camera, func() {
			rl.DrawRectangleGradientV(
				0,
				0,
				int32(1024),
				int32(768),
				rl.NewColor(c(.1), c(.1), c(.4), 255),
				rl.NewColor(c(.2), c(.5), c(.7), 255),
			)
		})

		shipsRenderTexture.Func3D(camera, func() {
			for z := range 6 {
				shipsRenderTexture.QuadFunc(
					rl.Vector3{Y: 3 + float32(z)/2, Z: float32(z)},
					rl.Vector3{X: 2, Y: 1, Z: 2},
					uint32(z),
				)
			}

			shipsRenderTexture.QuadFunc(
				playerPos,
				rl.Vector3{X: 2, Y: 1, Z: 2},
				10,
				func(v rl.Vector3) rl.Vector3 {
					return rl.Vector3RotateByAxisAngle(v, rl.Vector3{X: playerPos.X, Y: playerPos.Y, Z: playerPos.Z}, playerRot)
				},
			)
		})

		// rl.EndShaderMode()

		gameRenderTexture.BlendTextures(bgRenderTexture, shipsRenderTexture)

		gameRenderTexture.Draw()
	}
}

func c(v float32) uint8 {
	return uint8(v * 255)
}
