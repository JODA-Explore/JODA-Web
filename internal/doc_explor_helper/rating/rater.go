package rating

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
)

type Rater interface {
	Score() float64
	Desc() template.HTML
}

func PrettyScore(r Rater) template.HTML {
	return desc.Rating(r.Score())
}

func cut(value, bottom, top float64) float64 {
	if top >= 0 && value > top {
		return top
	}
	if bottom >= 0 && value < bottom {
		return bottom
	}
	return value
}

func normalize(score float64) float64 {
	return cut(score, 0, 0.9999)
}
