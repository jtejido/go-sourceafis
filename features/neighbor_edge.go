package features

import (
	"sort"

	"github.com/jtejido/sourceafis/config"
)

type NeighborEdge struct {
	*EdgeShape
	Neighbor int
}

func NewNeighborEdge(minutiae []*SearchMinutia, reference, neighbor int) *NeighborEdge {
	es := NewEdgeShape(minutiae[reference], minutiae[neighbor])
	return &NeighborEdge{
		EdgeShape: es,
		Neighbor:  neighbor,
	}
}

type EdgeTableLogger interface {
	Log(key string, data interface{}) error
}

type NeighborhoodBuilder struct {
	logger EdgeTableLogger
}

func NewNeighborhoodBuilder(logger EdgeTableLogger) *NeighborhoodBuilder {
	return &NeighborhoodBuilder{
		logger: logger,
	}
}

func (b *NeighborhoodBuilder) Build(minutiae []*SearchMinutia) [][]*NeighborEdge {
	edges := make([][]*NeighborEdge, len(minutiae))
	star := make([]*NeighborEdge, 0)
	allSqDistances := make([]int, len(minutiae))
	for reference := 0; reference < len(edges); reference++ {
		rminutia := minutiae[reference]
		maxSqDistance := int(^uint(0) >> 1)
		if len(minutiae)-1 > config.Config.EdgeTableNeighbors {
			for neighbor := 0; neighbor < len(minutiae); neighbor++ {
				nminutia := minutiae[neighbor]
				allSqDistances[neighbor] = (rminutia.X-nminutia.X)*(rminutia.X-nminutia.X) + (rminutia.Y-nminutia.Y)*(rminutia.Y-nminutia.Y)
			}

			sort.Ints(allSqDistances)
			maxSqDistance = allSqDistances[config.Config.EdgeTableNeighbors]
		}
		for neighbor := 0; neighbor < len(minutiae); neighbor++ {
			nminutia := minutiae[neighbor]
			distanceSq := (rminutia.X-nminutia.X)*(rminutia.X-nminutia.X) + (rminutia.Y-nminutia.Y)*(rminutia.Y-nminutia.Y)
			if neighbor != reference && distanceSq <= maxSqDistance {
				star = append(star, NewNeighborEdge(minutiae, reference, neighbor))
			}
		}

		sort.Slice(star, func(i, j int) bool {
			if star[i].Length == star[j].Length {
				return star[i].Neighbor < star[j].Neighbor
			}
			return star[i].Length < star[j].Length
		})

		for len(star) > config.Config.EdgeTableNeighbors {
			star = star[:len(star)-1]
		}

		edges[reference] = make([]*NeighborEdge, len(star))
		copy(edges[reference], star)
		star = star[:0]
	}
	b.logger.Log("edge-table", edges)
	return edges
}
