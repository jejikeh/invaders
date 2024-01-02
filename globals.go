package main

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
