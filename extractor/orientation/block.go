package orientation

import (
	"sourceafis/config"
	"sourceafis/extractor/logger"
	"sourceafis/extractor/orientation/pixelwise"
	"sourceafis/primitives"
)

type BlockOrientations struct {
	logger    logger.TransparencyLogger
	pixelwise *pixelwise.PixelwiseOrientations
}

func New(logger logger.TransparencyLogger) *BlockOrientations {
	return &BlockOrientations{
		logger:    logger,
		pixelwise: pixelwise.New(logger),
	}
}

func (o *BlockOrientations) Compute(image *primitives.Matrix, mask *primitives.BooleanMatrix, blocks *primitives.BlockMap) (*primitives.Matrix, error) {
	accumulated := o.pixelwise.Compute(image, mask, blocks)

	byBlock, err := o.aggregate(accumulated, blocks, mask)
	if err != nil {
		return nil, err
	}

	smooth, err := o.smooth(byBlock, mask)
	if err != nil {
		return nil, err
	}

	return angles(smooth, mask), nil
}

func (o *BlockOrientations) aggregate(orientation *primitives.FloatPointMatrix, blocks *primitives.BlockMap, mask *primitives.BooleanMatrix) (*primitives.FloatPointMatrix, error) {
	sums := primitives.NewFloatPointMatrixFromPoint(blocks.Primary.Blocks)
	it := blocks.Primary.Blocks.Iterator()
	for it.HasNext() {
		block := it.Next()
		if mask.GetPoint(block) {
			area := blocks.Primary.BlockPoint(block)
			for y := area.Top(); y < area.Bottom(); y++ {
				for x := area.Left(); x < area.Right(); x++ {
					sums.AddPoint(block, orientation.Get(x, y))
				}
			}
		}
	}

	return sums, o.logger.Log("block-orientation", sums)
}

func (o *BlockOrientations) smooth(orientation *primitives.FloatPointMatrix, mask *primitives.BooleanMatrix) (*primitives.FloatPointMatrix, error) {
	size := mask.Size()
	smoothed := primitives.NewFloatPointMatrixFromPoint(size)
	it := size.Iterator()
	for it.HasNext() {
		block := it.Next()
		if mask.GetPoint(block) {
			neighbors := primitives.IntRectAroundIntPoint(block, config.Config.OrientationSmoothingRadius).Intersect(primitives.IntRectFromPoint(size))
			for ny := neighbors.Top(); ny < neighbors.Bottom(); ny++ {
				for nx := neighbors.Left(); nx < neighbors.Right(); nx++ {
					if mask.Get(nx, ny) {
						smoothed.AddPoint(block, orientation.Get(nx, ny))
					}
				}
			}
		}
	}

	return smoothed, o.logger.Log("smoothed-orientation", smoothed)
}

func angles(vectors *primitives.FloatPointMatrix, mask *primitives.BooleanMatrix) *primitives.Matrix {
	size := mask.Size()
	angles := primitives.NewMatrixFromPoint(size)
	it := size.Iterator()
	for it.HasNext() {
		block := it.Next()
		if mask.GetPoint(block) {
			angles.SetPoint(block, float64(primitives.AtanFromFloatPointVector(vectors.GetPoint(block))))
		}

	}
	return angles
}
