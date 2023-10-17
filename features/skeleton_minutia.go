package features

import (
	"reflect"
	"sourceafis/primitives"
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
	for _, r := range m.Ridges {
		if ridge == r {
			return
		}
	}

	m.Ridges = append(m.Ridges, ridge)
	ridge.Start(m)
}
func (m *SkeletonMinutia) DetachStart(ridge *SkeletonRidge) {
	var contains bool
	var index int
	for i, r := range m.Ridges {
		if r == ridge {
			contains = true
			index = i
			break
		}
	}

	if contains {
		//m.Ridges = remove(m.Ridges, index)
		m.Ridges[index] = m.Ridges[len(m.Ridges)-1]
		m.Ridges[len(m.Ridges)-1] = nil
		m.Ridges = m.Ridges[:len(m.Ridges)-1]
		if reflect.DeepEqual(ridge.StartMinutia(), m) {
			ridge.Start(nil)
		}
	}
}

func remove(slice []*SkeletonRidge, s int) []*SkeletonRidge {
	return append(slice[:s], slice[s+1:]...)
}
