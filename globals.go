package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// @Todo: Maybe it`s a dumb way to have a separate global values for layers.
// But, maybe it wont be some dumb if we have all globals here? ant it will be
// called globals.go

const (
	PlayerLayer  = 10
	InvaderLayer = 10

	EmitterLayer = 9

	ShadowLayer = 8
)

type GameMode int

const (
	Game GameMode = iota
	Menu
	// DebugMode
)

type Display struct {
	VerticalPixels int
	Width          int
	Height         int
	VSync          bool
	MaxFPS         int
	Title          string
	Fullscreen     bool
	HiDPI          bool
}

// Some values will be applied from the start, but some will be changed in the first tick of game loop
var GameDisplay Display = Display{
	VerticalPixels: 720,
	VSync:          true,
	MaxFPS:         60,
	Title:          "Invaders",
	Fullscreen:     false,
	HiDPI:          false,
}

func InitGameDisplay() {
	GameDisplay.Height = GameDisplay.VerticalPixels
	GameDisplay.Width = GameDisplay.Height * Aspect
}

func (Display) Reload() {
	// @Incomplete: Handle change of resolution, fullscreen and vsync here

	// Called on the start while render.go may not be initialized yet
	if !rl.IsWindowReady() {
		return
	}

	// Check if in main thread
	GameDisplay.Height = GameDisplay.VerticalPixels
	GameDisplay.Width = GameDisplay.Height * Aspect

	// This need to be called in main thread.
	// @Cleanup: Or move vars.go assigment in main thread. Idk which solution will be less ugly
	ShouldUpdateWindowInMainThread = true

	rl.SetTargetFPS(int32(GameDisplay.MaxFPS))
	SetWindowFlags()
}

var GameVolume Volume

func InitGameVolume() {
	GameVolume = Volume{
		All:      1.0,
		Props:    1.0,
		Ambience: 1.0,
		Movement: 1.0,
		UI:       1.0,
		Music:    1.0,
	}
}

func (Volume) Reload() {
	// Called on the start while AudioManager may not be initialized yet
	if AudioManager == nil {
		return
	}

	AudioManager.UpdateVolume()
}
