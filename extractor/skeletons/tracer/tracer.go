package tracer

import (
	"sort"
	"sourceafis/extractor/logger"
	"sourceafis/features"
	"sourceafis/primitives"

	"golang.org/x/exp/slices"
)

type SkeletonTracing struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *SkeletonTracing {
	return &SkeletonTracing{
		logger: logger,
	}
}

func (tr *SkeletonTracing) Trace(thinned *primitives.BooleanMatrix, t features.SkeletonType) (*features.Skeleton, error) {
	skeleton := features.NewSkeleton(t, thinned.Size())
	minutiaPoints := findMinutiae(thinned)
	linking := linkNeighboringMinutiae(minutiaPoints)
	minutiaMap := minutiaCenters(skeleton, linking)
	err := traceRidges(thinned, minutiaMap)
	if err != nil {
		return nil, err
	}
	err = fixLinkingGaps(skeleton)
	if err != nil {
		return nil, err
	}

	return skeleton, tr.logger.LogSkeleton("traced-skeleton", skeleton)
}

func findMinutiae(thinned *primitives.BooleanMatrix) []primitives.IntPoint {
	result := make([]primitives.IntPoint, 0)
	it := thinned.Size().Iterator()
	for it.HasNext() {
		at := it.Next()
		if thinned.GetPoint(at) {
			var count int
			for _, relative := range primitives.CORNER_NEIGHBORS {
				if thinned.GetPointWithFallback(at.Plus(relative), false) {
					count++
				}
			}
			if count == 1 || count > 2 {
				result = append(result, at)
			}
		}
	}
	return result
}

func eq(a, b []primitives.IntPoint) bool {
	if len(a) != len(b) {
		return false
	}

	return slices.Equal(a, b)
}

func linkNeighboringMinutiae(minutiae []primitives.IntPoint) map[primitives.IntPoint][]primitives.IntPoint {
	linking := make(map[primitives.IntPoint][]primitives.IntPoint)
	for _, minutiaPos := range minutiae {
		var ownLinks []primitives.IntPoint
		for _, neighborRelative := range primitives.CORNER_NEIGHBORS {
			neighborPos := minutiaPos.Plus(neighborRelative)
			if neighborLinks, ok := linking[neighborPos]; ok {
				if !eq(neighborLinks, ownLinks) {
					if ownLinks != nil {
						neighborLinks = append(neighborLinks, ownLinks...)
						for _, mergedPos := range ownLinks {
							linking[mergedPos] = neighborLinks
						}
					}
					ownLinks = neighborLinks
				}
			}
		}
		if ownLinks == nil {
			ownLinks = make([]primitives.IntPoint, 0)
		}
		ownLinks = append(ownLinks, minutiaPos)
		linking[minutiaPos] = ownLinks
	}
	return linking
}

type list []primitives.IntPoint

func (e list) Len() int {
	return len(e)
}

func (e list) Less(i, j int) bool {
	return e[i].CompareTo(e[j]) < 0
}

func (e list) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func minutiaCenters(skeleton *features.Skeleton, linking map[primitives.IntPoint][]primitives.IntPoint) map[primitives.IntPoint]*features.SkeletonMinutia {
	centers := make(map[primitives.IntPoint]*features.SkeletonMinutia)
	keys := make([]primitives.IntPoint, 0, len(linking))
	for k := range linking {
		keys = append(keys, k)
	}
	sort.Sort(list(keys))
	for _, currentPos := range keys {
		linkedMinutiae := linking[currentPos]
		primaryPos := linkedMinutiae[0]
		if _, ok := centers[primaryPos]; !ok {
			sum := primitives.ZeroIntPoint()
			for _, linkedPos := range linkedMinutiae {
				sum = sum.Plus(linkedPos)
			}
			center := primitives.IntPoint{
				X: sum.X / len(linkedMinutiae),
				Y: sum.Y / len(linkedMinutiae),
			}
			minutia := features.NewSkeletonMinutia(center)
			skeleton.AddMinutia(minutia)
			centers[primaryPos] = minutia
		}
		centers[currentPos] = centers[primaryPos]
	}
	return centers
}

func fixLinkingGaps(skeleton *features.Skeleton) error {
	for _, minutia := range skeleton.Minutiae {
		for _, ridge := range minutia.Ridges {
			v, err := ridge.Points.Get(0)
			if err != nil {
				return err
			}
			if !v.Equals(minutia.Position) {
				filling := v.LineTo(minutia.Position)
				for i := 1; i < len(filling); i++ {
					ridge.Reversed.Points.Add(filling[i])
				}
			}
		}
	}
	return nil
}

func traceRidges(thinned *primitives.BooleanMatrix, minutiaePoints map[primitives.IntPoint]*features.SkeletonMinutia) error {
	leads := make(map[primitives.IntPoint]*features.SkeletonRidge)
	keys := make([]primitives.IntPoint, 0, len(minutiaePoints))
	for k := range minutiaePoints {
		keys = append(keys, k)
	}
	sort.Sort(list(keys))
	for _, minutiaPoint := range keys {
		for _, startRelative := range primitives.CORNER_NEIGHBORS {
			start := minutiaPoint.Plus(startRelative)
			if _, ok := minutiaePoints[start]; thinned.GetPointWithFallback(start, false) && !ok && leads[start] == nil {
				ridge := features.NewSkeletonRidge()
				ridge.Points.Add(minutiaPoint)
				ridge.Points.Add(start)
				previous := minutiaPoint
				current := start
				for ok := true; ok; ok = minutiaePoints[current] == nil {
					next := primitives.ZeroIntPoint()
					for _, nextRelative := range primitives.CORNER_NEIGHBORS {
						next = current.Plus(nextRelative)
						if thinned.GetPointWithFallback(next, false) && !next.Equals(previous) {
							break
						}
					}
					previous = current
					current = next
					ridge.Points.Add(current)
				}
				end := current
				ridge.SetStart(minutiaePoints[minutiaPoint])
				ridge.SetEnd(minutiaePoints[end])
				v, err := ridge.Points.Get(1)
				if err != nil {
					return err
				}
				leads[v] = ridge
				v, err = ridge.Reversed.Points.Get(1)
				if err != nil {
					return err
				}
				leads[v] = ridge
			}
		}
	}

	return nil
}
