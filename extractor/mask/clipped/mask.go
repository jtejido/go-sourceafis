package clipped

import (
	"math"
	"sourceafis/config"
	"sourceafis/extractor/logger"
	"sourceafis/primitives"
)

type ClippedContrast struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *ClippedContrast {
	return &ClippedContrast{
		logger: logger,
	}
}

func (m *ClippedContrast) Compute(blocks *primitives.BlockMap, histogram *primitives.HistogramCube) *primitives.Matrix {
	result := primitives.NewMatrixFromPoint(blocks.Primary.Blocks)
	it := blocks.Primary.Blocks.Iterator()
	for it.HasNext() {
		block := it.Next()
		volume := histogram.SumPoint(block)
		clipLimit := int(math.Floor((float64(volume) * config.Config.ClippedContrast) + .5))
		accumulator := 0
		lowerBound := histogram.Bins - 1
		for i := 0; i < histogram.Bins; i++ {
			accumulator += histogram.GetPoint(block, i)
			if accumulator > clipLimit {
				lowerBound = i
				break
			}
		}

		accumulator = 0
		upperBound := 0
		for i := histogram.Bins - 1; i >= 0; i-- {
			accumulator += histogram.GetPoint(block, i)
			if accumulator > clipLimit {
				upperBound = i
				break
			}
		}
		result.SetPoint(block, float64(upperBound-lowerBound)*(1.0/float64(histogram.Bins-1)))
	}
	m.logger.Log("contrast", result)
	return result
}
