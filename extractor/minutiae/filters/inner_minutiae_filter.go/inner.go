package inner

import (
	"sourceafis/config"
	"sourceafis/features"
	"sourceafis/primitives"
)

func Apply(minutiae []*features.FeatureMinutia, mask *primitives.BooleanMatrix) []*features.FeatureMinutia {
	filteredMinutiae := []*features.FeatureMinutia{}

	for _, minutia := range minutiae {
		arrow := primitives.FloatAngle(minutia.Direction).ToVector().Multiply(-config.Config.MaskDisplacement).Round()

		if mask.GetPointWithFallback(minutia.Position.Plus(arrow), false) {
			filteredMinutiae = append(filteredMinutiae, minutia)
		}
	}

	return filteredMinutiae
}
