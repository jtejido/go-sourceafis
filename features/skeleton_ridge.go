package features

import (
	"sourceafis/config"
	"sourceafis/primitives"
)

type SkeletonRidge struct {
	Points                   primitives.List[primitives.IntPoint]
	Reversed                 *SkeletonRidge
	startMinutia, endMinutia *SkeletonMinutia
}

func NewSkeletonRidge() *SkeletonRidge {
	s := new(SkeletonRidge)
	s.Points = primitives.NewCircularList[primitives.IntPoint]()
	s.Reversed = NewSkeletonRidgeFromReversed(s)
	return s
}
func NewSkeletonRidgeFromReversed(reversed *SkeletonRidge) *SkeletonRidge {
	s := new(SkeletonRidge)
	s.Points = primitives.NewReversedList(reversed.Points)
	s.Reversed = reversed
	return s
}

func (r *SkeletonRidge) Start() *SkeletonMinutia {
	return r.startMinutia
}

func (r *SkeletonRidge) SetStart(value *SkeletonMinutia) {
	if r.startMinutia != value {
		if r.startMinutia != nil {
			detachFrom := r.startMinutia
			r.startMinutia = nil
			detachFrom.DetachStart(r)
		}
		r.startMinutia = value
		if r.startMinutia != nil {
			r.startMinutia.AttachStart(r)
		}
		r.Reversed.endMinutia = value
	}
}

func (r *SkeletonRidge) End() *SkeletonMinutia {
	return r.endMinutia
}
func (r *SkeletonRidge) SetEnd(value *SkeletonMinutia) {
	if r.endMinutia != value {
		r.endMinutia = value
		r.Reversed.SetStart(value)
	}
}

func (r *SkeletonRidge) Detach() {
	r.SetStart(nil)
	r.SetEnd(nil)
}

func (r *SkeletonRidge) Direction() (float64, error) {
	first := config.Config.RidgeDirectionSkip
	last := config.Config.RidgeDirectionSkip + config.Config.RidgeDirectionSample - 1
	if last >= r.Points.Size() {
		shift := last - r.Points.Size() + 1
		last -= shift
		first -= shift
	}
	if first < 0 {
		first = 0
	}
	center, err := r.Points.Get(first)
	if err != nil {
		return 0, err
	}
	point, err := r.Points.Get(last)
	if err != nil {
		return 0, err
	}
	return float64(primitives.AtanFromIntPoints(center, point)), nil
}
