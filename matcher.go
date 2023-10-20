package sourceafis

import (
	"context"

	"github.com/jtejido/sourceafis/matcher"
	"github.com/jtejido/sourceafis/templates"
)

type Matcher struct {
	probe   *matcher.Probe
	matcher *matcher.Matcher
}

func NewMatcher(logger matcher.MatcherLogger, probe *templates.SearchTemplate) (*Matcher, error) {
	hashBuilder := matcher.NewEdgeHashBuilder(logger.(matcher.HashTableLogger))
	hash, err := hashBuilder.Build(probe)
	if err != nil {
		return nil, err
	}
	return &Matcher{
		matcher: matcher.NewMatcher(logger),
		probe:   matcher.NewProbe(probe, hash),
	}, nil
}

func (m *Matcher) Match(ctx context.Context, candidate *templates.SearchTemplate) float64 {
	return m.matcher.Match(ctx, m.probe, candidate)
}
