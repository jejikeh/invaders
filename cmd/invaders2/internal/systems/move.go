package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jejikeh/invaders/cmd/invaders2/internal/components"
	"github.com/jejikeh/invaders/pkg/goecs"
	"github.com/jejikeh/invaders/pkg/goengine"
	"github.com/jejikeh/invaders/pkg/gomath"
)

var keys = []ebiten.Key{}

func Move(layer *goecs.Layer) {
	input := inputVector()

	for _, entity := range layer.Request(goecs.GetComponentID[components.Player](layer)) {
		transform, _ := goecs.GetComponent[goengine.Transfrom](layer, entity)
		movable, _ := goecs.GetComponent[components.Movable](layer, entity)

		if input.X == 0 && input.Y == 0 {
			movable.Velocity = movable.Velocity.Sub(
				movable.Velocity.Mul(movable.Friction *),
			)
		}
	}
}

func inputVector() gomath.Vec2 {
	input := gomath.NewVec2(0, 0)

	keys = inpututil.AppendJustPressedKeys(keys[:0])

	for _, k := range keys {
		switch k {
		case ebiten.KeyRight, ebiten.KeyD:
			input.X = 1
		case ebiten.KeyLeft, ebiten.KeyA:
			input.X = -1
		case ebiten.KeyDown, ebiten.KeyS:
			input.Y = 1
		case ebiten.KeyUp, ebiten.KeyW:
			input.Y = -1
		}
	}

	return input
}
