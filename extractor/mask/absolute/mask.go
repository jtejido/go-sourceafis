package absolute

import (
	"sourceafis/config"
	"sourceafis/extractor/logger"
	"sourceafis/primitives"
)

type AbsoluteContrastMask struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *AbsoluteContrastMask {
	return &AbsoluteContrastMask{
		logger: logger,
	}
}

func (m *AbsoluteContrastMask) Compute(contrast *primitives.Matrix) *primitives.BooleanMatrix {
	result := primitives.NewBooleanMatrixFromPoint(contrast.Size())
	it := contrast.Size().Iterator()
	for it.HasNext() {
		block := it.Next()
		if contrast.GetPoint(block) < config.Config.MinAbsoluteContrast {
			result.SetPoint(block, true)
		}
	}
	m.logger.Log("absolute-contrast-mask", result)
	return result
}
