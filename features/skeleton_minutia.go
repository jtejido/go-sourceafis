package features

import (
	"github.com/jtejido/sourceafis/primitives"
)

type SkeletonMinutia struct {
	Position primitives.IntPoint
	Ridges   []*SkeletonRidge
}

func NewSkeletonMinutia(position primitives.IntPoint) *SkeletonMinutia {
	return &SkeletonMinutia{
		Position: position,
		Ridges:   make([]*SkeletonRidge, 0),
	}
}

func (m *SkeletonMinutia) AttachStart(ridge *SkeletonRidge) {
	if !m.containsRidge(ridge) {
		m.Ridges = append(m.Ridges, ridge)
		ridge.SetStart(m)
	}
}
func (m *SkeletonMinutia) DetachStart(ridge *SkeletonRidge) {
	if m.containsRidge(ridge) {
		m.removeRidge(ridge)
		if ridge.Start() == m {
			ridge.SetStart(nil)
		}
	}
}

func (minutia *SkeletonMinutia) containsRidge(ridge *SkeletonRidge) bool {
	for _, r := range minutia.Ridges {
		if r == ridge {
			return true
		}
	}
	return false
}

func (minutia *SkeletonMinutia) removeRidge(ridge *SkeletonRidge) {
	for i, r := range minutia.Ridges {
		if r == ridge {
			minutia.Ridges[i] = minutia.Ridges[len(minutia.Ridges)-1]
			minutia.Ridges[len(minutia.Ridges)-1] = nil
			minutia.Ridges = minutia.Ridges[:len(minutia.Ridges)-1]
			break
		}
	}
}
