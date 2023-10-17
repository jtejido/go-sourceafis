package cloud

import (
	"sourceafis/config"
	"sourceafis/features"
	"sourceafis/utils"
)

func Apply(minutiae []*features.FeatureMinutia) []*features.FeatureMinutia {
	radiusSq := utils.SquareInt(config.Config.MinutiaCloudRadius)

	var kept []*features.FeatureMinutia
	for _, minutia := range minutiae {
		count := 0
		for _, neighbor := range minutiae {
			if neighbor != minutia && neighbor.Position.Minus(minutia.Position).LengthSq() <= radiusSq {
				count++
			}
		}

		if config.Config.MaxCloudSize >= count-1 {
			kept = append(kept, minutia)
		}
	}

	minutiae = append(minutiae, kept...)
	return kept
}
