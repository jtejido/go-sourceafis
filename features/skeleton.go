package features

import (
	"sourceafis/primitives"
)

type SkeletonType int

const (
	RIDGES SkeletonType = iota
	VALLEYS
)

func (t SkeletonType) String() string {
	switch t {
	case RIDGES:
		return "ridges-"
	case VALLEYS:
		return "valleys-"
	default:
		return ""
	}
}

type Skeleton struct {
	T        SkeletonType
	Size     primitives.IntPoint
	Minutiae []*SkeletonMinutia
}

func NewSkeleton(t SkeletonType, size primitives.IntPoint) *Skeleton {
	return &Skeleton{
		T:        t,
		Size:     size,
		Minutiae: make([]*SkeletonMinutia, 0),
	}
}

func (s *Skeleton) AddMinutia(minutia *SkeletonMinutia) {
	s.Minutiae = append(s.Minutiae, minutia)
}
func (s *Skeleton) RemoveMinutia(minutia *SkeletonMinutia) {
	for i, m := range s.Minutiae {
		if m == minutia {
			s.Minutiae = append(s.Minutiae[:i], s.Minutiae[i+1:]...)
			break
		}
	}
}

func (s *Skeleton) Shadow() (*primitives.BooleanMatrix, error) {
	shadow := primitives.NewBooleanMatrixFromPoint(s.Size)
	for _, minutia := range s.Minutiae {
		shadow.SetPoint(minutia.Position, true)

		for _, ridge := range minutia.Ridges {
			if ridge.Start().Position.Y <= ridge.End().Position.Y {
				it := ridge.Points.Iterator()
				for it.HasNext() {
					point, err := it.Next()
					if err != nil {
						return nil, err
					}
					shadow.SetPoint(point, true)

				}
			}
		}
	}
	return shadow, nil
}
