package features

import (
	"reflect"
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
	for i := 0; i < len(s.Minutiae); i++ {
		if reflect.DeepEqual(s.Minutiae[i], minutia) {
			s.Minutiae = removeMinutia(s.Minutiae, i)
		}
	}
}

func removeMinutia(s []*SkeletonMinutia, i int) []*SkeletonMinutia {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (s *Skeleton) Shadow() (*primitives.BooleanMatrix, error) {
	shadow := primitives.NewBooleanMatrixFromPoint(s.Size)
	for _, minutia := range s.Minutiae {

		shadow.SetPoint(minutia.Position, true)

		for _, ridge := range minutia.Ridges {
			if ridge.StartMinutia().Position.Y <= ridge.EndMinutia().Position.Y {
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
