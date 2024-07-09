package components

import (
	"github.com/jejikeh/invaders/pkg/goecs"
	"github.com/jejikeh/invaders/pkg/goengine"
	"github.com/jejikeh/invaders/pkg/gomath"
)

type Player struct {
	_ bool
}

func NewPlayer(engine *goengine.Engine, pos gomath.Vec2) goecs.EntityID {
	player := engine.ECS.NewEntity()

	t := goecs.Attach[goengine.Transfrom](engine.ECS, player)
	t.Position = pos
	t.Scale = gomath.NewVec2(1, 1)

	sprite := goecs.Attach[goengine.EbitenSprite](engine.ECS, player)
	sprite.Path = "/Volumes/Dev/Projects/invaders/resources/player.png"

	goecs.Attach[Player](engine.ECS, player)

	movable := goecs.Attach[Movable](engine.ECS, player)
	movable.Acceleration = 1500
	movable.Friction = 600 / 2
	movable.Speed = 5

	return player
}
