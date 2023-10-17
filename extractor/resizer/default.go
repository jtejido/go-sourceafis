package resizer

import (
	"math"
	"sourceafis/primitives"
)

func Resize(input *primitives.Matrix, dpi float64) *primitives.Matrix {
	return resize(input, int(math.Floor((500.0/dpi*float64(input.Width))+.5)), int(math.Floor((500.0/dpi*float64(input.Height))+.5)))
}

func resize(input *primitives.Matrix, newWidth, newHeight int) *primitives.Matrix {
	if newWidth == input.Width && newHeight == input.Height {
		return input
	}

	output := primitives.NewMatrix(newWidth, newHeight)
	scaleX := float64(newWidth) / float64(input.Width)
	scaleY := float64(newHeight) / float64(input.Height)
	descaleX := 1 / scaleX
	descaleY := 1 / scaleY
	for y := 0; y < newHeight; y++ {
		y1 := float64(y) * descaleY
		y2 := y1 + descaleY
		y1i := int(y1)
		y2i := int(math.Min(math.Ceil(y2), float64(input.Height)))
		for x := 0; x < newWidth; x++ {
			x1 := float64(x) * descaleX
			x2 := x1 + descaleX
			x1i := int(x1)
			x2i := int(math.Min(math.Ceil(x2), float64(input.Width)))
			var sum float64
			for oy := y1i; oy < y2i; oy++ {
				ry := math.Min(float64(oy)+1, y2) - math.Max(float64(oy), y1)
				for ox := x1i; ox < x2i; ox++ {
					rx := math.Min(float64(ox)+1, x2) - math.Max(float64(ox), x1)
					sum += rx * ry * input.Get(ox, oy)
				}
			}
			output.Set(x, y, sum*(scaleX*scaleY))
		}
	}
	return output
}
