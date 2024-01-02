package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type MenuState struct {
	// @Incomplete: Handle user data here?

	MusicVolume float32
}

var menuState MenuState = MenuState{}

var choicesCount int
var currentMenuChoice int
var restartConfirmation bool
var quitConfirmation bool

func ToggleMenu() {
	if Mode == Game {
		Mode = Menu
	} else {
		Mode = Game
	}
}

func SimulateMenu() {
	handleChoiceInput := func(delta int) {
		currentMenuChoice += delta
		if currentMenuChoice < 0 {
			currentMenuChoice = 3
		}
		if currentMenuChoice > 3 {
			currentMenuChoice = 0
		}
	}

	if rl.IsKeyPressed(rl.KeyUp) {
		handleChoiceInput(-1)
	}

	if rl.IsKeyPressed(rl.KeyDown) {
		handleChoiceInput(1)
	}
}

func DrawMenu() {
	renderGradientMenuBackground()

	bigFont := Assets.FontsManager.BigFont

	const FontModifier = 1.4
	choicesCount = -1

	// @Cleanup: Create handy function to draw text on the center of the screen, also for measuring text dimensions

	//
	// Draw version
	//
	rl.DrawTextEx(*Assets.FontsManager.SmallFont, fmt.Sprintf("v%d.%d.%d", MajorVersion, MinorVersion, PatchVersion), rl.NewVector2(10, 10), SmallFontSize*0.8, 0, rl.RayWhite)

	var yy float32 = WindowHeight * 0.2

	//
	// Draw menu title
	//

	center := rl.MeasureTextEx(*bigFont, "Invaders", BigFontSize*0.7, 0)
	rl.DrawTextEx(*bigFont, "Invaders", rl.NewVector2(WindowWidth/2-center.X/2, yy), BigFontSize*0.7, 0, rl.RayWhite)
	yy += 128

	//
	// Draw resume
	//
	if drawItem(bigFont, "Resume", yy, SmallFontSize*FontModifier) {
		Mode = Game
	}
	yy += 48

	//
	// Draw restart
	//
	restartString := "Restart"
	if restartConfirmation {
		restartString = "Are you sure?"
	}

	if drawItem(bigFont, restartString, yy, SmallFontSize*FontModifier) {
		if !restartConfirmation {
			restartConfirmation = true
			quitConfirmation = false
		} else {
			Mode = Game
			restartConfirmation = false
		}
	}
	yy += 48

	//
	// Draw Music Option
	//
	musicString := "Music: On"
	drawItem(bigFont, musicString, yy, SmallFontSize*FontModifier)
	yy += 48

	//
	// Draw Exit
	//
	quitString := "Quit"
	if quitConfirmation {
		quitString = "Are you sure?"
	}

	if drawItem(bigFont, quitString, yy, SmallFontSize*FontModifier) {
		if !quitConfirmation {
			quitConfirmation = true
			restartConfirmation = false
		} else {
			ShouldClose = true
		}
	}
}

func drawItem(font *rl.Font, text string, yy float32, size float32) bool {
	choicesCount++
	center := rl.MeasureTextEx(*font, text, size, 0)

	// @Cleanup: Replace text shadow with shader stuff
	// rl.DrawTextEx(*font, text, rl.NewVector2((WindowWidth/2-center.X/2)+ShadowOffset, yy+ShadowOffset), size, 0, rl.Black)

	itemColor := rl.NewColor(156, 156, 156, 255)

	if choicesCount == currentMenuChoice {
		t := math.Cos(rl.GetTime() * 2)
		t *= t
		t = 0.4 + 0.55*t

		itemColor = linearColorFade(rl.Orange, rl.RayWhite, float32(t))
	}

	rl.DrawTextEx(*font, text, rl.NewVector2(WindowWidth/2-center.X/2, yy), size, 0, itemColor)

	return rl.IsKeyPressed(rl.KeyEnter) && choicesCount == currentMenuChoice
}

func renderGradientMenuBackground() {
	rl.ClearBackground(rl.Black)

	rl.DrawRectangleGradientV(
		0,
		0,
		WindowWidth,
		WindowHeight,
		rl.NewColor(c(.1), c(.1), c(.1), 255),
		rl.NewColor(c(.1), c(.1), c(.4), 255),
	)
}
