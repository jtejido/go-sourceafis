package templates

import (
	"sourceafis/features"
	"sourceafis/primitives"
)

type FeatureTemplate struct {
	Size     primitives.IntPoint
	Minutiae *primitives.GenericList[*features.FeatureMinutia]
}

func NewFeatureTemplate(size primitives.IntPoint, minutiae *primitives.GenericList[*features.FeatureMinutia]) *FeatureTemplate {
	return &FeatureTemplate{
		Size:     size,
		Minutiae: minutiae,
	}
}
