package histogram

import (
	"sourceafis/config"
	"sourceafis/extractor/logger"
	"sourceafis/primitives"
)

type LocalHistograms struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *LocalHistograms {
	return &LocalHistograms{
		logger: logger,
	}
}

func (h *LocalHistograms) Create(blocks *primitives.BlockMap, image *primitives.Matrix) *primitives.HistogramCube {
	histogram := primitives.NewHistogramCubeFromPoint(blocks.Primary.Blocks, config.Config.HistogramDepth)
	it := blocks.Primary.Blocks.Iterator()
	for it.HasNext() {
		block := it.Next()
		area := blocks.Primary.BlockPoint(block)
		for y := area.Top(); y < area.Bottom(); y++ {
			for x := area.Left(); x < area.Right(); x++ {
				depth := int(image.Get(x, y) * float64(histogram.Bins))

				histogram.IncrementPoint(block, histogram.Constrain(depth))
			}
		}
	}

	h.logger.Log("histogram", histogram)
	return histogram
}

func (h *LocalHistograms) Smooth(blocks *primitives.BlockMap, input *primitives.HistogramCube) *primitives.HistogramCube {
	blocksAround := []primitives.IntPoint{
		{X: 0, Y: 0},
		{X: -1, Y: 0},
		{X: 0, Y: -1},
		{X: -1, Y: -1},
	}
	output := primitives.NewHistogramCubeFromPoint(blocks.Secondary.Blocks, input.Bins)
	it := blocks.Secondary.Blocks.Iterator()
	for it.HasNext() {
		corner := it.Next()
		for _, relative := range blocksAround {
			block := corner.Plus(relative)
			if blocks.Primary.Blocks.Contains(block) {
				for i := 0; i < input.Bins; i++ {
					output.AddPoint(corner, i, input.GetPoint(block, i))
				}
			}
		}
	}
	h.logger.Log("smoothed-histogram", output)
	return output
}
