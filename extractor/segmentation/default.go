package segmentation

import (
	"github.com/jtejido/sourceafis/config"
	"github.com/jtejido/sourceafis/extractor/filter/vote"
	"github.com/jtejido/sourceafis/extractor/logger"
	"github.com/jtejido/sourceafis/extractor/mask/absolute"
	"github.com/jtejido/sourceafis/extractor/mask/clipped"
	"github.com/jtejido/sourceafis/extractor/mask/relative"
	"github.com/jtejido/sourceafis/primitives"
)

type SegmentationMask struct {
	logger   logger.TransparencyLogger
	clipped  *clipped.ClippedContrast
	absolute *absolute.AbsoluteContrastMask
	relative *relative.RelativeContrastMask
}

func New(logger logger.TransparencyLogger) *SegmentationMask {
	return &SegmentationMask{
		logger:   logger,
		clipped:  clipped.New(logger),
		absolute: absolute.New(logger),
		relative: relative.New(logger),
	}
}

func filter(input *primitives.BooleanMatrix) *primitives.BooleanMatrix {
	return vote.Apply(input, nil, config.Config.BlockErrorsVoteRadius, config.Config.BlockErrorsVoteMajority, config.Config.BlockErrorsVoteBorderDistance)
}

func (m *SegmentationMask) Compute(blocks *primitives.BlockMap, histogram *primitives.HistogramCube) (*primitives.BooleanMatrix, error) {
	contrast := m.clipped.Compute(blocks, histogram)
	mask := m.absolute.Compute(contrast)
	if err := mask.Merge(m.relative.Compute(contrast, blocks)); err != nil {
		return nil, err
	}
	if err := m.logger.Log("combined-mask", mask); err != nil {
		return nil, err
	}
	if err := mask.Merge(filter(mask)); err != nil {
		return nil, err
	}

	mask.Invert()

	if err := mask.Merge(filter(mask)); err != nil {
		return nil, err
	}

	if err := mask.Merge(filter(mask)); err != nil {
		return nil, err
	}
	if err := mask.Merge(vote.Apply(mask, nil, config.Config.MaskVoteRadius, config.Config.MaskVoteMajority, config.Config.MaskVoteBorderDistance)); err != nil {
		return nil, err
	}

	return mask, m.logger.Log("filtered-mask", mask)
}

func (m *SegmentationMask) Pixelwise(mask *primitives.BooleanMatrix, blocks *primitives.BlockMap) (*primitives.BooleanMatrix, error) {
	pixelized := primitives.NewBooleanMatrixFromPoint(blocks.Pixels)
	it := blocks.Primary.Blocks.Iterator()
	for it.HasNext() {
		block := it.Next()
		if mask.GetPoint(block) {
			pixelIterator := blocks.Primary.BlockPoint(block).Iterator()
			for pixelIterator.HasNext() {
				pixel := pixelIterator.Next()
				pixelized.SetPoint(pixel, true)
			}

		}
	}

	return pixelized, m.logger.Log("pixel-mask", pixelized)
}

func shrink(mask *primitives.BooleanMatrix, amount int) *primitives.BooleanMatrix {
	size := mask.Size()
	shrunk := primitives.NewBooleanMatrixFromPoint(size)
	for y := amount; y < size.Y-amount; y++ {
		for x := amount; x < size.X-amount; x++ {
			shrunk.Set(x, y, mask.Get(x, y-amount) && mask.Get(x, y+amount) && mask.Get(x-amount, y) && mask.Get(x+amount, y))
		}
	}
	return shrunk
}

func (m *SegmentationMask) Inner(outer *primitives.BooleanMatrix) (*primitives.BooleanMatrix, error) {
	size := outer.Size()
	inner := primitives.NewBooleanMatrixFromPoint(size)
	for y := 1; y < size.Y-1; y++ {
		for x := 1; x < size.X-1; x++ {
			inner.Set(x, y, outer.Get(x, y))
		}
	}
	if config.Config.InnerMaskBorderDistance >= 1 {
		inner = shrink(inner, 1)
	}

	total := 1
	for step := 1; total+step <= config.Config.InnerMaskBorderDistance; step *= 2 {
		inner = shrink(inner, step)
		total += step
	}
	if total < config.Config.InnerMaskBorderDistance {
		inner = shrink(inner, config.Config.InnerMaskBorderDistance-total)
	}

	return inner, m.logger.Log("inner-mask", inner)
}
