package binarizer

import (
	"sourceafis/config"
	"sourceafis/extractor/filter/vote"
	"sourceafis/extractor/logger"
	"sourceafis/primitives"
)

type BinarizedImage struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *BinarizedImage {
	return &BinarizedImage{
		logger: logger,
	}
}

func (b *BinarizedImage) Binarize(input, baseline *primitives.Matrix, mask *primitives.BooleanMatrix, blocks *primitives.BlockMap) *primitives.BooleanMatrix {
	binarized := primitives.NewBooleanMatrixFromPoint(input.Size())
	it := blocks.Primary.Blocks.Iterator()
	for it.HasNext() {
		block := it.Next()
		if mask.GetPoint(block) {
			rect := blocks.Primary.BlockPoint(block)
			for y := rect.Top(); y < rect.Bottom(); y++ {
				for x := rect.Left(); x < rect.Right(); x++ {
					if input.Get(x, y)-baseline.Get(x, y) > 0 {
						binarized.Set(x, y, true)
					}
				}
			}
		}
	}
	b.logger.Log("binarized-image", binarized)
	return binarized
}

func (b *BinarizedImage) Cleanup(binary, mask *primitives.BooleanMatrix) {
	size := binary.Size()
	inverted := primitives.NewBooleanMatrixFromBooleanMatrix(binary)
	inverted.Invert()
	islands := vote.Apply(inverted, mask, config.Config.BinarizedVoteRadius, config.Config.BinarizedVoteMajority, config.Config.BinarizedVoteBorderDistance)
	holes := vote.Apply(binary, mask, config.Config.BinarizedVoteRadius, config.Config.BinarizedVoteMajority, config.Config.BinarizedVoteBorderDistance)
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			binary.Set(x, y, binary.Get(x, y) && !islands.Get(x, y) || holes.Get(x, y))
		}
	}
	removeCrosses(binary)
	b.logger.Log("filtered-binary-image", binary)
}

func removeCrosses(input *primitives.BooleanMatrix) {
	size := input.Size()
	any := true
	for any {
		any = false
		for y := 0; y < size.Y-1; y++ {
			for x := 0; x < size.X-1; x++ {
				if input.Get(x, y) && input.Get(x+1, y+1) && !input.Get(x, y+1) && !input.Get(x+1, y) || input.Get(x, y+1) && input.Get(x+1, y) && !input.Get(x, y) && !input.Get(x+1, y+1) {
					input.Set(x, y, false)
					input.Set(x, y+1, false)
					input.Set(x+1, y, false)
					input.Set(x+1, y+1, false)
					any = true
				}
			}
		}
	}
}

func (b *BinarizedImage) Invert(binary, mask *primitives.BooleanMatrix) *primitives.BooleanMatrix {
	size := binary.Size()
	inverted := primitives.NewBooleanMatrixFromPoint(size)
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			inverted.Set(x, y, !binary.Get(x, y) && mask.Get(x, y))
		}
	}
	return inverted
}
