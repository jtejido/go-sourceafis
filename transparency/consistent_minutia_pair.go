package transparency

import "sourceafis/matcher"

type ConsistentMinutiaPair struct {
	Probe, Candidate int
}

func Roots(count int, roots []*matcher.MinutiaPair) []ConsistentMinutiaPair {
	var consistentPairs []ConsistentMinutiaPair
	for i := 0; i < count && i < len(roots); i++ {
		p := roots[i]
		consistentPairs = append(consistentPairs, ConsistentMinutiaPair{
			Probe:     p.Probe,
			Candidate: p.Candidate,
		})
	}

	return consistentPairs
}
