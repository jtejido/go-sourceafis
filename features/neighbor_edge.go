package features

import (
	"math"
	"sort"
	"sourceafis/config"
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
		maxSqDistance := math.MaxInt
		if len(minutiae)-1 > config.Config.EdgeTableNeighbors {
			for neighbor := 0; neighbor < len(minutiae); neighbor++ {
				nminutia := minutiae[neighbor]
				allSqDistances[neighbor] = int(math.Pow(float64(rminutia.X-nminutia.X), 2)) + int(math.Pow(float64(rminutia.Y-nminutia.Y), 2))
			}

			sort.Ints(allSqDistances)
			maxSqDistance = allSqDistances[config.Config.EdgeTableNeighbors]
		}
		for neighbor := 0; neighbor < len(minutiae); neighbor++ {
			nminutia := minutiae[neighbor]
			if neighbor != reference && math.Pow(float64(rminutia.X-nminutia.X), 2)+math.Pow(float64(rminutia.Y-nminutia.Y), 2) <= float64(maxSqDistance) {
				star = append(star, NewNeighborEdge(minutiae, reference, neighbor))
			}
		}

		sort.Slice(star, func(i, j int) bool {
			return star[i].Length < star[j].Length || star[i].Neighbor < star[j].Neighbor
		})
		for len(star) > config.Config.EdgeTableNeighbors {
			star[config.Config.EdgeTableNeighbors] = star[len(star)-1]
			star[len(star)-1] = nil
			star = star[:len(star)-1]
		}

		edges[reference] = append(edges[reference], star...)
		star = make([]*NeighborEdge, 0)
	}
	b.logger.Log("edge-table", edges)
	return edges
}
