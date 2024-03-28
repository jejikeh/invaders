package components

type Transform struct {
	X, Y          float64
	Width, Height float64
	Angle         float64
}

func NewTransform(x, y, width, height, angle float64) *Transform {
	t := &Transform{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		Angle:  angle,
	}

	return t
}
