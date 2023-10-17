package vote

import (
	"math"
	"sourceafis/primitives"
	"sourceafis/utils"
)

func Apply(input, mask *primitives.BooleanMatrix, radius int, majority float64, borderDistance int) *primitives.BooleanMatrix {
	size := input.Size()
	rect := primitives.IntRect{
		X:      borderDistance,
		Y:      borderDistance,
		Width:  size.X - 2*borderDistance,
		Height: size.Y - 2*borderDistance,
	}
	thresholds := make([]int, 0)
	for x := 0; x < utils.SquareInt(2*radius+1)+1; x++ {
		thresholds = append(thresholds, int(math.Ceil(majority*float64(x))))
	}
	counts := primitives.NewIntMatrixFromPoint(size)
	output := primitives.NewBooleanMatrixFromPoint(size)
	for y := rect.Top(); y < rect.Bottom(); y++ {
		superTop := y - radius - 1
		superBottom := y + radius
		yMin := int(math.Max(0, float64(y-radius)))
		yMax := int(math.Min(float64(size.Y)-1, float64(y+radius)))
		yRange := yMax - yMin + 1
		for x := rect.Left(); x < rect.Right(); x++ {
			if mask == nil || (mask != nil && mask.Get(x, y)) {
				var left, top, diagonal int
				var isLeft, isTop bool
				if x > 0 {
					left = counts.Get(x-1, y)
					isLeft = true
				}
				if y > 0 {
					top = counts.Get(x, y-1)
					isTop = true
				}

				if isLeft && isTop {
					diagonal = counts.Get(x-1, y-1)
				}

				xMin := int(math.Max(0, float64(x-radius)))
				xMax := int(math.Min(float64(size.X)-1, float64(x+radius)))
				var ones int
				if left > 0 && top > 0 && diagonal > 0 {
					ones = top + left - diagonal - 1
					superLeft := x - radius - 1
					superRight := x + radius
					if superLeft >= 0 && superTop >= 0 && input.Get(superLeft, superTop) {
						ones++
					}
					if superLeft >= 0 && superBottom < size.Y && input.Get(superLeft, superBottom) {
						ones--
					}
					if superRight < size.X && superTop >= 0 && input.Get(superRight, superTop) {
						ones--
					}
					if superRight < size.X && superBottom < size.Y && input.Get(superRight, superBottom) {
						ones++
					}
				} else {
					ones = 0
					for ny := yMin; ny <= yMax; ny++ {
						for nx := xMin; nx <= xMax; nx++ {
							if input.Get(nx, ny) {
								ones++
							}
						}
					}
				}
				counts.Set(x, y, ones+1)
				if ones >= thresholds[yRange*(xMax-xMin+1)] {
					output.Set(x, y, true)
				}
			}
		}
	}
	return output
}
