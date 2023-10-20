package thinner

import (
	"math/bits"

	"github.com/jtejido/sourceafis/config"
	"github.com/jtejido/sourceafis/extractor/logger"
	"github.com/jtejido/sourceafis/features"
	"github.com/jtejido/sourceafis/primitives"
)

type NeighborhoodType int

const (
	SKELETON NeighborhoodType = iota
	ENDING
	REMOVABLE
)

type BinaryThinning struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *BinaryThinning {
	return &BinaryThinning{
		logger: logger,
	}
}

func (t *BinaryThinning) Thin(input *primitives.BooleanMatrix, ty features.SkeletonType) (*primitives.BooleanMatrix, error) {
	neighborhoodTypes := neighborhoodTypes()
	size := input.Size()
	partial := primitives.NewBooleanMatrixFromPoint(size)
	for y := 1; y < size.Y-1; y++ {
		for x := 1; x < size.X-1; x++ {
			partial.Set(x, y, input.Get(x, y))
		}
	}
	thinned := primitives.NewBooleanMatrixFromPoint(size)
	removedAnything := true
	for i := 0; i < config.Config.ThinningIterations && removedAnything; i++ {
		removedAnything = false
		for evenY := 0; evenY < 2; evenY++ {
			for evenX := 0; evenX < 2; evenX++ {
				for y := 1 + evenY; y < size.Y-1; y += 2 {
					for x := 1 + evenX; x < size.X-1; x += 2 {
						if partial.Get(x, y) && !thinned.Get(x, y) && !(partial.Get(x, y-1) && partial.Get(x, y+1) && partial.Get(x-1, y) && partial.Get(x+1, y)) {
							var neighbors int
							if partial.Get(x+1, y+1) {
								neighbors = 128
							}
							if partial.Get(x, y+1) {
								neighbors |= 64
							}

							if partial.Get(x-1, y+1) {
								neighbors |= 32
							}

							if partial.Get(x+1, y) {
								neighbors |= 16
							}

							if partial.Get(x-1, y) {
								neighbors |= 8
							}

							if partial.Get(x+1, y-1) {
								neighbors |= 4
							}

							if partial.Get(x, y-1) {
								neighbors |= 2
							}

							if partial.Get(x-1, y-1) {
								neighbors |= 1
							}

							if (neighborhoodTypes[neighbors] == REMOVABLE || neighborhoodTypes[neighbors] == ENDING && isFalseEnding(partial, primitives.IntPoint{X: x, Y: y})) {
								removedAnything = true
								partial.Set(x, y, false)
							} else {
								thinned.Set(x, y, true)
							}
						}
					}
				}
			}
		}
	}

	return thinned, t.logger.Log(ty.String()+"thinned-skeleton", thinned)
}

func isFalseEnding(binary *primitives.BooleanMatrix, ending primitives.IntPoint) bool {
	for _, relativeNeighbor := range primitives.CORNER_NEIGHBORS {
		neighbor := ending.Plus(relativeNeighbor)
		if binary.GetPoint(neighbor) {
			var count int
			for _, relative2 := range primitives.CORNER_NEIGHBORS {
				if binary.GetPointWithFallback(neighbor.Plus(relative2), false) {
					count++
				}
			}
			if count > 2 {
				return true
			}
		}
	}
	return false
}

func neighborhoodTypes() []NeighborhoodType {
	types := make([]NeighborhoodType, 256)
	for mask := 0; mask < 256; mask++ {
		TL := (mask & 1) != 0
		TC := (mask & 2) != 0
		TR := (mask & 4) != 0
		CL := (mask & 8) != 0
		CR := (mask & 16) != 0
		BL := (mask & 32) != 0
		BC := (mask & 64) != 0
		BR := (mask & 128) != 0
		count := bits.OnesCount(uint(mask))
		diagonal := !TC && !CL && TL || !CL && !BC && BL || !BC && !CR && BR || !CR && !TC && TR
		horizontal := !TC && !BC && (TR || CR || BR) && (TL || CL || BL)
		vertical := !CL && !CR && (TL || TC || TR) && (BL || BC || BR)
		end := (count == 1)
		if end {
			types[mask] = ENDING
		} else if !diagonal && !horizontal && !vertical {
			types[mask] = REMOVABLE
		} else {
			types[mask] = SKELETON
		}
	}
	return types
}
