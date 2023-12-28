package invaders

import "core:fmt"
import "vendor:raylib"
import "core:strings"
import "core:math"

// https://www.youtube.com/watch?v=As5XL0C4rR4&t=978s
// https://www.youtube.com/watch?v=AwbLjR46HkY
ASPECT :: (224.0 / 288.0)
VPIXELS :: 720
WINDOW_WIDTH :: (VPIXELS * ASPECT)
WINDOW_HEIGHT :: VPIXELS
WINDOW_MINIMAL_DEL :: 1

RESOURCES_PATH :: "./resources/"

main_font: raylib.Font 
render_texture: raylib.RenderTexture2D
textures: ^InvadersTextures

InvadersTextures :: struct {
	player: raylib.Texture2D,
}

Entity :: struct {
	texture: raylib.Texture2D,
	pos: raylib.Vector2,
	size: raylib.Vector2,
	rotation: f32,
	tint: raylib.Color,
	shadow: bool,
	height: f32
}

make_entity :: proc ( texture: raylib.Texture2D, pos: raylib.Vector2, 
	size:= raylib.Vector2{1, 1}, rotation:f32= .0, tint:=raylib.RAYWHITE, 
	shadow:=true, height:f32=5) -> (entity: Entity) {
	entity.texture = texture
	entity.pos = pos
	entity.size = size
	entity.rotation = rotation
	entity.tint = tint
	entity.shadow = shadow
	entity.height = height

	return entity
}

main :: proc() {
	raylib.SetConfigFlags({.WINDOW_RESIZABLE, .VSYNC_HINT})
	raylib.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "[odin] Jai Invaders")
	raylib.SetWindowMinSize(WINDOW_WIDTH / WINDOW_MINIMAL_DEL, WINDOW_HEIGHT / WINDOW_MINIMAL_DEL)

	render_texture := raylib.LoadRenderTexture(WINDOW_WIDTH, WINDOW_HEIGHT)
	raylib.SetTextureFilter(render_texture.texture, .POINT)

    main_font = raylib.LoadFont(RESOURCES_PATH + "robtronika.ttf");

	textures = new(InvadersTextures)
	invaders_load_textures(textures)

	player := make_entity(textures.player, raylib.Vector2{
		(WINDOW_WIDTH / 2),
		(WINDOW_HEIGHT / 2),
	},
	raylib.Vector2{
		.5,
		.5,
	},
	0.0,
	raylib.WHITE)

	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		if raylib.IsKeyPressed(.ESCAPE) {
			raylib.CloseWindow()
		}

		// Texture Render Mode
		raylib.BeginTextureMode(render_texture)
		raylib.ClearBackground(raylib.BLACK)

		render_background()

		player.rotation += 1
		render_entity(player)
		
		raylib.EndTextureMode()

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.BLACK)
		raylib.DrawTexturePro(
			render_texture.texture,
			raylib.Rectangle {
				0,
				0,
				f32(render_texture.texture.width),
				-f32(render_texture.texture.height),
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

	invaders_unload()
}

invaders_load_textures :: proc(invaders_textures: ^InvadersTextures) {
	invaders_textures.player = raylib.LoadTexture(RESOURCES_PATH + "player.png")
}

invaders_unload_textures :: proc(invaders_textures: ^InvadersTextures) {
	raylib.UnloadTexture(invaders_textures.player)
}

invaders_unload :: proc() {
	invaders_unload_textures(textures)

	raylib.UnloadRenderTexture(render_texture)
    raylib.UnloadFont(main_font)
	raylib.CloseWindow()
}

render_background :: proc() {
	raylib.DrawRectangleGradientV(
		0, 
		WINDOW_HEIGHT/1.5, 
		WINDOW_WIDTH, 
		WINDOW_HEIGHT / 3, 
		color(.1, .1, .9, 1), 
		color(.2, .5, .7, 1))

	raylib.DrawRectangleGradientV(
		0, 
		0, 
		WINDOW_WIDTH, 
		WINDOW_HEIGHT/1.5, 
		color(.1, .1, .2, 1), 
		color(.1, .1, .9, 1))
}

color :: proc(r: f32, g: f32, b: f32, a: f32) -> raylib.Color {
	return raylib.Color{
		cast(u8)(255.0 * r),
		cast(u8)(255.0 * g), 
		cast(u8)(255.0 * b),
		cast(u8)(255.0 * a),
	}
}

render_texture_centered :: proc(texture: raylib.Texture2D, pos: raylib.Vector2, size: raylib.Vector2, rotation: f32, color: raylib.Color) {
	// TODO: format it
	raylib.DrawTexturePro(texture, raylib.Rectangle {
		0,
		0,
		f32(texture.width),
		-f32(texture.height),
	}, raylib.Rectangle{pos.x, pos.y, f32(texture.width) * size.x, f32(texture.height) * size.y}, raylib.Vector2 { f32(texture.width / 2) * size.x, f32(texture.height / 2) * size.y}, rotation, color)
}

render_entity :: proc(entity: Entity) {
	if entity.shadow {
		render_texture_centered(entity.texture, raylib.Vector2{entity.pos.x + entity.height, entity.pos.y - entity.height}, entity.size, entity.rotation, raylib.BLACK)
	}

	render_texture_centered(entity.texture, entity.pos, entity.size, entity.rotation, entity.tint)
}