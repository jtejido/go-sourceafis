package features

import "github.com/jtejido/sourceafis/primitives"

type MinutiaType int

const (
	ENDING MinutiaType = iota
	BIFURCATION
)

type SearchMinutia struct {
	X, Y      int
	Direction float64
	T         MinutiaType
}

func NewSearchMinutia(feature *FeatureMinutia) *SearchMinutia {
	return &SearchMinutia{
		X:         feature.Position.X,
		Y:         feature.Position.Y,
		Direction: feature.Direction,
		T:         feature.T,
	}
}

func (sm *SearchMinutia) Feature() *FeatureMinutia {
	return NewFeatureMinutia(primitives.IntPoint{X: int(sm.X), Y: int(sm.Y)}, sm.Direction, sm.T)
}
