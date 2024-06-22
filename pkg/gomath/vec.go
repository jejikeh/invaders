package gomath

type Vec2 struct {
	X float64
	Y float64
}

func NewVec2(x, y float64) *Vec2 {
	return &Vec2{
		X: x,
		Y: y,
	}
}

func NewVec2FromVec(v *Vec2) *Vec2 {
	return &Vec2{
		X: v.X,
		Y: v.Y,
	}
}

func (v *Vec2) Add(v2 *Vec2) *Vec2 {
	v.X += v2.X
	v.Y += v2.Y

	return v
}

func (v *Vec2) Scale(scalar float64) *Vec2 {
	if scalar == 0.0 {
		return v
	}
	v.X *= scalar
	v.Y *= scalar

	return v
}
