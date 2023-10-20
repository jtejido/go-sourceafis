package pixelwise

import (
	"math"

	"github.com/jtejido/sourceafis/config"
	"github.com/jtejido/sourceafis/extractor/logger"
	"github.com/jtejido/sourceafis/primitives"
	"github.com/jtejido/sourceafis/utils"
)

type ConsideredOrientation struct {
	offset      primitives.IntPoint
	orientation primitives.FloatPoint
}

type PixelwiseOrientations struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *PixelwiseOrientations {
	return &PixelwiseOrientations{
		logger: logger,
	}
}

func (o *PixelwiseOrientations) Compute(input *primitives.Matrix, mask *primitives.BooleanMatrix, blocks *primitives.BlockMap) *primitives.FloatPointMatrix {
	neighbors := plan()

	orientation := primitives.NewFloatPointMatrixFromPoint(input.Size())
	for blockY := 0; blockY < blocks.Primary.Blocks.Y; blockY++ {
		maskRange := maskRange(mask, blockY)
		if maskRange.Length() > 0 {
			validXRange := primitives.IntRange{
				Start: blocks.Primary.Block(maskRange.Start, blockY).Left(),
				End:   blocks.Primary.Block(maskRange.End-1, blockY).Right(),
			}
			for y := blocks.Primary.Block(0, blockY).Top(); y < blocks.Primary.Block(0, blockY).Bottom(); y++ {
				for _, neighbor := range neighbors[y%len(neighbors)] {
					radius := int(math.Max(math.Abs(float64(neighbor.offset.X)), math.Abs(float64(neighbor.offset.Y))))
					if y-radius >= 0 && y+radius < input.Height {
						xRange := primitives.IntRange{
							Start: int(math.Max(float64(radius), float64(validXRange.Start))),
							End:   int(math.Min(float64(input.Width-radius), float64(validXRange.End))),
						}
						for x := xRange.Start; x < xRange.End; x++ {
							before := input.Get(x-neighbor.offset.X, y-neighbor.offset.Y)
							at := input.Get(x, y)
							after := input.Get(x+neighbor.offset.X, y+neighbor.offset.Y)
							strength := at - math.Max(before, after)
							if strength > 0 {
								orientation.AddFloatPoint(x, y, neighbor.orientation.Multiply(strength))
							}
						}
					}
				}
			}
		}
	}

	o.logger.Log("pixelwise-orientation", orientation)
	return orientation
}

func plan() [][]*ConsideredOrientation {
	random := utils.NewOrientationRandom()
	splits := make([][]*ConsideredOrientation, config.Config.OrientationSplit)
	for i := 0; i < config.Config.OrientationSplit; i++ {
		orientations := make([]*ConsideredOrientation, config.Config.OrientationsChecked)
		for j := 0; j < config.Config.OrientationsChecked; j++ {
			sample := new(ConsideredOrientation)
			orientations[j] = sample
			for {
				angle := primitives.FloatAngle(random.Float64() * primitives.Pi)
				distance := utils.InterpolateExponential(float64(config.Config.MinOrientationRadius), float64(config.Config.MaxOrientationRadius), random.Float64())
				sample.offset = angle.ToVector().Multiply(distance).Round()

				if sample.offset.Equals(primitives.ZeroIntPoint()) {
					continue
				}
				if sample.offset.Y < 0 {
					continue
				}
				var duplicate bool
				for jj := 0; jj < j; jj++ {
					if orientations[jj].offset.Equals(sample.offset) {
						duplicate = true
					}
				}
				if duplicate {
					continue
				}

				break
			}

			sample.orientation = primitives.AngleAdd(float64(primitives.AtanFromFloatPointVector(sample.offset.ToFloat()).ToOrientation()), primitives.Pi).ToVector()
		}
		splits[i] = orientations
	}
	return splits
}

func maskRange(mask *primitives.BooleanMatrix, y int) primitives.IntRange {
	first := -1
	last := -1
	for x := 0; x < mask.Width; x++ {
		if mask.Get(x, y) {
			last = x
			if first < 0 {
				first = x
			}

		}
	}
	if first >= 0 {
		return primitives.IntRange{
			Start: first,
			End:   last + 1,
		}
	}

	return primitives.ZeroIntRange()
}
