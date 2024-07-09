package main

import (
	"github.com/jejikeh/invaders/cmd/invaders2/internal/components"
	"github.com/jejikeh/invaders/cmd/invaders2/internal/systems"
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
	components.NewPlayer(engine, gomath.NewVec2(100, 100))

	engine.ECS.AddSystems(systems.Move)

	engine.Run()
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
