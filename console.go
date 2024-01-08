package main

import rl "github.com/gen2brain/raylib-go/raylib"

var openT float32 = 0.0
var openTarget float32 = 0.0
var openDelta float32 = 100

var inputBackgroundColor rl.Color = rl.NewColor(0, 0, 0, 255)
var inputTextColor rl.Color = rl.NewColor(255, 255, 255, 255)

var outputBackgroundColor rl.Color = rl.NewColor(0, 0, 0, 255)
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

func DrawConsole() {
	UpdateOpeness()

	x0 := consoleWindowX0 * BackBufferHeight
	x1 := GameDisplay.Width

	y0 := 0.0
	y1 := openT

	inputY0 := y0

	const Spacing = 48
	y0 += Spacing

	// cursorInOutput := false

	rl.DrawRectangleGradientH(int32(x0), int32(inputY0), int32(x1), int32(y1), inputBackgroundColor, outputBackgroundColor)
}

func Open(state ConsoleState) {
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
