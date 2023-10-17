package skeletons

import (
	"sourceafis/extractor/logger"
	"sourceafis/extractor/skeletons/filters"
	"sourceafis/extractor/skeletons/thinner"
	"sourceafis/extractor/skeletons/tracer"
	"sourceafis/features"
	"sourceafis/primitives"
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
	g.logger.Log(t.String()+"binarized-skeleton", binary)
	thinned := g.thinner.Thin(binary, t)

	skeleton, err := g.tracer.Trace(thinned, t)
	if err != nil {
		return nil, err
	}

	if err := g.filters.Apply(skeleton); err != nil {
		return nil, err
	}

	return skeleton, nil
}
