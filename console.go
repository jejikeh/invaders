package main

import rl "github.com/gen2brain/raylib-go/raylib"

var openT float32 = 0.0
var openTarget float32 = 0.0
var openDelta float32 = 500

type ConsoleConfig struct {
	InputBackgroundColor rl.Color
	InputTextColor       rl.Color

	OutputBackgroundColor rl.Color
	OutputTextColor       rl.Color

	VerticalGradient bool
}

type Console struct {
	ConsoleConfig
	CurrentState ConsoleState
}

var GameConsole Console = Console{}

func (Console) Reload() {
	// @Incomplete
}

var inputTextColor rl.Color = rl.NewColor(255, 255, 255, 255)

var outputTextColor rl.Color = rl.NewColor(255, 255, 255, 255)

var consoleWindowX0 float32 = 0.0
var consoleWindowX1 float32 = 1.0
var consoleWindowYb float32 = 0.6

var BackBufferHeight float32 = 1.0

type ConsoleState int

const (
	Closed ConsoleState = iota
	OpenSmall
	OpenBig
)

func (c *Console) DrawConsole() {
	UpdateOpeness()

	x0 := consoleWindowX0 * BackBufferHeight
	x1 := GameDisplay.Width

	y0 := 0.0
	y1 := openT

	inputY0 := y0

	const Spacing = 48
	y0 += Spacing

	// cursorInOutput := false

	if c.VerticalGradient {
		rl.DrawRectangleGradientV(int32(x0), int32(inputY0), int32(x1), int32(y1), c.InputBackgroundColor, c.OutputBackgroundColor)
	} else {
		rl.DrawRectangleGradientH(int32(x0), int32(inputY0), int32(x1), int32(y1), c.InputBackgroundColor, c.OutputBackgroundColor)
	}
}

func Open(state ConsoleState) {
	GameConsole.CurrentState = state
	switch state {
	case OpenSmall:
		openTarget = float32(GameDisplay.Height) * 0.2
	case OpenBig:
		openTarget = float32(GameDisplay.Height) * 0.7
	case Closed:
		openTarget = 0
	}

	// to make it snappy
	UpdateOpeness()
}

func ToggleState() {
	GameConsole.CurrentState = (GameConsole.CurrentState + 1) % 3
	Open(GameConsole.CurrentState)
}

func UpdateOpeness() {
	dOpen := rl.GetFrameTime() * openDelta
	if openT < openTarget {
		openT += dOpen

		if openT > openTarget {
			openT = openTarget
		}
	} else if openT > openTarget {
		openT -= dOpen
		if openT < 0 {
			openT = 0
		}
	}
}

func GetConsoleBottom() float32 {
	return lerp(1.0, consoleWindowYb, openT) * BackBufferHeight
}

// @Incomplete: Normal lerp function.
func lerp(v1, v2 float32, fraction float32) float32 {
	return v1 + (v2-v1)*fraction
}
