package primitives

import "math"

type IntRect struct {
	X, Y, Width, Height int
}

func IntRectFromPoint(size IntPoint) IntRect {
	return IntRect{
		X:      0,
		Y:      0,
		Width:  size.X,
		Height: size.Y,
	}
}

func IntRectBetween(startX, startY, endX, endY int) IntRect {
	return IntRect{
		X:      startX,
		Y:      startY,
		Width:  endX - startX,
		Height: endY - startY,
	}
}

func IntRectAround(x, y, radius int) IntRect {
	return IntRectBetween(x-radius, y-radius, x+radius+1, y+radius+1)
}

func IntRectBetweenIntPoint(start, end IntPoint) IntRect {
	return IntRectBetween(start.X, start.Y, end.X, end.Y)
}

func IntRectAroundIntPoint(center IntPoint, radius int) IntRect {
	return IntRectAround(center.X, center.Y, radius)
}

func (r IntRect) Left() int {
	return r.X
}
func (r IntRect) Top() int {
	return r.Y
}
func (r IntRect) Right() int {
	return r.X + r.Width
}
func (r IntRect) Bottom() int {
	return r.Y + r.Height
}
func (r IntRect) Area() int {
	return r.Width * r.Height
}

func (r IntRect) Center() IntPoint {
	return IntPoint{
		X: (r.Left() + r.Right()) / 2,
		Y: (r.Top() + r.Bottom()) / 2,
	}
}

func (r IntRect) Move(delta IntPoint) IntRect {
	return IntRect{r.X + delta.X, r.Y + delta.Y, r.Width, r.Height}
}

func (r IntRect) Intersect(other IntRect) IntRect {
	return IntRectBetweenIntPoint(
		IntPoint{int(math.Max(float64(r.Left()), float64(other.Left()))), int(math.Max(float64(r.Top()), float64(other.Top())))},
		IntPoint{int(math.Min(float64(r.Right()), float64(other.Right()))), int(math.Min(float64(r.Bottom()), float64(other.Bottom())))},
	)
}

func (p IntRect) Iterator() IntPointIterator {
	return newBlockIterator(p)
}
