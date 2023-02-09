package rating

import "html/template"

type ConstRating float64

func (cr ConstRating) Score() float64 {
	return float64(cr)
}

func (cr ConstRating) Desc() template.HTML {
	return ""
}

func (cr ConstRating) Update() {
}
