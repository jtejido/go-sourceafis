package thinner

import (
	"math/bits"
	"sourceafis/config"
	"sourceafis/extractor/logger"
	"sourceafis/features"
	"sourceafis/primitives"
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

func (t *BinaryThinning) Thin(input *primitives.BooleanMatrix, ty features.SkeletonType) *primitives.BooleanMatrix {
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
							} else {
								neighbors |= 0
							}

							if partial.Get(x-1, y+1) {
								neighbors |= 32
							} else {
								neighbors |= 0
							}

							if partial.Get(x+1, y) {
								neighbors |= 16
							} else {
								neighbors |= 0
							}

							if partial.Get(x-1, y) {
								neighbors |= 8
							} else {
								neighbors |= 0
							}

							if partial.Get(x+1, y-1) {
								neighbors |= 4
							} else {
								neighbors |= 0
							}

							if partial.Get(x, y-1) {
								neighbors |= 2
							} else {
								neighbors |= 0
							}

							if partial.Get(x-1, y-1) {
								neighbors |= 1
							} else {
								neighbors |= 0
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
	t.logger.Log(ty.String()+"thinned-skeleton", thinned)
	return thinned
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
			return count > 2
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