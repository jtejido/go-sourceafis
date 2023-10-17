package primitives

import "math"

type FloatPoint struct {
	X, Y float64
}

func ZeroFloatPoint() FloatPoint {
	return FloatPoint{0, 0}
}

func (f FloatPoint) Multiply(factor float64) FloatPoint {
	return FloatPoint{factor * f.X, factor * f.Y}
}
func (f FloatPoint) Round() IntPoint {
	return IntPoint{int(math.Floor(f.X + .5)), int(math.Floor(f.Y + .5))}
}
