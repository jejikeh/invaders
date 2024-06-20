package main

import (
	"github.com/jejikeh/invaders/pkg/goengine"
	"github.com/jejikeh/invaders/pkg/gomath"
)

func main() {
	engine, err := goengine.NewEngine(&goengine.Config{
		Window: &goengine.WindowConfig{
			Title:    "Invaders",
			Provider: goengine.Ebiten,
			Size:     gomath.NewVec2(800, 600),
			Scale:	  2,
		},
	})
	
	if err != nil {
		panic(err)
	}
	
	engine.UpdateConfig(func (c *goengine.Config) {
		c.Window.Title = "in"
	})
	
	
}
