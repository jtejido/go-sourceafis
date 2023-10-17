package primitives

import (
	"math"
)

const (
	Pi     = 3.141592653589793
	Pi2    = 2 * Pi
	HalfPi = 0.5 * Pi
	InvPi2 = 1.0 / Pi2
)

type FloatAngle float64

func (f FloatAngle) Difference(second FloatAngle) FloatAngle {
	angle := f - second
	if angle >= 0 {
		return angle
	}

	return angle + Pi2
}

func (f FloatAngle) Complementary() FloatAngle {
	complement := Pi2 - f
	if complement < Pi2 {
		return complement
	}
	return complement - Pi2
}

func (f FloatAngle) Opposite() FloatAngle {
	if f < FloatAngle(Pi) {
		return f + FloatAngle(Pi)
	}

	return f - FloatAngle(Pi)
}

func (f FloatAngle) ToVector() FloatPoint {
	return FloatPoint{
		X: math.Cos(float64(f)),
		Y: math.Sin(float64(f)),
	}
}

func (f FloatAngle) FromOrientation() FloatAngle {
	return FloatAngle(0.5) * f
}

func (f FloatAngle) ToOrientation() FloatAngle {
	if f < Pi {
		return 2 * f
	}

	return 2 * (f - Pi)
}

func (f FloatAngle) Quantize(resolution int) int {
	result := int(float64(f) * InvPi2 * float64(resolution))
	if result < 0 {
		return 0
	} else if result >= resolution {
		return resolution - 1
	} else {
		return result
	}
}

func (f FloatAngle) Distance(second FloatAngle) float64 {
	delta := math.Abs(float64(f) - float64(second))
	if delta <= Pi {
		return delta
	}
	return Pi2 - delta
}

func AngleAdd(start, delta float64) FloatAngle {
	angle := FloatAngle(start + delta)
	if angle < Pi2 {
		return angle
	}
	return angle - Pi2
}

func AtanFromFloatPointVector(vector FloatPoint) FloatAngle {
	angle := FloatAngle(math.Atan2(vector.Y, vector.X))
	if angle >= 0 {
		return angle
	}
	return angle + Pi2
}
func AtanFromIntPointVector(vector IntPoint) FloatAngle {
	return AtanFromFloatPointVector(vector.ToFloat())
}

func AtanFromIntPoints(center, point IntPoint) FloatAngle {
	return AtanFromIntPointVector(point.Minus(center))
}

func BucketCenter(bucket, resolution int) FloatAngle {
	return FloatAngle(Pi2) * (2*FloatAngle(bucket) + 1) / (2 * FloatAngle(resolution))
}
