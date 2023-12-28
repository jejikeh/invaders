package invaders

import "core:fmt"
import "vendor:raylib"
import "core:strings"

// https://www.youtube.com/watch?v=As5XL0C4rR4&t=978s
ASPECT :: (224.0 / 288.0)
VPIXELS :: 720

// TODO: Make a normal struct with window data
WINDOW_WIDTH :: (VPIXELS * ASPECT)
WINDOW_HEIGHT :: VPIXELS

WINDOW_MINIMAL_DEL :: 1

main :: proc() {
	display_text: cstring = "Hello, Raylib!"

	// config, err_init_window := init_window()
	// if err_init_window != nil {
	// 	fmt.eprint(err_init_window)
	// 	return
	// }

	raylib.SetConfigFlags({.WINDOW_RESIZABLE, .VSYNC_HINT})
	raylib.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "[odin] Jai Invaders")
	raylib.SetWindowMinSize(WINDOW_WIDTH / WINDOW_MINIMAL_DEL, WINDOW_HEIGHT / WINDOW_MINIMAL_DEL)

	renderTarget := raylib.LoadRenderTexture(WINDOW_WIDTH, WINDOW_HEIGHT)
	raylib.SetTextureFilter(renderTarget.texture, .POINT)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
        screen_start_width := (f32(
            f32(raylib.GetScreenWidth()) -
            (f32(WINDOW_WIDTH)) *
                (f32(raylib.GetScreenHeight()) / f32(WINDOW_HEIGHT))
        ));

		if raylib.IsKeyPressed(.ESCAPE) {
			raylib.CloseWindow()
		}

		// Texture Render Mode
		raylib.BeginTextureMode(renderTarget)

		raylib.ClearBackground(raylib.RAYWHITE)
		raylib.DrawText(
			fmt.caprintf("hello, world"),
			WINDOW_WIDTH / 2 - 50,
            WINDOW_HEIGHT / 2,
			20,
			raylib.DARKGRAY,
		)

		raylib.EndTextureMode()


		// Window Render Mode
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.BLACK)
		raylib.DrawTexturePro(
			renderTarget.texture,
			raylib.Rectangle {
				0,
				0,
				f32(renderTarget.texture.width),
				-f32(renderTarget.texture.height),
			},
			raylib.Rectangle {
				(f32(
						f32(raylib.GetScreenWidth()) -
						(f32(WINDOW_WIDTH)) *
							(f32(raylib.GetScreenHeight()) / f32(WINDOW_HEIGHT))
					) /
					2),
				0,
				(f32(WINDOW_WIDTH) *
                (f32(raylib.GetScreenHeight()) / f32(WINDOW_HEIGHT))),
				f32(raylib.GetScreenHeight()),
			},
			raylib.Vector2{0, 0},
			0,
			raylib.WHITE,
		)

		raylib.EndDrawing()
	}

	raylib.UnloadRenderTexture(renderTarget)
	raylib.CloseWindow()
}
