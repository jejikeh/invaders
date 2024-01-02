package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var choicesIterator int
var currentMenuChoice int
var choicesCount int

var restartConfirmation bool
var quitConfirmation bool

var optionColorFlashTime float32
var lastKeyPressTime float64

const OptionColorFlashDuration = 1

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
		lastKeyPressTime = rl.GetTime()
	}

	if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) {
		handleChoiceInput(1)
		lastKeyPressTime = rl.GetTime()
	}

	// Change color of selected option
	if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyA) || rl.IsKeyPressed(rl.KeyD) || rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyRight) {
		optionColorFlashTime = OptionColorFlashDuration
	}

	if optionColorFlashTime > 0 {
		optionColorFlashTime -= rl.GetFrameTime()
	}
}

func DrawMenu() {
	renderGradientMenuBackground()

	font := Assets.FontsManager.BigFont

	const FontModifier = 1.4
	choicesIterator = -1

	// @Cleanup: Create handy function to draw text on the center of the screen, also for measuring text dimensions

	//
	// Draw version
	//
	rl.DrawTextEx(*Assets.FontsManager.SmallFont, fmt.Sprintf("v%s.%s", MajorVersion, MinorVersion), rl.NewVector2(10, 10), SmallFontSize*0.8, 0, rl.RayWhite)

	var yy float32 = WindowHeight * 0.2
	const Spacing = 48

	//
	// Draw menu title
	//
	center := rl.MeasureTextEx(*font, "Invaders", BigFontSize*0.7, 0)
	rl.DrawTextEx(*font, "Invaders", rl.NewVector2(WindowWidth/2-center.X/2, yy), BigFontSize*0.7, 0, rl.RayWhite)
	yy += 128

	//
	// Draw resume
	//

	font = Assets.FontsManager.SmallFont

	if drawBoolItem(font, "Resume", yy, SmallFontSize*FontModifier) {
		Mode = Game

		quitConfirmation = false
		restartConfirmation = false
	}
	yy += Spacing

	//
	// Draw restart
	//
	restartString := "Restart"
	if restartConfirmation {
		restartString = "Are you sure?"
	}

	if drawBoolItem(font, restartString, yy, SmallFontSize*FontModifier) {
		if !restartConfirmation {
			restartConfirmation = true
			quitConfirmation = false
		} else {
			Mode = Game

			restartConfirmation = false
			quitConfirmation = false
		}
	}
	yy += Spacing

	//
	// Draw Music Option
	//
	musicString := fmt.Sprintf("Music: %.0f%%", Volumes[Music]*100)
	musicInputDirection := drawIntItem(font, musicString, yy, SmallFontSize*FontModifier)
	switch musicInputDirection {
	case 1:
		SetVolume(Music, Volumes[Music]+0.1)
	case -1:
		SetVolume(Music, Volumes[Music]-0.1)
	}

	yy += Spacing

	//
	// Draw Debug Option
	//
	//

	debugString := "Debug: On"
	if !Debug.Visible {
		debugString = "Debug: Off"
	}

	if drawBoolItem(font, debugString, yy, SmallFontSize*FontModifier) {
		if !Debug.Visible {
			Debug.Visible = true
		} else {
			Debug.Visible = false
		}
	}

	yy += Spacing

	// Draw Exit
	//
	quitString := "Quit"
	if quitConfirmation {
		quitString = "Are you sure?"
	}

	if drawBoolItem(font, quitString, yy, SmallFontSize*FontModifier) {
		if !quitConfirmation {
			quitConfirmation = true
			restartConfirmation = false
		} else {
			ShouldClose = true

			// @Cleanup: Make new Update() function in saves.go
			Users.UpdateSettings()
			Users.SaveUser()

			quitConfirmation = false
			restartConfirmation = false
		}
	}

	// Set choices count equal to the number of items we iterate in this function
	// draw[*]Item each call will increment the choicesIterator
	choicesCount = choicesIterator
}

func drawBoolItem(font *rl.Font, text string, yy float32, size float32) bool {
	choicesIterator++
	center := rl.MeasureTextEx(*font, text, size, 0)

	// @Cleanup: Replace text shadow with shader stuff
	// rl.DrawTextEx(*font, text, rl.NewVector2((WindowWidth/2-center.X/2)+ShadowOffset, yy+ShadowOffset), size, 0, rl.Black)

	itemColor := rl.NewColor(156, 156, 156, 255)

	// Handle picked item
	if choicesIterator == currentMenuChoice {
		tBase := 0.2
		tRange := 0.6

		fromColor := rl.Orange
		toColor := rl.White

		t := math.Cos((rl.GetTime() - lastKeyPressTime) * 3)
		t *= t
		t = tBase + tRange*t

		// Color Flash selected option a little bit
		if optionColorFlashTime > 0 {
			s := optionColorFlashTime / OptionColorFlashDuration
			s *= s

			s = 1 - s

			fromColor = linearColorFade(rl.Black, fromColor, s)
			toColor = linearColorFade(rl.Black, toColor, s)
		}

		itemColor = linearColorFade(fromColor, toColor, float32(t))
	}

	rl.DrawTextEx(*font, text, rl.NewVector2(WindowWidth/2-center.X/2, yy), size, 0, itemColor)

	return (rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace)) && choicesIterator == currentMenuChoice
}

func drawIntItem(font *rl.Font, text string, yy float32, size float32) int {
	choicesIterator++
	center := rl.MeasureTextEx(*font, text, size, 0)

	// @Cleanup: Replace text shadow with shader stuff
	// rl.DrawTextEx(*font, text, rl.NewVector2((WindowWidth/2-center.X/2)+ShadowOffset, yy+ShadowOffset), size, 0, rl.Black)

	itemColor := rl.NewColor(156, 156, 156, 255)

	// Handle picked item
	if choicesIterator == currentMenuChoice {
		tBase := 0.2
		tRange := 0.6

		fromColor := rl.Orange
		toColor := rl.White

		t := math.Cos((rl.GetTime() - lastKeyPressTime) * 3)
		t *= t
		t = tBase + tRange*t

		// Color Flash selected option a little bit
		if optionColorFlashTime > 0 {
			s := optionColorFlashTime / OptionColorFlashDuration
			s *= s

			s = 1 - s

			fromColor = linearColorFade(rl.Black, fromColor, s)
			toColor = linearColorFade(rl.Black, toColor, s)
		}

		itemColor = linearColorFade(fromColor, toColor, float32(t))
	}

	rl.DrawTextEx(*font, text, rl.NewVector2(WindowWidth/2-center.X/2, yy), size, 0, itemColor)

	dir := 0
	if choicesIterator == currentMenuChoice {
		if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD) {
			dir = 1
		}
		if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA) {
			dir = -1
		}
	}

	return dir
}

func renderGradientMenuBackground() {
	rl.DrawRectangleGradientV(
		0,
		0,
		WindowWidth,
		WindowHeight,
		rl.NewColor(0, 0, 0, 240),
		rl.NewColor(0, 0, 0, 180),
	)
}
