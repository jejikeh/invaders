package main

import (
	"math/rand/v2"

	"github.com/jejikeh/invaders/pkg/goecs"
	"github.com/jejikeh/invaders/pkg/goengine"
	"github.com/jejikeh/invaders/pkg/gomath"
)

func main() {
	engine, err := goengine.NewEngine(&goengine.Config{
		Window: &goengine.WindowConfig{
			Title: "Invaders",
			Size:  gomath.NewVec2(1024, 768),
			Scale: 1.0,
		},
	})

	if err != nil {
		panic(err)
	}

	for range 1 {
		spawnPlayer(engine)
		spawnEnemy(engine)
	}

	engine.ECS.AddSystems(movePlayer)

	engine.Run()
}

type PlayerTag struct {
	_ bool
}

func spawnPlayer(engine *goengine.Engine) {
	player := engine.ECS.NewEntity()
	t := goecs.Attach[goengine.Transfrom](engine.ECS, player)
	t.Position.X = float64(rand.Int64N(53))
	t.Position.Y = float64(rand.Int64N(50))
	t.Scale.X = 1
	t.Scale.Y = 1

	sprite := goecs.Attach[goengine.EbitenSprite](engine.ECS, player)
	sprite.Path = "/Volumes/Dev/Projects/invaders/resources/player.png"

	goecs.Attach[PlayerTag](engine.ECS, player)
}

func spawnEnemy(engine *goengine.Engine) {
	player := engine.ECS.NewEntity()
	t := goecs.Attach[goengine.Transfrom](engine.ECS, player)
	t.Position.X = float64(rand.Int64N(103))
	t.Position.Y = float64(rand.Int64N(100))
	t.Scale.X = 1
	t.Scale.Y = 1

	sprite := goecs.Attach[goengine.EbitenSprite](engine.ECS, player)
	sprite.Path = "/Volumes/Dev/Projects/invaders/resources/alien.png"

	goecs.Attach[PlayerTag](engine.ECS, player)
}

func movePlayer(layer *goecs.Layer) {
	entities := layer.Request(
		goecs.GetComponentID[goengine.Transfrom](layer),
		goecs.GetComponentID[PlayerTag](layer),
	)

	for _, entity := range entities {
		// pt, _ := goecs.GetComponent[PlayerTag](layer, entity)
		transform, _ := goecs.GetComponent[goengine.Transfrom](layer, entity)
		// transform.Position.Add(gomath.NewVec2(randomVector.X, randomVector.Y).Scale(float64(pt.m)))
		transform.Rotation += 0.01
	}
}
