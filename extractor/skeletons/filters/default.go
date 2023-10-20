package filters

import (
	"github.com/jtejido/sourceafis/extractor/logger"
	"github.com/jtejido/sourceafis/extractor/skeletons/filters/dot"
	"github.com/jtejido/sourceafis/extractor/skeletons/filters/fragment"
	"github.com/jtejido/sourceafis/extractor/skeletons/filters/gap"
	"github.com/jtejido/sourceafis/extractor/skeletons/filters/pore"
	"github.com/jtejido/sourceafis/extractor/skeletons/filters/tail"
	"github.com/jtejido/sourceafis/features"
)

type SkeletonFilters struct {
	logger   logger.TransparencyLogger
	pore     *pore.SkeletonPoreFilter
	gap      *gap.SkeletonGapFilter
	tail     *tail.SkeletonTailFilter
	fragment *fragment.SkeletonFragmentFilter
}

func New(logger logger.TransparencyLogger) *SkeletonFilters {
	return &SkeletonFilters{
		logger:   logger,
		pore:     pore.New(logger),
		gap:      gap.New(logger),
		tail:     tail.New(logger),
		fragment: fragment.New(logger),
	}
}

func (f *SkeletonFilters) Apply(skeleton *features.Skeleton) error {
	if err := dot.Apply(skeleton); err != nil {
		return err
	}

	if err := f.logger.LogSkeleton("removed-dots", skeleton); err != nil {
		return err
	}
	if err := f.pore.Apply(skeleton); err != nil {
		return err
	}

	if err := f.gap.Apply(skeleton); err != nil {
		return err
	}

	if err := f.tail.Apply(skeleton); err != nil {
		return err
	}

	return f.fragment.Apply(skeleton)
}
