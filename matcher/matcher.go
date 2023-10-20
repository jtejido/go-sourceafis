package matcher

import (
	"context"

	"github.com/jtejido/sourceafis/templates"
)

type MatcherLogger interface {
	LogRootPairs(count int, roots []*MinutiaPair) error
	LogPairing(pairing *PairingGraph) error
	LogBestPairing(pairing *PairingGraph) error
	LogScore(score *ScoringData) error
	LogBestScore(score *ScoringData) error
	LogBestMatch(int) error
}

type Matcher struct {
	logger MatcherLogger
}

func NewMatcher(logger MatcherLogger) *Matcher {
	return &Matcher{
		logger: logger,
	}
}

func (m *Matcher) Match(ctx context.Context, probe *Probe, candidate *templates.SearchTemplate) (high float64) {
	thread := CurrentThread(ctx)
	defer ctx.Done()

	thread.Pairing.ReserveProbe(probe)
	thread.Pairing.ReserveCandidate(candidate)
	thread.Pairing.supportEnabled = true
	Enumerate(probe, candidate, thread.Roots)
	m.logger.LogRootPairs(thread.Roots.count, thread.Roots.pairs)

	best := -1
	for i := 0; i < thread.Roots.count; i++ {
		Crawl(probe.template.Edges, candidate.Edges, thread.Pairing, thread.Roots.pairs[i], &thread.Queue)
		m.logger.LogPairing(thread.Pairing)
		Compute(probe.template, candidate, thread.Pairing, thread.Score)
		m.logger.LogScore(thread.Score)
		partial := thread.Score.ShapedScore
		if best < 0 || partial > high {
			high = partial
			best = i
		}
		thread.Pairing.Clear()
	}

	if best >= 0 {
		thread.Pairing.supportEnabled = true
		Crawl(probe.template.Edges, candidate.Edges, thread.Pairing, thread.Roots.pairs[best], &thread.Queue)
		m.logger.LogBestPairing(thread.Pairing)
		Compute(probe.template, candidate, thread.Pairing, thread.Score)

		m.logger.LogBestScore(thread.Score)
		thread.Pairing.Clear()
	}

	thread.Roots.Discard()
	m.logger.LogBestMatch(best)
	return
}
