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
	if err := e.logger.Log("decoded-image", raw); err != nil {
		return nil, err
	}
	raw = resizer.Resize(raw, dpi)
	if err := e.logger.Log("scaled-image", raw); err != nil {
		return nil, err
	}
	blocks := primitives.NewBlockMap(raw.Width, raw.Height, config.Config.BlockSize)
	if err := e.logger.Log("blocks", blocks); err != nil {
		return nil, err
	}
	histogram, err := e.localHistogram.Create(blocks, raw)
	if err != nil {
		return nil, err
	}

	smoothHistogram, err := e.localHistogram.Smooth(blocks, histogram)
	if err != nil {
		return nil, err
	}
	mask, err := e.segmentationMask.Compute(blocks, histogram)
	if err != nil {
		return nil, err
	}

	equalized, err := e.equalizer.Equalize(blocks, raw, smoothHistogram, mask)
	if err != nil {
		return nil, err
	}
	orientation, err := e.orientations.Compute(equalized, mask, blocks)
	if err != nil {
		return nil, err
	}
	smoothed, err := e.smoothing.Parallel(equalized, orientation, mask, blocks)
	if err != nil {
		return nil, err
	}
	orthogonal, err := e.smoothing.Orthogonal(smoothed, orientation, mask, blocks)
	if err != nil {
		return nil, err
	}
	binary, err := e.binarizer.Binarize(smoothed, orthogonal, mask, blocks)
	if err != nil {
		return nil, err
	}

	pixelMask, err := e.segmentationMask.Pixelwise(mask, blocks)
	if err != nil {
		return nil, err
	}
	if err := e.binarizer.Cleanup(binary, pixelMask); err != nil {
		return nil, err
	}

	if err := e.logger.Log("pixel-mask", pixelMask); err != nil {
		return nil, err
	}
	inverted := e.binarizer.Invert(binary, pixelMask)

	innerMask, err := e.segmentationMask.Inner(pixelMask)
	if err != nil {
		return nil, err
	}
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

	if err := e.logger.Log("skeleton-minutiae", template); err != nil {
		return nil, err
	}
	inner.Apply(template.Minutiae, innerMask)

	if err := e.logger.Log("inner-minutiae", template); err != nil {
		return nil, err
	}
	cloud.Apply(template.Minutiae)
	if err := e.logger.Log("removed-minutia-clouds", template); err != nil {
		return nil, err
	}

	template = templates.NewFeatureTemplate(template.Size, top.Apply(template.Minutiae))
	if err := e.logger.Log("top-minutia", template); err != nil {
		return nil, err
	}
	return template, nil
}
