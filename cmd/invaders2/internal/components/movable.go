package components

import "github.com/jejikeh/invaders/pkg/gomath"

type Movable struct {
	Direction    gomath.Vec2
	Speed        float64
	Acceleration float64
	Friction     float64
	Velocity     gomath.Vec2
}
