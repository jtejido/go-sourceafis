package templates

import (
	"sourceafis/features"
	"sourceafis/primitives"
)

type FeatureTemplate struct {
	Size     primitives.IntPoint
	Minutiae []*features.FeatureMinutia
}

func NewFeatureTemplate(size primitives.IntPoint, minutiae []*features.FeatureMinutia) *FeatureTemplate {
	return &FeatureTemplate{
		Size:     size,
		Minutiae: minutiae,
	}
}
