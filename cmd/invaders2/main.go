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

	createNewBackground(engine)
	spawnPlayer(engine)
	// spawnEnemy(engine)

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

	shader := goecs.Attach[goengine.Shader](engine.ECS, player)
	shader.Path = "/Volumes/Dev/Projects/invaders/resources/shaders/dissolve.kage"
	shader.Name = "dissolve"
}

func createNewBackground(engine *goengine.Engine) {
	bg := engine.ECS.NewEntity()
	s := goecs.Attach[goengine.EbitenSprite](engine.ECS, bg)
	goengine.NewScreenImage(engine, s)

	t := goecs.Attach[goengine.Transfrom](engine.ECS, bg)
	t.Scale = *gomath.NewVec2(1, 1)

	shader := goecs.Attach[goengine.Shader](engine.ECS, bg)
	shader.Path = "/Volumes/Dev/Projects/invaders/resources/shaders/bg.kage"
	shader.Name = "bg"
}

func movePlayer(layer *goecs.Layer) {
	entities := layer.Request(
		goecs.GetComponentID[goengine.Transfrom](layer),
		goecs.GetComponentID[PlayerTag](layer),
	)

	for _, entity := range entities {
		transform, _ := goecs.GetComponent[goengine.Transfrom](layer, entity)
		transform.Rotation += 0.01
	}
}
