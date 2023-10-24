package primitives

import (
	"math"

	"github.com/jtejido/sourceafis/utils"
)

type IntPoint struct {
	X, Y int
}

var CORNER_NEIGHBORS = []IntPoint{
	{-1, -1},
	{0, -1},
	{1, -1},
	{-1, 0},
	{1, 0},
	{-1, 1},
	{0, 1},
	{1, 1},
}

var EDGE_NEIGHBORS = []IntPoint{
	{0, -1},
	{-1, 0},
	{1, 0},
	{0, 1},
}

func ZeroIntPoint() IntPoint {
	return IntPoint{0, 0}
}

func (p IntPoint) Equals(obj interface{}) bool {
	if v, ok := obj.(IntPoint); !ok {
		return false
	} else {
		return p.X == v.X && p.Y == v.Y
	}
}

func (p IntPoint) Negate() IntPoint {
	return IntPoint{-p.X, -p.Y}
}

func (p IntPoint) Contains(other IntPoint) bool {
	return other.X >= 0 && other.Y >= 0 && other.X < p.X && other.Y < p.Y
}

func (p IntPoint) Plus(other IntPoint) IntPoint {
	return IntPoint{X: p.X + other.X, Y: p.Y + other.Y}
}
func (p IntPoint) Minus(other IntPoint) IntPoint {
	return IntPoint{X: p.X - other.X, Y: p.Y - other.Y}
}

func (p IntPoint) ToFloat() FloatPoint {
	return FloatPoint{X: float64(p.X), Y: float64(p.Y)}
}

func (p IntPoint) LengthSq() int {
	return utils.SquareInt(p.X) + utils.SquareInt(p.Y)
}

func (p IntPoint) Area() int {
	return p.X * p.Y
}

func (p IntPoint) CompareTo(other IntPoint) int {
	resultY := p.Y - other.Y
	if resultY != 0 {
		return resultY
	}
	return p.X - other.X
}

func (p IntPoint) LineTo(to IntPoint) []IntPoint {
	var result []IntPoint
	relative := to.Minus(p)
	if int(math.Abs(float64(relative.X))) >= int(math.Abs(float64(relative.Y))) {
		result = make([]IntPoint, int(math.Abs(float64(relative.X)))+1)
		if relative.X > 0 {
			for i := 0; i <= relative.X; i++ {
				v := float64(i) * (float64(relative.Y) / float64(relative.X))
				result[i] = IntPoint{X: p.X + i, Y: p.Y + int(math.Floor(v+.5))}
			}
		} else if relative.X < 0 {
			for i := 0; i <= -relative.X; i++ {
				v := float64(i) * (float64(relative.Y) / float64(relative.X))
				result[i] = IntPoint{X: p.X - i, Y: p.Y - int(math.Floor(v+.5))}
			}
		} else {
			result[0] = p
		}
	} else {
		result = make([]IntPoint, int(math.Abs(float64(relative.Y)))+1)
		if relative.Y > 0 {
			for i := 0; i <= relative.Y; i++ {
				v := float64(i) * (float64(relative.X) / float64(relative.Y))
				result[i] = IntPoint{X: p.X + int(math.Floor(v+.5)), Y: p.Y + i}
			}
		} else if relative.Y < 0 {
			for i := 0; i <= -relative.Y; i++ {
				v := float64(i) * (float64(relative.X) / float64(relative.Y))
				result[i] = IntPoint{X: p.X - int(math.Floor(v+.5)), Y: p.Y - i}
			}
		} else {
			result[0] = p
		}
	}
	return result
}

func (p IntPoint) Iterator() IntPointIterator {
	return newIterator(p)
}
