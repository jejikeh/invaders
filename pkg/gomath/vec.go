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

func (v *Vec2) DivFloat(scalar float64) *Vec2 {
	return NewVec2(v.X/scalar, v.Y/scalar)
}
