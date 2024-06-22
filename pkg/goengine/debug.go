package goengine

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jejikeh/invaders/pkg/goecs"
)

type DebufInfo struct {
	FPS float64
}

func spawnDebugInfo(layer *goecs.Layer) {
	debugInfo := layer.NewEntity()
	goecs.Attach[DebufInfo](layer, debugInfo)
}

func updateDebugInfo(layer *goecs.Layer) {
	for _, entity := range layer.Request(goecs.GetComponentID[DebufInfo](layer)) {
		info, _ := goecs.GetComponent[DebufInfo](layer, entity)

		info.FPS = ebiten.ActualFPS()
	}
}

func drawDebugInfo(screen *ebiten.Image, layer *goecs.Layer) {
	for _, entity := range layer.Request(goecs.GetComponentID[DebufInfo](layer)) {
		info, _ := goecs.GetComponent[DebufInfo](layer, entity)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", info.FPS))
	}
}
