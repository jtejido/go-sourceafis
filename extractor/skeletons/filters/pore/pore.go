package pore

import (
	"github.com/jtejido/sourceafis/config"
	"github.com/jtejido/sourceafis/extractor/logger"
	"github.com/jtejido/sourceafis/extractor/skeletons/filters/knot"
	"github.com/jtejido/sourceafis/features"
)

type SkeletonPoreFilter struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *SkeletonPoreFilter {
	return &SkeletonPoreFilter{
		logger: logger,
	}
}

func (f *SkeletonPoreFilter) Apply(skeleton *features.Skeleton) error {

	for _, minutia := range skeleton.Minutiae {
		if len(minutia.Ridges) == 3 {
			for exit := 0; exit < 3; exit++ {
				exitRidge := minutia.Ridges[exit]
				arm1 := minutia.Ridges[(exit+1)%3]
				arm2 := minutia.Ridges[(exit+2)%3]
				if arm1.End() == arm2.End() && exitRidge.End() != arm1.End() && arm1.End() != minutia && exitRidge.End() != minutia {
					end := arm1.End()
					if len(end.Ridges) == 3 && arm1.Points.Size() <= config.Config.MaxPoreArm && arm2.Points.Size() <= config.Config.MaxPoreArm {
						arm1.Detach()
						arm2.Detach()
						merged := features.NewSkeletonRidge()
						merged.SetStart(minutia)
						merged.SetEnd(end)
						for _, point := range minutia.Position.LineTo(end.Position) {
							merged.Points.Add(point)
						}
					}
					break
				}
			}
		}
	}

	err := knot.Apply(skeleton)
	if err != nil {
		return err
	}

	return f.logger.LogSkeleton("removed-pores", skeleton)
}
