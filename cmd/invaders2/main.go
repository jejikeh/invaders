package main

import (
	"github.com/jejikeh/invaders/pkg/goecs"
	"github.com/jejikeh/invaders/pkg/goengine"
	"github.com/jejikeh/invaders/pkg/gomath"
)

func main() {
	engine, err := goengine.NewEngine(&goengine.Config{
		Window: &goengine.WindowConfig{
			Title: "Invaders",
			Size:  gomath.NewVec2(800, 600),
			Scale: 2,
		},
	})

	if err != nil {
		panic(err)
	}

	engine.ECS.AddSystems(movePlayer)

	NewPlayer(engine)

	engine.Run()
}

type PlayerTag struct{}

func NewPlayer(engine *goengine.Engine) goecs.EntityID {
	player := engine.ECS.NewEntity()
	goecs.Attach[goengine.Transfrom](engine.ECS, player)

	sprite := goecs.Attach[goengine.EbitenSprite](engine.ECS, player)
	sprite.Path = "/Volumes/Dev/Projects/invaders/resources/player.png"

	goecs.Attach[PlayerTag](engine.ECS, player)

	return player
}

func movePlayer(layer *goecs.Layer) {
	entities := layer.Request(
		goecs.GetComponentID[goengine.Transfrom](layer),
		goecs.GetComponentID[PlayerTag](layer),
	)

	for _, entity := range entities {
		transform, _ := goecs.GetComponent[goengine.Transfrom](layer, entity)
		transform.Position.X += 1
	}
}
