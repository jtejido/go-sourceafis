package relative

import (
	"math"
	"sort"

	"github.com/jtejido/sourceafis/config"
	"github.com/jtejido/sourceafis/extractor/logger"
	"github.com/jtejido/sourceafis/primitives"
)

type RelativeContrastMask struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *RelativeContrastMask {
	return &RelativeContrastMask{
		logger: logger,
	}
}

func (m *RelativeContrastMask) Compute(contrast *primitives.Matrix, blocks *primitives.BlockMap) *primitives.BooleanMatrix {
	sortedContrast := make([]float64, 0)
	it := contrast.Size().Iterator()
	for it.HasNext() {
		block := it.Next()
		sortedContrast = append(sortedContrast, contrast.GetPoint(block))
	}
	sort.Slice(sortedContrast, func(i, j int) bool {
		return sortedContrast[i] > sortedContrast[j]
	})

	pixelsPerBlock := blocks.Pixels.Area() / blocks.Primary.Blocks.Area()
	sampleCount := math.Min(float64(len(sortedContrast)), float64(config.Config.RelativeContrastSample)/float64(pixelsPerBlock))
	consideredBlocks := int(math.Max(math.Floor((sampleCount*config.Config.RelativeContrastPercentile)+.5), 1))
	averageContrast := calculateAverage(sortedContrast, consideredBlocks)
	limit := averageContrast * config.Config.MinRelativeContrast
	result := primitives.NewBooleanMatrixFromPoint(blocks.Primary.Blocks)
	it = blocks.Primary.Blocks.Iterator()
	for it.HasNext() {
		block := it.Next()
		if contrast.GetPoint(block) < limit {
			result.SetPoint(block, true)
		}
	}
	m.logger.Log("relative-contrast-mask", result)
	return result
}

func calculateAverage(sortedContrast []float64, consideredBlocks int) float64 {
	if consideredBlocks <= 0 || len(sortedContrast) == 0 {
		return 0.0
	}

	// Take the first "consideredBlocks" elements from the sortedContrast slice
	slice := sortedContrast[:consideredBlocks]

	// Calculate the sum of the elements in the slice
	sum := 0.0
	for _, value := range slice {
		sum += value
	}

	// Calculate the average by dividing the sum by the number of elements
	average := sum / float64(consideredBlocks)

	return average
}
