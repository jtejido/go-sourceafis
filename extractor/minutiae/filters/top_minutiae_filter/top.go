package top

import (
	"sort"
	"sourceafis/config"
	"sourceafis/features"
)

func Apply(minutiae []*features.FeatureMinutia) []*features.FeatureMinutia {
	if len(minutiae) <= config.Config.MaxMinutiae {
		return minutiae
	}

	sort.Slice(minutiae, func(i, j int) bool {
		return compareMinutiae(minutiae[i], minutiae) > compareMinutiae(minutiae[j], minutiae)
	})

	return minutiae[:config.Config.MaxMinutiae]
}

func compareMinutiae(minutia *features.FeatureMinutia, allMinutiae []*features.FeatureMinutia) int {
	neighborDistances := make([]int, 0)

	for _, neighbor := range allMinutiae {
		if neighbor != minutia {
			neighborDistances = append(neighborDistances, minutia.Position.Minus(neighbor.Position).LengthSq())
		}
	}

	sort.Ints(neighborDistances)
	return neighborDistances[config.Config.SortByNeighbor]
}
