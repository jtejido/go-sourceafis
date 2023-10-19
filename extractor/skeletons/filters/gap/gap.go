package gap

import (
	"container/heap"
	"sourceafis/config"
	"sourceafis/extractor/logger"
	"sourceafis/extractor/skeletons/filters/knot"
	"sourceafis/features"
	"sourceafis/primitives"
	"sourceafis/utils"
)

type SkeletonGap struct {
	distance   int
	end1, end2 *features.SkeletonMinutia
}

func (gap *SkeletonGap) CompareTo(other *SkeletonGap) int {
	distanceCmp := gap.distance - other.distance
	if distanceCmp != 0 {
		return distanceCmp
	}
	end1Cmp := gap.end1.Position.CompareTo(other.end1.Position)
	if end1Cmp != 0 {
		return end1Cmp
	}
	return gap.end2.Position.CompareTo(other.end2.Position)
}

type PriorityQueue []*SkeletonGap

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].CompareTo(pq[j]) < 0
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]

}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*SkeletonGap)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

type SkeletonGapFilter struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *SkeletonGapFilter {
	return &SkeletonGapFilter{
		logger: logger,
	}
}

func (f *SkeletonGapFilter) Apply(skeleton *features.Skeleton) error {
	queue := make(PriorityQueue, 0)
	heap.Init(&queue)

	for _, end1 := range skeleton.Minutiae {
		if len(end1.Ridges) == 1 && end1.Ridges[0].Points.Size() >= config.Config.ShortestJoinedEnding {
			for _, end2 := range skeleton.Minutiae {
				withinGapLimits, err := isWithinGapLimits(end1, end2)
				if err != nil {
					return err
				}
				if end2 != end1 && len(end2.Ridges) == 1 && end1.Ridges[0].End() != end2 && end2.Ridges[0].Points.Size() >= config.Config.ShortestJoinedEnding && withinGapLimits {
					gap := new(SkeletonGap)
					gap.distance = end1.Position.Minus(end2.Position).LengthSq()
					gap.end1 = end1
					gap.end2 = end2
					heap.Push(&queue, gap)
				}
			}
		}
	}

	shadow, err := skeleton.Shadow()
	if err != nil {
		return err
	}
	for queue.Len() > 0 {
		gap := heap.Pop(&queue).(*SkeletonGap)
		if len(gap.end1.Ridges) == 1 && len(gap.end2.Ridges) == 1 {
			line := gap.end1.Position.LineTo(gap.end2.Position)
			if !isRidgeOverlapping(line, shadow) {
				addGapRidge(shadow, gap, line)
			}
		}
	}

	if err := knot.Apply(skeleton); err != nil {
		return err
	}
	f.logger.LogSkeleton("removed-gaps", skeleton)
	return nil
}

func addGapRidge(shadow *primitives.BooleanMatrix, gap *SkeletonGap, line []primitives.IntPoint) {
	ridge := features.NewSkeletonRidge()
	for _, point := range line {
		ridge.Points.Add(point)
	}
	ridge.SetStart(gap.end1)
	ridge.SetEnd(gap.end2)
	for _, point := range line {
		shadow.SetPoint(point, true)
	}
}

func isRidgeOverlapping(line []primitives.IntPoint, shadow *primitives.BooleanMatrix) bool {
	for i := config.Config.ToleratedGapOverlap; i < len(line)-config.Config.ToleratedGapOverlap; i++ {
		if shadow.GetPoint(line[i]) {
			return true
		}
	}
	return false
}

func isWithinGapLimits(end1, end2 *features.SkeletonMinutia) (bool, error) {
	distanceSq := end1.Position.Minus(end2.Position).LengthSq()
	if distanceSq <= utils.SquareInt(config.Config.MaxRuptureSize) {
		return true, nil
	}
	if distanceSq > utils.SquareInt(config.Config.MaxGapSize) {
		return false, nil
	}
	gapDirection := primitives.AtanFromIntPoints(end1.Position, end2.Position)
	angle, err := angleSampleForGapRemoval(end1)
	if err != nil {
		return false, err
	}
	direction1 := primitives.AtanFromIntPoints(end1.Position, angle)
	if direction1.Distance(gapDirection.Opposite()) > config.Config.MaxGapAngle {
		return false, nil
	}
	angle2, err := angleSampleForGapRemoval(end1)
	if err != nil {
		return false, err
	}
	direction2 := primitives.AtanFromIntPoints(end2.Position, angle2)
	if direction2.Distance(gapDirection) > config.Config.MaxGapAngle {
		return false, nil
	}
	return true, nil
}

func angleSampleForGapRemoval(minutia *features.SkeletonMinutia) (primitives.IntPoint, error) {
	ridge := minutia.Ridges[0]
	if config.Config.GapAngleOffset < ridge.Points.Size() {
		return ridge.Points.Get(config.Config.GapAngleOffset)
	}

	return ridge.End().Position, nil
}
