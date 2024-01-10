package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ConsoleConfig struct {
	InputBackgroundColor rl.Color
	InputTextColor       rl.Color

	OutputBackgroundColor rl.Color
	OutputTextColor       rl.Color

	VerticalGradient bool

	ScreenPercentX  float32
	ScreenPercentX1 float32

	ScreenPercentY float32

	ScreenPercentYSmall  float32
	ScreenPercentYBig    float32
	ScreenPercentYClosed float32

	Speed float32

	StopTheGame bool
}

type Console struct {
	ConsoleConfig
	CurrentState ConsoleState

	X float32
	Y float32

	X1 float32
	Y1 float32
}

var GameConsole Console = Console{
	X: 0,
	ConsoleConfig: ConsoleConfig{
		Speed: 500,
	},
}

func (c *Console) Init() {
	c.X = float32(GameDisplay.Width) * c.ScreenPercentX
	c.X1 = (float32(GameDisplay.Width) - c.X) * c.ScreenPercentX1

	// @Incomplete: Handle c.Y in GetTargetHeight?
	c.Y = float32(GameDisplay.Height) * c.ScreenPercentY
	c.Y1 = c.GetTargetHeight()
}

func (g *Console) GetTargetHeight() float32 {
	switch g.CurrentState {
	case OpenSmall:
		return float32(GameDisplay.Height) * g.ScreenPercentYSmall
	case OpenBig:
		return float32(GameDisplay.Height) * g.ScreenPercentYBig
	case Closed:
		return float32(GameDisplay.Height) * g.ScreenPercentYClosed
	}

	return 0
}

func (Console) Reload() {
	// @Incomplete
}

type ConsoleState int

const (
	Closed ConsoleState = iota
	OpenSmall
	OpenBig
)

func (c *Console) Draw() {
	c.UpdateOpeness()

	if c.VerticalGradient {
		rl.DrawRectangleGradientV(int32(c.X), int32(c.Y), int32(c.X1), int32(c.Y1), c.InputBackgroundColor, c.OutputBackgroundColor)
	} else {
		rl.DrawRectangleGradientH(int32(c.X), int32(c.Y), int32(c.X1), int32(c.Y1), c.InputBackgroundColor, c.OutputBackgroundColor)
	}
}

func (c *Console) SetState(state ConsoleState) {
	GameConsole.CurrentState = state
	c.UpdateOpeness()
}

func (c *Console) ToggleState() {
	GameConsole.CurrentState = (GameConsole.CurrentState + 1) % 3
	c.SetState(GameConsole.CurrentState)
}

func (c *Console) UpdateOpeness() {
	dOpen := rl.GetFrameTime() * c.Speed
	if c.Y1 < c.GetTargetHeight() {
		c.Y1 += dOpen

		if c.Y1 > c.GetTargetHeight() {
			c.Y1 = c.GetTargetHeight()
		}
	} else if c.Y1 > c.GetTargetHeight() {
		c.Y1 -= dOpen
		if c.Y1 < 0 {
			c.Y1 = 0
		}
	}
}

// @Incomplete: Normal lerp function.
func lerp(v1, v2 float32, fraction float32) float32 {
	return v1 + (v2-v1)*fraction
}
