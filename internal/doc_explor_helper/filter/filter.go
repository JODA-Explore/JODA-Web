package filter

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
)

type Filter struct {
	MaxNumber          int
	MinRating, MaxDiff float64
}

var filterTmpl = template.Must(
	template.New("filter").
		Parse(`
<ul> 
  <li> MaxNumber: <b>{{.MaxNumber}}</b> </li> 
  <li> MinRating: <b>{{.MinRating}}</b> </li> 
</ul>
`),
)

func (f Filter) Desc() template.HTML {
	return desc.ExecuteTemplateWeb(*filterTmpl, f)
}
