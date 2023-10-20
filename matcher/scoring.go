package matcher

import (
	"math"

	"github.com/jtejido/sourceafis/config"
	"github.com/jtejido/sourceafis/features"
	"github.com/jtejido/sourceafis/templates"
)

func Compute(probe, candidate *templates.SearchTemplate, pairing *PairingGraph, score *ScoringData) {
	pminutiae := probe.Minutiae
	cminutiae := candidate.Minutiae
	score.MinutiaCount = pairing.Count
	score.MinutiaScore = config.Config.MinutiaScore * float64(score.MinutiaCount)
	score.MinutiaFractionInProbe = float64(pairing.Count) / float64(len(pminutiae))
	score.MinutiaFractionInCandidate = float64(pairing.Count) / float64(len(cminutiae))
	score.MinutiaFraction = 0.5 * (score.MinutiaFractionInProbe + score.MinutiaFractionInCandidate)
	score.MinutiaFractionScore = config.Config.MinutiaFractionScore * score.MinutiaFraction
	score.SupportingEdgeSum = 0
	score.SupportedMinutiaCount = 0
	score.MinutiaTypeHits = 0
	for i := 0; i < pairing.Count; i++ {
		pair := pairing.Tree[i]
		score.SupportingEdgeSum += pair.supportingEdges
		if pair.supportingEdges >= config.Config.MinSupportingEdges {
			score.SupportedMinutiaCount++
		}
		if pminutiae[pair.Probe].T == cminutiae[pair.Candidate].T {
			score.MinutiaTypeHits++
		}
	}
	score.EdgeCount = pairing.Count + score.SupportingEdgeSum
	score.EdgeScore = config.Config.EdgeScore * float64(score.EdgeCount)
	score.SupportedMinutiaScore = config.Config.SupportedMinutiaScore * float64(score.SupportedMinutiaCount)
	score.MinutiaTypeScore = config.Config.MinutiaTypeScore * float64(score.MinutiaTypeHits)
	innerDistanceRadius := math.Floor((config.Config.DistanceErrorFlatness * float64(config.Config.MaxDistanceError)) + .5)
	innerAngleRadius := (config.Config.AngleErrorFlatness * config.Config.MaxAngleError)
	score.DistanceErrorSum = 0
	score.AngleErrorSum = 0

	for i := 1; i < pairing.Count; i++ {
		pair := pairing.Tree[i]
		probeEdge := features.NewEdgeShape(pminutiae[pair.ProbeRef], pminutiae[pair.Probe])
		candidateEdge := features.NewEdgeShape(cminutiae[pair.CandidateRef], cminutiae[pair.Candidate])
		score.DistanceErrorSum += int(math.Max(innerDistanceRadius, math.Abs(float64(probeEdge.Length)-float64(candidateEdge.Length))))
		score.AngleErrorSum += math.Max(innerAngleRadius, probeEdge.ReferenceAngle.Distance(candidateEdge.ReferenceAngle))
		score.AngleErrorSum += math.Max(innerAngleRadius, probeEdge.NeighborAngle.Distance(candidateEdge.NeighborAngle))
	}
	score.DistanceAccuracyScore = 0
	score.AngleAccuracyScore = 0
	distanceErrorPotential := config.Config.MaxDistanceError * int(math.Max(0, float64(pairing.Count-1)))
	score.DistanceAccuracySum = distanceErrorPotential - score.DistanceErrorSum
	if distanceErrorPotential > 0 {
		score.DistanceAccuracyScore = config.Config.DistanceAccuracyScore * (float64(score.DistanceAccuracySum) / float64(distanceErrorPotential))
	} else {
		score.DistanceAccuracyScore = 0
	}
	angleErrorPotential := config.Config.MaxAngleError * math.Max(0, float64(pairing.Count-1)) * 2
	score.AngleAccuracySum = angleErrorPotential - score.AngleErrorSum
	if angleErrorPotential > 0 {
		score.AngleAccuracyScore = config.Config.AngleAccuracyScore * (score.AngleAccuracySum / angleErrorPotential)
	} else {
		score.AngleAccuracyScore = 0
	}

	score.TotalScore = score.MinutiaScore + score.MinutiaFractionScore + score.SupportedMinutiaScore + score.EdgeScore + score.MinutiaTypeScore + score.DistanceAccuracyScore + score.AngleAccuracyScore

	score.ShapedScore = shape(score.TotalScore)
}

func shape(raw float64) float64 {
	if raw < config.Config.ThresholdFMRMax {
		return 0
	}

	if raw < config.Config.ThresholdFMR2 {
		return interpolate(raw, config.Config.ThresholdFMRMax, config.Config.ThresholdFMR2, 0, 3)
	}
	if raw < config.Config.ThresholdFMR10 {
		return interpolate(raw, config.Config.ThresholdFMR2, config.Config.ThresholdFMR10, 3, 7)
	}
	if raw < config.Config.ThresholdFMR100 {
		return interpolate(raw, config.Config.ThresholdFMR10, config.Config.ThresholdFMR100, 10, 10)
	}
	if raw < config.Config.ThresholdFMR1000 {
		return interpolate(raw, config.Config.ThresholdFMR100, config.Config.ThresholdFMR1000, 20, 10)
	}
	if raw < config.Config.ThresholdFMR10000 {
		return interpolate(raw, config.Config.ThresholdFMR1000, config.Config.ThresholdFMR10000, 30, 10)
	}
	if raw < config.Config.ThresholdFMR100000 {
		return interpolate(raw, config.Config.ThresholdFMR10000, config.Config.ThresholdFMR100000, 40, 10)
	}
	return (raw-config.Config.ThresholdFMR100000)/(config.Config.ThresholdFMR100000-config.Config.ThresholdFMR100)*30 + 50
}
func interpolate(raw, min, max, start, length float64) float64 {
	return (raw-min)/(max-min)*length + start
}
