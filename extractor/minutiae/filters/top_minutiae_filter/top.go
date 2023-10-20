package top

import (
	"sort"
	"sourceafis/config"
	"sourceafis/features"
	"sourceafis/primitives"
)

func Apply(minutiae *primitives.GenericList[*features.FeatureMinutia]) *primitives.GenericList[*features.FeatureMinutia] {
	if minutiae.Len() <= config.Config.MaxMinutiae {
		return minutiae
	}
	sortedList := primitives.NewGenericList[*features.FeatureMinutia]()
	for e := minutiae.Front(); e != nil; e = e.Next() {
		minutia := e.Value.(*features.FeatureMinutia)

		// Create a slice to store neighbor distances
		neighborDistances := make([]int, 0, minutiae.Len())

		// Iterate over neighbors and calculate distances
		for e2 := minutiae.Front(); e2 != nil; e2 = e2.Next() {
			neighbor := e2.Value.(*features.FeatureMinutia)
			distanceSq := minutia.Position.Minus(neighbor.Position).LengthSq()
			neighborDistances = append(neighborDistances, distanceSq)
		}

		// Sort neighbor distances and skip the specified number of neighbors
		sort.Ints(neighborDistances)
		if len(neighborDistances) > config.Config.SortByNeighbor {
			neighborDistances = neighborDistances[config.Config.SortByNeighbor:]
		}

		// Find the first distance (or set it to Integer.MAX_VALUE if none found)
		firstDistance := neighborDistances[0]
		if len(neighborDistances) == 0 {
			firstDistance = int(^uint(0) >> 1) // Max int value
		}

		// Add the FeatureMinutia with the firstDistance as the key to the sorted list
		sortedList.PushBack(FeatureMinutiaWithDistance{minutia, firstDistance})
	}

	// Sort the list in reverse order by the firstDistance key
	sort.Sort(sort.Reverse(ByDistance{sortedList}))

	if sortedList.Len() > config.Config.MaxMinutiae {
		sortedList.Init()
	}
	top := primitives.NewGenericList[*features.FeatureMinutia]()
	for temp := sortedList.Front(); temp != nil; temp = temp.Next() {
		minutia := temp.Value.(*features.FeatureMinutia)
		top.PushBack(minutia)
	}

	return top
}

type FeatureMinutiaWithDistance struct {
	Minutia  *features.FeatureMinutia
	Distance int
}

type ByDistance struct {
	list *primitives.GenericList[*features.FeatureMinutia]
}

func (s ByDistance) Len() int {
	return s.list.Len()
}

func (s ByDistance) Swap(i, j int) {
	// Not needed for a singly linked list
}

func (s ByDistance) Less(i, j int) bool {
	minutia1 := s.list.Front().Value.(FeatureMinutiaWithDistance)
	minutia2 := s.list.Back().Value.(FeatureMinutiaWithDistance)
	return minutia1.Distance < minutia2.Distance
}
