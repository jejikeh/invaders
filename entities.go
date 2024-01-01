package main

import (
	"math"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type EntityInterface interface {
	Start()
	FirstDraw()
	Draw()
	Update()
	GetLayer() int
	GetRectangle() rl.Rectangle
}

type Entity struct {
	Texture      *rl.Texture2D
	Position     rl.Vector2
	Size         rl.Vector2
	Rotation     float32
	Tint         rl.Color
	ShadowHeight float32
	ShadowTint   rl.Color
	Visible      bool
	Velocity     rl.Vector2
	MaxSpeed     float32
	Layer        int
}

func newEntity(
	texture rl.Texture2D,
	position rl.Vector2,
	size rl.Vector2,
	rotation float32,
	tint rl.Color,
) *Entity {
	entity := new(Entity)
	entity.Texture = &texture
	entity.Position = position
	entity.Size = size
	entity.Rotation = rotation
	entity.Tint = tint
	entity.ShadowHeight = 5
	entity.ShadowTint = rl.NewColor(c(.0), c(.0), c(.0), c(1))
	entity.Visible = true
	entity.MaxSpeed = 400
	// entity.addEffect(&RespanwEffect{
	// 	Entity:   entity,
	// 	Lifetime: RespawnEffectDefaultLifetime,
	// })

	return entity
}

func (e *Entity) Start() {
}

func (e *Entity) FirstDraw() {
	entityHeight := float32(math.Abs(float64(e.ShadowHeight))) + 2.5
	shadowColor := e.ShadowTint

	// if !e.containsEffectOfType(Respawn) {
	shadowColor.A = c(1 / (entityHeight + 1))
	// }

	// render shadow
	Renderer.RenderTexture2D(
		e.Texture,
		rl.Vector2AddValue(e.Position, entityHeight),
		rl.Vector2AddValue(e.Size, entityHeight/100),
		e.Rotation,
		shadowColor,
	)
}

func (e *Entity) Draw() {
	if e == nil || !e.Visible {
		return
	}

	entityHeight := float32(math.Abs(float64(e.ShadowHeight))) + 2.5

	// Render entity
	Renderer.RenderTexture2D(
		e.Texture,
		rl.Vector2Subtract(e.Position, rl.NewVector2(0, entityHeight)),
		rl.Vector2AddValue(e.Size, entityHeight/100),
		e.Rotation,
		e.Tint,
	)
}

func (e *Entity) Update() {
	e.ShadowHeight = float32((math.Sin(float64(rl.GetTime()) / 1.2)) * 8)

	e.Position.X += float32(math.Sin(float64(float32(rl.GetTime()*2)))) / 10
	e.Position.Y += float32(math.Sin(float64(float32(rl.GetTime()*2)))) / 10

	e.Position.X += e.Velocity.X * rl.GetFrameTime()
	e.Position.Y += e.Velocity.Y * rl.GetFrameTime()
}

func (e *Entity) GetRectangle() rl.Rectangle {

	entityHeight := float32(math.Abs(float64(e.ShadowHeight))) + 2.5
	size := rl.Vector2AddValue(e.Size, entityHeight/100)

	textureSizeX := float32(e.Texture.Width) * size.X
	textureSizeY := float32(e.Texture.Height) * size.Y

	sizeRectangle := rl.NewVector2(textureSizeX, textureSizeY)

	position := rl.Vector2Subtract(e.Position, rl.NewVector2(0, entityHeight))

	return rl.NewRectangle(position.X, position.Y, sizeRectangle.X, sizeRectangle.Y)
}

func (e *Entity) GetLayer() int {
	return e.Layer
}

type EntityManager struct {
	Entities []EntityInterface
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		Entities: []EntityInterface{},
	}
}

func (em *EntityManager) Add(entity EntityInterface) {
	em.Entities = append(em.Entities, entity)

	sort.Slice(em.Entities, func(i0, i1 int) bool {
		return em.Entities[i0].GetLayer() < em.Entities[i1].GetLayer()
	})
}

func (em *EntityManager) Start() {
	for _, entity := range em.Entities {
		entity.Start()
	}
}

func (em *EntityManager) FirstDraw() {
	for _, entity := range em.Entities {
		entity.FirstDraw()
	}
}

func (em *EntityManager) Draw() {
	for _, entity := range em.Entities {
		entity.Draw()
	}
}

func (em *EntityManager) Update() {
	for _, entity := range em.Entities {
		entity.Update()
	}
}

type Player struct {
	*Entity
	Speed        float32
	Score        int
	Lives        int
	Acceleration float32
	Friction     float32
}

func NewPlayer() *Player {
	player := &Player{
		Entity: newEntity(
			*Assets.TexturesManager.Player,
			rl.NewVector2(WindowWidth/2, WindowHeight/1.2),
			rl.NewVector2(EntitiesBaseSize, EntitiesBaseSize),
			.0,
			rl.White,
		),
		Speed:        5,
		Lives:        3,
		Acceleration: 1500,
		Friction:     600 / 2,
		// EngineEmitters: initEngineParticleSystem(rl.NewVector2(WindowWidth/2, WindowHeight/1.2)),
	}

	// player.EngineEmitters.Start()

	return player
}

func (p *Player) Start() {
	p.Entity.Start()
}

func (p *Player) Update() {
	p.Entity.Update()

	p.updateMovement()

	if p.Position.Y > WindowHeight {
		p.Position.Y = WindowHeight
	}

	if p.Position.Y < 0 {
		p.Position.Y = 0
	}

	if p.Position.X > WindowWidth {
		p.Position.X = 0
	}

	if p.Position.X < 0 {
		p.Position.X = WindowWidth
	}
}

func (p *Player) GetLayer() int {
	return PlayerLayer
}
