package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/jejikeh/invaders/pkg/goecs"
	"github.com/jejikeh/invaders/pkg/goengine"
	"github.com/jejikeh/invaders/pkg/gomath"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

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

	player := engine.ECS.NewEntity()
	t := goecs.Attach[goengine.Transfrom](engine.ECS, player)
	t.Position.X = 100
	t.Position.Y = 100
	t.Scale.X = 1
	t.Scale.Y = 1

	sprite := goecs.Attach[goengine.EbitenSprite](engine.ECS, player)
	sprite.Path = "/Volumes/Dev/Projects/invaders/resources/player.png"

	goecs.Attach[PlayerTag](engine.ECS, player)

	// runtime.SetFinalizer(sprite, func(sprite *goengine.EbitenSprite) {
	// 	runtime.KeepAlive(sprite)
	// 	runtime.KeepAlive(sprite.Image)
	// })

	// disable gc
	// debug.SetGCPercent(-1)

	engine.Run()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

type PlayerTag struct {
}

func NewPlayer(engine *goengine.Engine) goecs.EntityID {
	runtime.KeepAlive(engine)

	return goecs.EntityID(1)
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
