package dot

import (
	"sourceafis/features"
)

func Apply(skeleton *features.Skeleton) error {

	for i := 0; i < len(skeleton.Minutiae); i++ {
		minutia := skeleton.Minutiae[i]
		if len(minutia.Ridges) == 0 {
			skeleton.RemoveMinutia(minutia)
		}
	}

	return nil
}
