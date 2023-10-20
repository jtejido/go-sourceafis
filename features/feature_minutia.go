package features

import "github.com/jtejido/sourceafis/primitives"

type FeatureMinutia struct {
	Position  primitives.IntPoint
	Direction float64
	T         MinutiaType
}

func NewFeatureMinutia(position primitives.IntPoint, direction float64, t MinutiaType) *FeatureMinutia {
	return &FeatureMinutia{
		Position:  position,
		Direction: direction,
		T:         t,
	}
}
