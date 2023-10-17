package sourceafis

import (
	"sourceafis/matcher"
	"sourceafis/templates"
)

type Matcher struct {
	probe   *matcher.Probe
	matcher *matcher.Matcher
}

func NewMatcher(logger matcher.MatcherLogger, probe *templates.SearchTemplate) *Matcher {
	hashBuilder := matcher.NewEdgeHashBuilder(logger.(matcher.HashTableLogger))
	return &Matcher{
		matcher: matcher.NewMatcher(logger),
		probe:   matcher.NewProbe(probe, hashBuilder.Build(probe)),
	}
}

func (m *Matcher) Match(candidate *templates.SearchTemplate) float64 {
	return m.matcher.Match(m.probe, candidate)
}
