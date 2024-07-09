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
	for _, entity := range layer.Request(goecs.GetComponentID[components.Player](layer)) {
		transform, _ := goecs.GetComponent[goengine.Transfrom](layer, entity)
		movable, _ := goecs.GetComponent[components.Movable](layer, entity)

		transform.Position = transform.Position.Add(
			movable.Direction.Scale(movable.Speed),
		)
	}
}

func MoveWithInput(layer *goecs.Layer) {
	input := inputVector()

	for _, entity := range layer.Request(goecs.GetComponentID[components.Player](layer)) {
		movable, _ := goecs.GetComponent[components.Movable](layer, entity)

		movable.Direction = input
	}
}

func inputVector() gomath.Vec2 {
	input := gomath.NewVec2(0, 0)

	keys = inpututil.AppendPressedKeys(keys[:0])

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
