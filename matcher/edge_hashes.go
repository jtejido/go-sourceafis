package matcher

import (
	"math"
	"sourceafis/config"
	"sourceafis/features"
	"sourceafis/primitives"
	"sourceafis/templates"
)

func ComplementaryMaxAngleError() float64 {
	return float64(primitives.FloatAngle(config.Config.MaxAngleError).Complementary())
}

func Hash(edge *features.EdgeShape) int {
	lengthBin := edge.Length / config.Config.MaxDistanceError
	referenceAngleBin := int(float64(edge.ReferenceAngle) / config.Config.MaxAngleError)
	neighborAngleBin := int(float64(edge.NeighborAngle) / config.Config.MaxAngleError)
	return (referenceAngleBin << 24) + (neighborAngleBin << 16) + lengthBin
}

func Matching(probe *features.IndexedEdge, candidate *features.EdgeShape) bool {
	lengthDelta := probe.Length - candidate.Length
	if lengthDelta >= -config.Config.MaxDistanceError && lengthDelta <= config.Config.MaxDistanceError {
		referenceDelta := float64(probe.ReferenceAngle.Difference(candidate.ReferenceAngle))
		if referenceDelta <= config.Config.MaxAngleError || referenceDelta >= ComplementaryMaxAngleError() {
			neighborDelta := float64(probe.NeighborAngle.Difference(candidate.NeighborAngle))
			if neighborDelta <= config.Config.MaxAngleError || neighborDelta >= ComplementaryMaxAngleError() {
				return true
			}
		}
	}
	return false
}

func coverage(edge *features.IndexedEdge) []int {
	minLengthBin := (edge.Length - config.Config.MaxDistanceError) / config.Config.MaxDistanceError
	maxLengthBin := (edge.Length + config.Config.MaxDistanceError) / config.Config.MaxDistanceError
	angleBins := int(math.Ceil(2 * primitives.Pi / config.Config.MaxAngleError))
	minReferenceBin := int(math.Floor(float64(edge.ReferenceAngle.Difference(primitives.FloatAngle(config.Config.MaxAngleError))) / config.Config.MaxAngleError))
	maxReferenceBin := int(math.Floor(float64(primitives.AngleAdd(float64(edge.ReferenceAngle), config.Config.MaxAngleError)) / config.Config.MaxAngleError))
	endReferenceBin := (maxReferenceBin + 1) % angleBins
	minNeighborBin := int(math.Floor(float64(edge.NeighborAngle.Difference(primitives.FloatAngle(config.Config.MaxAngleError))) / config.Config.MaxAngleError))
	maxNeighborBin := int(math.Floor(float64(primitives.AngleAdd(float64(edge.NeighborAngle), config.Config.MaxAngleError)) / config.Config.MaxAngleError))
	endNeighborBin := (maxNeighborBin + 1) % angleBins
	coverage := make([]int, 0)
	for lengthBin := minLengthBin; lengthBin <= maxLengthBin; lengthBin++ {
		for referenceBin := minReferenceBin; referenceBin != endReferenceBin; referenceBin = (referenceBin + 1) % angleBins {
			for neighborBin := minNeighborBin; neighborBin != endNeighborBin; neighborBin = (neighborBin + 1) % angleBins {
				coverage = append(coverage, (referenceBin<<24)+(neighborBin<<16)+lengthBin)
			}
		}
	}
	return coverage
}

type HashTableLogger interface {
	LogEdgeHash(map[int][]*features.IndexedEdge) error
}

type EdgeHashBuilder struct {
	logger HashTableLogger
}

func NewEdgeHashBuilder(logger HashTableLogger) *EdgeHashBuilder {
	return &EdgeHashBuilder{
		logger: logger,
	}
}

func (b *EdgeHashBuilder) Build(template *templates.SearchTemplate) map[int][]*features.IndexedEdge {
	m := make(map[int][]*features.IndexedEdge)
	for reference := 0; reference < len(template.Minutiae); reference++ {
		for neighbor := 0; neighbor < len(template.Minutiae); neighbor++ {
			if reference != neighbor {
				edge := features.NewIndexedEdge(template.Minutiae, reference, neighbor)
				for _, hash := range coverage(edge) {
					if m[hash] == nil {
						m[hash] = make([]*features.IndexedEdge, 0)
					}
					m[hash] = append(m[hash], edge)
				}
			}
		}
	}
	b.logger.LogEdgeHash(m)
	return m
}
