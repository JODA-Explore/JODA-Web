package rating

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
)

type GlobalScale float64

func (s GlobalScale) ToScale() float64 {
	return float64(s)
}

func (s GlobalScale) Desc() template.HTML {
	return desc.Collapse(
		desc.Scale(s.Name(), s.ToScale()),
		"", "", "",
	)
}

func (s GlobalScale) Name() string {
	return "Global Scale"
}
