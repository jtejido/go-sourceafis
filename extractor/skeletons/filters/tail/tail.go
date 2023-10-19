package tail

import (
	"sourceafis/config"
	"sourceafis/extractor/logger"
	"sourceafis/extractor/skeletons/filters/dot"
	"sourceafis/extractor/skeletons/filters/knot"
	"sourceafis/features"
)

type SkeletonTailFilter struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *SkeletonTailFilter {
	return &SkeletonTailFilter{
		logger: logger,
	}
}

func (f *SkeletonTailFilter) Apply(skeleton *features.Skeleton) error {
	for _, minutia := range skeleton.Minutiae {
		if len(minutia.Ridges) == 1 && len(minutia.Ridges[0].End().Ridges) >= 3 {
			if minutia.Ridges[0].Points.Size() < config.Config.MinTailLength {
				minutia.Ridges[0].Detach()
			}
		}
	}
	if err := dot.Apply(skeleton); err != nil {
		return err
	}

	if err := knot.Apply(skeleton); err != nil {
		return err
	}

	f.logger.LogSkeleton("removed-tails", skeleton)
	return nil
}
