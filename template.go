package sourceafis

import (
	"sourceafis/extractor"
	"sourceafis/extractor/logger"
	"sourceafis/primitives"
	"sourceafis/templates"
)

type Extractor interface {
	Extract(raw *primitives.Matrix, dpi float64) (*templates.FeatureTemplate, error)
}

type TemplateCreator struct {
	logger    logger.TransparencyLogger
	extractor Extractor
}

func NewTemplateCreator(logger logger.TransparencyLogger) *TemplateCreator {
	return &TemplateCreator{
		logger:    logger,
		extractor: extractor.New(logger),
	}
}

func (c *TemplateCreator) Template(img *Image) (*templates.SearchTemplate, error) {
	ft, err := c.extractor.Extract(img.matrix, img.dpi)
	if err != nil {
		return nil, err
	}
	if c.logger != nil {
		err := c.logger.Log("shuffled-minutiae", ft)
		if err != nil {
			return nil, err
		}
	}

	return templates.NewSearchTemplate(c.logger, ft), nil
}
