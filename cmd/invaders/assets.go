package main

import rl "github.com/gen2brain/raylib-go/raylib"

type AssetsManager struct {
	TexturesManager *TexturesManager
	FontsManager    *FontsManager
}

func NewAssetsManager() *AssetsManager {
	return &AssetsManager{
		TexturesManager: NewTexturesManager(),
		FontsManager:    NewFontsManager(),
	}
}

func (am *AssetsManager) Unload() {
	am.TexturesManager.Unload()
	am.FontsManager.Unload()
}

type FontsManager struct {
	BigFont   *rl.Font
	SmallFont *rl.Font
}

func NewFontsManager() *FontsManager {
	bigFont := rl.LoadFontEx(FontFolder+"Martel-Regular.ttf", BigFontSize, nil)
	smallFont := rl.LoadFontEx(FontFolder+"Martel-Regular.ttf", SmallFontSize, nil)

	return &FontsManager{
		BigFont:   &bigFont,
		SmallFont: &smallFont,
	}
}

func (fm *FontsManager) Unload() {
	rl.UnloadFont(*fm.BigFont)
	rl.UnloadFont(*fm.SmallFont)
}

type TexturesManager struct {
	Player *rl.Texture2D
	Dude   *rl.Texture2D
	Alien  *rl.Texture2D
	Rocket *rl.Texture2D
	Circle *rl.Texture2D
}

func NewTexturesManager() *TexturesManager {
	player := rl.LoadTexture(ResourseFolder + "player.png")
	dude := rl.LoadTexture(ResourseFolder + "dude.png")
	alien := rl.LoadTexture(ResourseFolder + "alien.png")
	rocket := rl.LoadTexture(ResourseFolder + "rocket.png")

	imageCircle := rl.GenImageGradientRadial(16, 16, 0.3, rl.White, rl.Black)
	textureCircle := rl.LoadTextureFromImage(imageCircle)

	return &TexturesManager{
		Player: &player,
		Dude:   &dude,
		Alien:  &alien,
		Rocket: &rocket,
		Circle: &textureCircle,
	}
}

func (tm *TexturesManager) Unload() {
	rl.UnloadTexture(*tm.Player)
	rl.UnloadTexture(*tm.Dude)
	rl.UnloadTexture(*tm.Alien)
	rl.UnloadTexture(*tm.Rocket)
	rl.UnloadTexture(*tm.Circle)
}
