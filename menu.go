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

var menuState MenuState = MenuState{
	MusicVolume: 1,
}

var choicesIterator int
var currentMenuChoice int
var choicesCount int

var restartConfirmation bool
var quitConfirmation bool

func ToggleMenu() {
	if Mode == Game {
		Mode = Menu
	} else {
		Mode = Game
	}

	quitConfirmation = false
	restartConfirmation = false
}

func SimulateMenu() {
	handleChoiceInput := func(delta int) {
		currentMenuChoice += delta
		if currentMenuChoice < 0 {
			currentMenuChoice = choicesCount
		}
		if currentMenuChoice > choicesCount {
			currentMenuChoice = 0
		}
	}

	if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) {
		handleChoiceInput(-1)
	}

	if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) {
		handleChoiceInput(1)
	}
}

func DrawMenu() {
	renderGradientMenuBackground()

	bigFont := Assets.FontsManager.BigFont

	const FontModifier = 1.4
	choicesIterator = -1

	// @Cleanup: Create handy function to draw text on the center of the screen, also for measuring text dimensions

	//
	// Draw version
	//
	rl.DrawTextEx(*Assets.FontsManager.SmallFont, fmt.Sprintf("v%s.%s", MajorVersion, MinorVersion), rl.NewVector2(10, 10), SmallFontSize*0.8, 0, rl.RayWhite)

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
			quitConfirmation = false
		}
	}
	yy += 48

	//
	// Draw Music Option
	//
	musicString := "Music: On"
	if menuState.MusicVolume == 0 {
		musicString = "Music: Off"
	}
	if drawItem(bigFont, musicString, yy, SmallFontSize*FontModifier) {
		if menuState.MusicVolume == 0 {
			// @Cleanup: Create music manager...
			rl.SetMusicVolume(music, 0.4)
			menuState.MusicVolume = 1
		} else {
			// @Cleanup: Create music manager...
			rl.SetMusicVolume(music, 0.0)
			menuState.MusicVolume = 0
		}
	}
	yy += 48

	//
	// Draw Debug Option
	//
	//

	debugString := "Debug: On"
	if !Debug.Visible {
		debugString = "Debug: Off"
	}

	if drawItem(bigFont, debugString, yy, SmallFontSize*FontModifier) {
		if !Debug.Visible {
			Debug.Visible = true
		} else {
			Debug.Visible = false
		}
	}

	yy += 48

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
			quitConfirmation = false
			restartConfirmation = false
		}
	}

	// Set choices count equal to the number of items we iterate in this function
	// drawItem each call will increment the choicesIterator
	choicesCount = choicesIterator
}

func drawItem(font *rl.Font, text string, yy float32, size float32) bool {
	choicesIterator++
	center := rl.MeasureTextEx(*font, text, size, 0)

	// @Cleanup: Replace text shadow with shader stuff
	// rl.DrawTextEx(*font, text, rl.NewVector2((WindowWidth/2-center.X/2)+ShadowOffset, yy+ShadowOffset), size, 0, rl.Black)

	itemColor := rl.NewColor(156, 156, 156, 255)

	if choicesIterator == currentMenuChoice {
		t := math.Cos(rl.GetTime() * 2)
		t *= t
		t = 0.2 + 0.5*t

		itemColor = linearColorFade(rl.White, rl.Orange, float32(t))
	}

	rl.DrawTextEx(*font, text, rl.NewVector2(WindowWidth/2-center.X/2, yy), size, 0, itemColor)

	return (rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace)) && choicesIterator == currentMenuChoice
}

func renderGradientMenuBackground() {
	rl.DrawRectangleGradientV(
		0,
		0,
		WindowWidth,
		WindowHeight,
		rl.NewColor(0, 0, 0, 200),
		rl.NewColor(0, 0, 0, 100),
	)
}
