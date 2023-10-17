package utils

import "math"

func SquareInt(x int) int {
	return x * x
}

func SquareFloat64(x float64) float64 {
	return x * x
}

func RoundUpDiv(dividend, divisor int) int {
	return (dividend + divisor - 1) / divisor
}

func Interpolate(start, end, position float64) float64 {
	return start + position*(end-start)
}

func InterpolateRect(bottomleft, bottomright, topleft, topright, x, y float64) float64 {
	left := Interpolate(topleft, bottomleft, y)
	right := Interpolate(topright, bottomright, y)
	return Interpolate(left, right, x)
}

func InterpolateExponential(start, end, position float64) float64 {
	return math.Pow(end/start, position) * start
}
