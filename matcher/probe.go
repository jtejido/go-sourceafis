package matcher

import (
	"github.com/jtejido/sourceafis/features"
	"github.com/jtejido/sourceafis/templates"
)

type Probe struct {
	template *templates.SearchTemplate
	hash     map[int][]*features.IndexedEdge
}

func NewProbe(template *templates.SearchTemplate, hash map[int][]*features.IndexedEdge) *Probe {
	return &Probe{
		template: template,
		hash:     hash,
	}
}
