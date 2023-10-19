package extractor

import (
	"sourceafis/config"
	"sourceafis/extractor/binarizer"
	"sourceafis/extractor/equalizer"
	localHistogram "sourceafis/extractor/histogram"
	"sourceafis/extractor/logger"
	"sourceafis/extractor/minutiae"
	cloud "sourceafis/extractor/minutiae/filters/cloud_minutia_filter"
	inner "sourceafis/extractor/minutiae/filters/inner_minutiae_filter.go"
	top "sourceafis/extractor/minutiae/filters/top_minutiae_filter"
	"sourceafis/extractor/orientation"
	"sourceafis/extractor/resizer"
	"sourceafis/extractor/segmentation"
	"sourceafis/extractor/skeletons"
	"sourceafis/extractor/smoothing"
	"sourceafis/features"
	"sourceafis/primitives"
	"sourceafis/templates"
)

type Extractor struct {
	logger           logger.TransparencyLogger
	localHistogram   *localHistogram.LocalHistograms
	segmentationMask *segmentation.SegmentationMask
	equalizer        *equalizer.ImageEqualization
	orientations     *orientation.BlockOrientations
	smoothing        *smoothing.OrientedSmoothing
	binarizer        *binarizer.BinarizedImage
	skeletons        *skeletons.SkeletonGraphs
}

func New(logger logger.TransparencyLogger) *Extractor {
	return &Extractor{
		logger:           logger,
		localHistogram:   localHistogram.New(logger),
		segmentationMask: segmentation.New(logger),
		equalizer:        equalizer.New(logger),
		orientations:     orientation.New(logger),
		smoothing:        smoothing.New(logger),
		binarizer:        binarizer.New(logger),
		skeletons:        skeletons.New(logger),
	}
}

func (e *Extractor) Extract(raw *primitives.Matrix, dpi float64) (*templates.FeatureTemplate, error) {
	e.logger.Log("decoded-image", raw)
	raw = resizer.Resize(raw, dpi)
	e.logger.Log("scaled-image", raw)
	blocks := primitives.NewBlockMap(raw.Width, raw.Height, config.Config.BlockSize)
	e.logger.Log("blocks", blocks)
	histogram := e.localHistogram.Create(blocks, raw)

	smoothHistogram := e.localHistogram.Smooth(blocks, histogram)

	mask, err := e.segmentationMask.Compute(blocks, histogram)
	if err != nil {
		return nil, err
	}

	equalized := e.equalizer.Equalize(blocks, raw, smoothHistogram, mask)

	orientation := e.orientations.Compute(equalized, mask, blocks)

	smoothed := e.smoothing.Parallel(equalized, orientation, mask, blocks)

	orthogonal := e.smoothing.Orthogonal(smoothed, orientation, mask, blocks)

	binary := e.binarizer.Binarize(smoothed, orthogonal, mask, blocks)

	pixelMask := e.segmentationMask.Pixelwise(mask, blocks)

	e.binarizer.Cleanup(binary, pixelMask)

	e.logger.Log("pixel-mask", pixelMask)
	inverted := e.binarizer.Invert(binary, pixelMask)

	innerMask := e.segmentationMask.Inner(pixelMask)

	ridges, err := e.skeletons.Create(binary, features.RIDGES)
	if err != nil {
		return nil, err
	}

	valleys, err := e.skeletons.Create(inverted, features.VALLEYS)
	if err != nil {
		return nil, err
	}

	minutiae, err := minutiae.Collect(ridges, valleys)
	if err != nil {
		return nil, err
	}

	var template = templates.NewFeatureTemplate(raw.Size(), minutiae)

	e.logger.Log("skeleton-minutiae", template)
	inner.Apply(template.Minutiae, innerMask)

	e.logger.Log("inner-minutiae", template)
	cloud.Apply(template.Minutiae)
	e.logger.Log("removed-minutia-clouds", template)

	template = templates.NewFeatureTemplate(template.Size, top.Apply(template.Minutiae))
	e.logger.Log("top-minutia", template)
	return template, nil
}
