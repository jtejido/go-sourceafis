package skeletons

import (
	"github.com/jtejido/sourceafis/extractor/logger"
	"github.com/jtejido/sourceafis/extractor/skeletons/filters"
	"github.com/jtejido/sourceafis/extractor/skeletons/thinner"
	"github.com/jtejido/sourceafis/extractor/skeletons/tracer"
	"github.com/jtejido/sourceafis/features"
	"github.com/jtejido/sourceafis/primitives"
)

type SkeletonGraphs struct {
	logger  logger.TransparencyLogger
	thinner *thinner.BinaryThinning
	tracer  *tracer.SkeletonTracing
	filters *filters.SkeletonFilters
}

func New(logger logger.TransparencyLogger) *SkeletonGraphs {
	return &SkeletonGraphs{
		logger:  logger,
		thinner: thinner.New(logger),
		tracer:  tracer.New(logger),
		filters: filters.New(logger),
	}
}

func (g *SkeletonGraphs) Create(binary *primitives.BooleanMatrix, t features.SkeletonType) (*features.Skeleton, error) {
	if err := g.logger.Log(t.String()+"binarized-skeleton", binary); err != nil {
		return nil, err
	}

	thinned, err := g.thinner.Thin(binary, t)
	if err != nil {
		return nil, err
	}

	skeleton, err := g.tracer.Trace(thinned, t)
	if err != nil {
		return nil, err
	}

	if err := g.filters.Apply(skeleton); err != nil {
		return nil, err
	}

	return skeleton, nil
}
