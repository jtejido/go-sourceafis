package dot

import (
	"sourceafis/features"
)

func Apply(skeleton *features.Skeleton) error {
	removed := make([]*features.SkeletonMinutia, 0)
	for _, minutia := range skeleton.Minutiae {
		if len(minutia.Ridges) == 0 {
			removed = append(removed, minutia)
		}
	}
	for _, minutia := range removed {
		skeleton.RemoveMinutia(minutia)
	}

	return nil
}
