package equalizer

import (
	"math"
	"sourceafis/config"
	"sourceafis/extractor/logger"
	"sourceafis/primitives"
	"sourceafis/utils"
)

type ImageEqualization struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *ImageEqualization {
	return &ImageEqualization{
		logger: logger,
	}
}

func (e *ImageEqualization) Equalize(blocks *primitives.BlockMap, image *primitives.Matrix, histogram *primitives.HistogramCube, blockMask *primitives.BooleanMatrix) *primitives.Matrix {
	rangeMin := -1.
	rangeMax := 1.
	rangeSize := rangeMax - rangeMin
	widthMax := rangeSize / 256 * config.Config.MaxEqualizationScaling
	widthMin := rangeSize / 256 * config.Config.MinEqualizationScaling
	limitedMin := make([]float64, histogram.Bins)
	limitedMax := make([]float64, histogram.Bins)
	dequantized := make([]float64, histogram.Bins)
	for i := 0; i < histogram.Bins; i++ {
		limitedMin[i] = math.Max(float64(i)*widthMin+rangeMin, rangeMax-(float64(histogram.Bins)-1-float64(i))*widthMax)
		limitedMax[i] = math.Min(float64(i)*widthMax+rangeMin, rangeMax-(float64(histogram.Bins)-1-float64(i))*widthMin)
		dequantized[i] = float64(i) / (float64(histogram.Bins) - 1)
	}
	mappings := make(map[primitives.IntPoint][]float64)
	it := blocks.Secondary.Blocks.Iterator()
	for it.HasNext() {
		corner := it.Next()
		mapping := make([]float64, histogram.Bins)
		mappings[corner] = mapping
		if blockMask.GetPointWithFallback(corner, false) || blockMask.GetWithFallback(corner.X-1, corner.Y, false) || blockMask.GetWithFallback(corner.X, corner.Y-1, false) || blockMask.GetWithFallback(corner.X-1, corner.Y-1, false) {
			step := rangeSize / float64(histogram.SumPoint(corner))
			top := rangeMin
			for i := 0; i < histogram.Bins; i++ {
				band := float64(histogram.GetPoint(corner, i)) * step
				equalized := top + dequantized[i]*band
				top += band
				if equalized < limitedMin[i] {
					equalized = limitedMin[i]
				}
				if equalized > limitedMax[i] {
					equalized = limitedMax[i]
				}
				mapping[i] = equalized
			}
		}
	}
	result := primitives.NewMatrixFromPoint(blocks.Pixels)
	it = blocks.Primary.Blocks.Iterator()
	for it.HasNext() {
		block := it.Next()
		area := blocks.Primary.BlockPoint(block)
		if blockMask.GetPoint(block) {
			topleft := mappings[block]
			topright := mappings[primitives.IntPoint{X: block.X + 1, Y: block.Y}]
			bottomleft := mappings[primitives.IntPoint{X: block.X, Y: block.Y + 1}]
			bottomright := mappings[primitives.IntPoint{X: block.X + 1, Y: block.Y + 1}]
			for y := area.Top(); y < area.Bottom(); y++ {
				for x := area.Left(); x < area.Right(); x++ {
					depth := histogram.Constrain(int(image.Get(x, y) * float64(histogram.Bins)))
					rx := (float64(x) - float64(area.X) + 0.5) / float64(area.Width)
					ry := (float64(y) - float64(area.Y) + 0.5) / float64(area.Height)
					result.Set(x, y, utils.InterpolateRect(bottomleft[depth], bottomright[depth], topleft[depth], topright[depth], rx, ry))
				}
			}
		} else {
			for y := area.Top(); y < area.Bottom(); y++ {
				for x := area.Left(); x < area.Right(); x++ {
					result.Set(x, y, -1)
				}
			}
		}
	}
	e.logger.Log("equalized-image", result)
	return result
}
