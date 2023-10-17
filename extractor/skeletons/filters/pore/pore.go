package pore

import (
	"sourceafis/config"
	"sourceafis/extractor/logger"
	"sourceafis/extractor/skeletons/filters/knot"
	"sourceafis/features"
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
				if arm1.EndMinutia() == arm2.EndMinutia() && exitRidge.EndMinutia() != arm1.EndMinutia() && arm1.EndMinutia() != minutia && exitRidge.EndMinutia() != minutia {
					end := arm1.EndMinutia()
					if len(end.Ridges) == 3 && arm1.Points.Size() <= config.Config.MaxPoreArm && arm2.Points.Size() <= config.Config.MaxPoreArm {
						arm1.Detach()
						arm2.Detach()
						merged := features.NewSkeletonRidge()
						merged.Start(minutia)
						merged.End(end)
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
	f.logger.LogSkeleton("removed-pores", skeleton)
	return nil
}
