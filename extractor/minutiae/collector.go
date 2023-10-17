package minutiae

import (
	"sourceafis/features"
)

func Collect(ridges, valleys *features.Skeleton) ([]*features.FeatureMinutia, error) {
	endingsM, err := collect(ridges, features.ENDING)
	if err != nil {
		return nil, err
	}
	valleysM, err := collect(valleys, features.BIFURCATION)
	if err != nil {
		return nil, err
	}

	return append(valleysM, endingsM...), nil
}

func collect(skeleton *features.Skeleton, t features.MinutiaType) ([]*features.FeatureMinutia, error) {
	minutiae := make([]*features.FeatureMinutia, 0)
	for _, minutia := range skeleton.Minutiae {
		if len(minutia.Ridges) == 1 {
			dir, err := minutia.Ridges[0].Direction()
			if err != nil {
				return nil, err
			}
			minutiae = append(minutiae, features.NewFeatureMinutia(minutia.Position, dir, t))
		}
	}

	return minutiae, nil
}
