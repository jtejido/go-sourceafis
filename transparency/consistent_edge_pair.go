package transparency

import "github.com/jtejido/sourceafis/matcher"

type ConsistentEdgePair struct {
	ProbeFrom, ProbeTo, CandidateFrom, CandidateTo int
}

func NewConsistentEdgePair(pair *matcher.MinutiaPair) *ConsistentEdgePair {
	return &ConsistentEdgePair{
		ProbeFrom:     pair.ProbeRef,
		ProbeTo:       pair.Probe,
		CandidateFrom: pair.CandidateRef,
		CandidateTo:   pair.Candidate,
	}
}
