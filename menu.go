package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const OptionColorFlashDuration = 1

// @Incomplete: This will be HandleEscape possible when modes will be incereased

func ToggleMenu() {
	if Mode == Game {
		Mode = Menu
	} else {
		Mode = Game
	}

	CurrentPage.Reset()
}

type PageElement interface {
	Draw()
	Simulate()
	Reset()
}

type Page struct {
	// @Cleanup: do we need this?
	ChoicesIterator   int
	CurrentMenuChoice int
	ChoicesCount      int

	ColorFlashTime   float32
	LastKeyPressTime float64
}

func (p *Page) Simulate() {
	handleChoiceInput := func(delta int) {
		p.CurrentMenuChoice += delta
		if p.CurrentMenuChoice < 0 {
			p.CurrentMenuChoice = p.ChoicesCount
		}
		if p.CurrentMenuChoice > p.ChoicesCount {
			p.CurrentMenuChoice = 0
		}
	}

	if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) {
		handleChoiceInput(-1)
		p.LastKeyPressTime = rl.GetTime()
	}

	if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) {
		handleChoiceInput(1)
		p.LastKeyPressTime = rl.GetTime()
	}

	// Change color of selected option
	// @Cleanup: Whe might want to not handle arrow keys for boolean options?
	if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyA) || rl.IsKeyPressed(rl.KeyD) || rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyRight) {
		p.ColorFlashTime = OptionColorFlashDuration
	}

	if p.ColorFlashTime > 0 {
		p.ColorFlashTime -= rl.GetFrameTime()
	}

	// @Incomplete: Add some lamda functions if need
}

func (p *Page) DrawBoolItem(font *rl.Font, text string, yy float32, size float32) bool {
	p.DrawItem(font, text, yy, size)

	// @Cleanup: Decide handle arrows here or not?
	return (rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyA) || rl.IsKeyPressed(rl.KeyD)) && p.ChoicesIterator == p.CurrentMenuChoice
}

func (p *Page) DrawIntItem(font *rl.Font, text string, yy float32, size float32) int {
	p.DrawItem(font, text, yy, size)

	dir := 0
	if p.ChoicesIterator == p.CurrentMenuChoice {
		if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD) {
			dir = 1
		}
		if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA) {
			dir = -1
		}
	}

	return dir
}

func (p *Page) DrawItem(font *rl.Font, text string, yy float32, size float32) {
	p.ChoicesIterator++
	center := rl.MeasureTextEx(*font, text, size, 0)

	// @Cleanup: Replace text shadow with shader stuff
	// rl.DrawTextEx(*font, text, rl.NewVector2((WindowWidth/2-center.X/2)+ShadowOffset, yy+ShadowOffset), size, 0, rl.Black)

	itemColor := rl.NewColor(156, 156, 156, 255)

	// Handle picked item
	if p.ChoicesIterator == p.CurrentMenuChoice {
		tBase := 0.2
		tRange := 0.6

		fromColor := rl.Orange
		toColor := rl.White

		t := math.Cos((rl.GetTime() - p.LastKeyPressTime) * 3)
		t *= t
		t = tBase + tRange*t

		// Color Flash selected option a little bit
		if p.ColorFlashTime > 0 {
			s := p.ColorFlashTime / OptionColorFlashDuration
			s *= s

			s = 1 - s

			fromColor = linearColorFade(rl.Black, fromColor, s)
			toColor = linearColorFade(rl.Black, toColor, s)
		}

		itemColor = linearColorFade(fromColor, toColor, float32(t))
	}

	rl.DrawTextEx(*font, text, rl.NewVector2(float32(GameDisplay.Width)/2-center.X/2, yy), size, 0, itemColor)
}

//
// MenuPage - for now it's just the pause
//

type MenuPage struct {
	Page
	RestartConfirmation bool
	QuitConfirmation    bool
}

func (p *MenuPage) Reset() {
	p.RestartConfirmation = false
	p.QuitConfirmation = false
}

func (p *MenuPage) Draw() {
	renderGradientMenuBackground := func() {
		rl.DrawRectangleGradientV(
			0,
			0,
			int32(GameDisplay.Width),
			int32(GameDisplay.Height),
			rl.NewColor(0, 0, 0, 240),
			rl.NewColor(0, 0, 0, 180),
		)
	}

	renderGradientMenuBackground()

	font := Assets.FontsManager.BigFont

	const FontModifier = 1.4
	p.ChoicesIterator = -1

	// @Cleanup: Create handy function to draw text on the center of the screen, also for measuring text dimensions

	//
	// Draw version
	//
	rl.DrawTextEx(*Assets.FontsManager.SmallFont, fmt.Sprintf("v%s.%s", MajorVersion, MinorVersion), rl.NewVector2(10, 10), SmallFontSize*0.8, 0, rl.RayWhite)

	var yy float32 = float32(GameDisplay.Height) * 0.2
	const Spacing = 48

	//
	// Draw menu title
	//
	center := rl.MeasureTextEx(*font, "Invaders", BigFontSize*0.7, 0)
	rl.DrawTextEx(*font, "Invaders", rl.NewVector2(float32(GameDisplay.Width)/2-center.X/2, yy), BigFontSize*0.7, 0, rl.RayWhite)
	yy += 128

	//
	// Draw resume
	//
	font = Assets.FontsManager.SmallFont

	if p.DrawBoolItem(font, "Resume", yy, SmallFontSize*FontModifier) {
		Mode = Game
		p.Reset()
	}
	yy += Spacing

	//
	// Draw restart
	//
	restartString := "Restart"
	if p.RestartConfirmation {
		restartString = "Are you sure?"
	}

	if p.DrawBoolItem(font, restartString, yy, SmallFontSize*FontModifier) {
		if !p.RestartConfirmation {
			p.RestartConfirmation = true
			p.QuitConfirmation = false
		} else {
			Mode = Game
			p.Reset()
		}
	}
	yy += Spacing

	//
	// Draw Options Option
	//
	//
	optionsString := "Options"
	if p.DrawBoolItem(font, optionsString, yy, SmallFontSize*FontModifier) {
		CurrentPage = &OptionsPage{}
	}

	yy += Spacing

	//
	// Draw Exit
	//
	quitString := "Quit"
	if p.QuitConfirmation {
		quitString = "Are you sure?"
	}

	if p.DrawBoolItem(font, quitString, yy, SmallFontSize*FontModifier) {
		if !p.QuitConfirmation {
			p.QuitConfirmation = true
			p.RestartConfirmation = false
		} else {
			ShouldClose = true

			// @Cleanup: Make new Update() function in saves.go
			Users.UpdateSettings()
			Users.SaveUser()
			p.Reset()
		}
	}

	// Set choices count equal to the number of items we iterate in this function
	// draw[*]Item each call will increment the choicesIterator
	p.ChoicesCount = p.ChoicesIterator
}

type OptionsPage struct {
	Page
}

func (p *OptionsPage) Reset() {
}

func (p *OptionsPage) Draw() {
	renderGradientMenuBackground := func() {
		rl.DrawRectangleGradientV(
			0,
			0,
			int32(GameDisplay.Width),
			int32(GameDisplay.Height),
			rl.NewColor(0, 0, 0, 240),
			rl.NewColor(0, 0, 0, 180),
		)
	}

	renderGradientMenuBackground()

	font := Assets.FontsManager.BigFont

	const FontModifier = 1.4
	p.ChoicesIterator = -1

	// @Cleanup: Create handy function to draw text on the center of the screen, also for measuring text dimensions

	//
	// Draw version
	//
	rl.DrawTextEx(*Assets.FontsManager.SmallFont, fmt.Sprintf("v%s.%s", MajorVersion, MinorVersion), rl.NewVector2(10, 10), SmallFontSize*0.8, 0, rl.RayWhite)

	var yy float32 = float32(GameDisplay.Height) * 0.2
	const Spacing = 48

	//
	// Draw menu title
	//
	center := rl.MeasureTextEx(*font, "Options", BigFontSize*0.7, 0)
	rl.DrawTextEx(*font, "Options", rl.NewVector2(float32(GameDisplay.Width)/2-center.X/2, yy), BigFontSize*0.7, 0, rl.RayWhite)
	yy += 128

	font = Assets.FontsManager.SmallFont

	//
	// Draw Music Option
	//
	musicString := fmt.Sprintf("Music: %.0f%%", UserVolume.Music*100)
	musicInputDirection := p.DrawIntItem(font, musicString, yy, SmallFontSize*FontModifier)
	switch musicInputDirection {
	case 1:
		// @Cleanup: Without support for maps, just make SetMusicVolume and etc...
		UserVolume.SetVolume(Music, UserVolume.Music+0.1)
	case -1:
		UserVolume.SetVolume(Music, UserVolume.Music-0.1)
	}

	yy += Spacing

	fullscreenString := "Fullscreen: On"
	if !GameDisplay.Fullscreen {
		fullscreenString = "Fullscreen: Off"
	}

	if p.DrawBoolItem(font, fullscreenString, yy, SmallFontSize*FontModifier) {
		if !GameDisplay.Fullscreen {
			GameDisplay.Fullscreen = true
		} else {
			GameDisplay.Fullscreen = false
		}

		// @Incomplete: Fullscreen mode
		// Renderer.UpdateFullscreen()
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

	if p.DrawBoolItem(font, debugString, yy, SmallFontSize*FontModifier) {
		if !Debug.Visible {
			Debug.Visible = true
		} else {
			Debug.Visible = false
		}
	}

	yy += Spacing

	// Draw Back
	//
	backString := "Back"
	if p.DrawBoolItem(font, backString, yy, SmallFontSize*FontModifier) {
		menuPage := &MenuPage{}
		// @Cleanup: Remove this hardcoded stuff?
		menuPage.CurrentMenuChoice = 2
		CurrentPage = menuPage
	}

	// Set choices count equal to the number of items we iterate in this function
	// draw[*]Item each call will increment the choicesIterator
	p.ChoicesCount = p.ChoicesIterator
}
