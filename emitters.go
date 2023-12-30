// This file from `https://github.com/jejikeh/go-libpartikel/blob/master/partikel/partikel.go`

package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Particle struct {
	Origin               rl.Vector2
	Position             rl.Vector2
	Velocity             rl.Vector2
	ExternalAcceleration rl.Vector2
	OriginAcceleration   float32
	Age                  float32
	TTL                  float32
	Active               bool
}

func (p Particle) IsDead() bool {
	return p.Age > p.TTL
}

func NewParticle(config EmitterConfig) *Particle {
	p := new(Particle)
	p.Age = 0
	p.Origin = config.Origin

	randomAngle := getRandomFloatRange(config.DirectionAngle)
	res := rotateVector2(config.Direction, randomAngle)

	randomVelocity := getRandomFloatRange(config.Velocity)
	p.Velocity = rl.Vector2{
		X: res.X * randomVelocity,
		Y: res.Y * randomVelocity,
	}

	randomVelocityAngle := getRandomFloatRange(config.VelocityAngle)
	p.Velocity = rotateVector2(p.Velocity, randomVelocityAngle)

	randomOffset := getRandomFloatRange(config.Offset)
	p.Position = rl.Vector2{
		X: config.Origin.X + res.X*randomOffset,
		Y: config.Origin.Y + res.Y*randomOffset,
	}

	randomOriginAcceleration := getRandomFloatRange(config.OriginAcceleration)
	p.OriginAcceleration = randomOriginAcceleration
	p.ExternalAcceleration = config.ExternalAcceleration
	p.TTL = getRandomFloatRange(config.Age)
	p.Active = true

	return p
}

func (p *Particle) Update() {
	if !p.Active {
		return
	}

	p.Age += rl.GetFrameTime()

	if p.IsDead() {
		p.Active = false
		return
	}

	toOrigin := rl.Vector2Normalize(rl.NewVector2(p.Origin.X-p.Position.X, p.Origin.Y-p.Position.Y))
	p.Velocity.X += toOrigin.X * p.OriginAcceleration * rl.GetFrameTime()
	p.Velocity.Y += toOrigin.Y * p.OriginAcceleration * rl.GetFrameTime()

	p.Velocity.X += p.ExternalAcceleration.X * rl.GetFrameTime()
	p.Velocity.Y += p.ExternalAcceleration.Y * rl.GetFrameTime()

	p.Position.X += p.Velocity.X * rl.GetFrameTime()
	p.Position.Y += p.Velocity.Y * rl.GetFrameTime()
}

func rotateVector2(vec rl.Vector2, angle float32) rl.Vector2 {
	radians := float64(angle * math.Pi / 180)
	return rl.Vector2{
		X: float32(math.Cos(radians))*vec.X - float32(math.Sin(radians))*vec.Y,
		Y: float32(math.Sin(radians))*vec.X + float32(math.Cos(radians))*vec.Y,
	}
}

type FloatRange [2]float32

type IntRange [2]int

type Vectror2Range [2]rl.Vector2

type EmitterConfig struct {
	// TODO: Make random StartSize and EndSize using Vector2Range
	// But in Draw() we need to use Vector2, so we need another field for CurrentSize
	StartSize            rl.Vector2
	EndSize              rl.Vector2
	Direction            rl.Vector2
	Velocity             FloatRange
	DirectionAngle       FloatRange
	VelocityAngle        FloatRange
	Offset               FloatRange
	OriginAcceleration   FloatRange
	Burst                IntRange
	Capacity             int
	EmmisionRate         int
	Origin               rl.Vector2
	ExternalAcceleration rl.Vector2
	StartColor           rl.Color
	EndColor             rl.Color
	Age                  FloatRange
	BlendMode            rl.BlendMode
	Texture              rl.Texture2D
}

func getRandomFloatRange(rng FloatRange) float32 {
	return float32(float64(rng[0]) + rand.Float64()*float64((rng[1]-rng[0])))
}

func getRandomIntRange(rng IntRange) int {
	return rand.Intn(rng[1]-rng[0]) + rng[0]
}

// func getRandomVector2Range(rng Vectror2Range) rl.Vector2 {
// 	var maxValueX float32
// 	var maxValueY float32

// 	minValueX := float32(math.Min(float64(rng[0].X), float64(rng[1].X)))
// 	minValueY := float32(math.Min(float64(rng[0].Y), float64(rng[1].Y)))

// 	if minValueX == rng[0].X {
// 		maxValueX = rng[1].X
// 	} else {
// 		maxValueX = rng[0].X
// 	}

// 	if minValueY == rng[0].Y {
// 		maxValueY = rng[1].Y
// 	} else {
// 		maxValueY = rng[0].Y
// 	}

// 	return rl.Vector2{
// 		X: getRandomFloatRange([2]float32{minValueX, maxValueX}),
// 		Y: getRandomFloatRange([2]float32{minValueY, maxValueY}),
// 	}
// }

type Emitter struct {
	Config    EmitterConfig
	EmitCount float32
	Offset    rl.Vector2
	Active    bool
	Particles []*Particle
}

func NewEmitter(config EmitterConfig) *Emitter {
	e := new(Emitter)

	e.Config = config

	e.Offset.X = float32(e.Config.Texture.Width) / 2
	e.Offset.Y = float32(e.Config.Texture.Height) / 2

	e.Particles = make([]*Particle, 0, e.Config.Capacity)

	e.Config.Direction = rl.Vector2Normalize(e.Config.Direction)

	for i := 0; i < e.Config.Capacity; i++ {
		e.Particles = append(e.Particles, NewParticle(e.Config))
	}

	return e
}

func (e *Emitter) Start() {
	e.Active = true
}

func (e *Emitter) Stop() {
	e.Active = false
}

func (e *Emitter) Burst() {
	emitted := 0
	amount := getRandomIntRange(e.Config.Burst)

	var p *Particle

	for i := 0; i < e.Config.Capacity; i++ {
		p = e.Particles[i]
		if !p.Active {
			// TODO: Remove allocation here?
			p = NewParticle(e.Config)
			e.Particles[i] = p
			p.Position = e.Config.Origin
			emitted += 1
		}

		if emitted >= int(amount) {
			break
		}
	}
}

func (e *Emitter) Update() {
	emitNow := 0

	var p *Particle
	counter := 0

	if e.Active {
		e.EmitCount += rl.GetFrameTime() * float32(e.Config.EmmisionRate)
		emitNow = int(e.EmitCount)
	}

	for i := 0; i < e.Config.Capacity; i++ {
		p = e.Particles[i]
		if p.Active {
			p.Update()
			counter += 1
		} else {
			p = NewParticle(e.Config)
			e.Particles[i] = p

			p.Update()

			emitNow -= 1
			e.EmitCount -= 1
			counter += 1
		}
	}
}

func (e *Emitter) Draw() {
	rl.BeginBlendMode(e.Config.BlendMode)

	for _, p := range e.Particles {
		if p.Active {
			size := linearVectorFade(e.Config.StartSize, e.Config.EndSize, p.Age/p.TTL)
			textureSizeX := float32(e.Config.Texture.Width) * size.X
			textureSizeY := float32(e.Config.Texture.Height) * size.Y
			origin := rl.NewVector2(
				((float32(e.Config.Texture.Width))*size.X)/2,
				((float32(e.Config.Texture.Height))*size.Y)/2,
			)
			rl.DrawTexturePro(
				e.Config.Texture,
				rl.Rectangle{
					X:      0,
					Y:      0,
					Width:  float32(e.Config.Texture.Width),
					Height: -float32(e.Config.Texture.Height)},
				rl.Rectangle{
					X:      p.Position.X - e.Offset.X,
					Y:      p.Position.Y - e.Offset.Y,
					Width:  textureSizeX,
					Height: textureSizeY,
				},
				origin,
				0,
				linearColorFade(
					e.Config.StartColor,
					e.Config.EndColor,
					p.Age/p.TTL))
		}
	}

	rl.EndBlendMode()
}

func linearColorFade(c1, c2 rl.Color, fraction float32) rl.Color {
	newR := int32(float32(c1.R) + (float32(c2.R)-float32(c1.R))*fraction)
	newG := int32(float32(c1.G) + (float32(c2.G)-float32(c1.G))*fraction)
	newB := int32(float32(c1.B) + (float32(c2.B)-float32(c1.B))*fraction)
	newA := int32(float32(c1.A) + (float32(c2.A)-float32(c1.A))*fraction)

	return rl.Color{R: uint8(newR), G: uint8(newG), B: uint8(newB), A: uint8(newA)}
}

func linearVectorFade(v1, v2 rl.Vector2, fraction float32) rl.Vector2 {
	return rl.Vector2{
		X: v1.X + (v2.X-v1.X)*fraction,
		Y: v1.Y + (v2.Y-v1.Y)*fraction,
	}
}

type ParticleSystem struct {
	Emitters []*Emitter
	Count    int
	Origin   rl.Vector2
	Active   bool
}

func (p *ParticleSystem) Add(emitter *Emitter) {
	p.Emitters = append(p.Emitters, emitter)
	p.Count += 1
}

func (p *ParticleSystem) Remove(emitter *Emitter) bool {
	for i, e := range p.Emitters {
		if e == emitter {
			p.Emitters = append(p.Emitters[:i], p.Emitters[i+1:]...)
			p.Count -= 1
			return true
		}
	}

	return false
}

func (p *ParticleSystem) SetOrigin(origin rl.Vector2) {
	p.Origin = origin
	for _, e := range p.Emitters {
		e.Config.Origin = p.Origin
	}
}

func (p *ParticleSystem) Start() {
	p.Active = true
	for _, e := range p.Emitters {
		e.Start()
	}
}

func (p *ParticleSystem) Stop() {
	p.Active = false
	for _, e := range p.Emitters {
		e.Stop()
	}
}

func (p *ParticleSystem) Burst() {
	for _, e := range p.Emitters {
		e.Burst()
	}
}

func (p *ParticleSystem) Draw() {
	for _, e := range p.Emitters {
		e.Draw()
	}
}

func (p *ParticleSystem) Update() {
	for _, e := range p.Emitters {
		e.Update()
	}
}
