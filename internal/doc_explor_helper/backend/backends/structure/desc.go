package structure

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/cond"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/slices"
)

var tmpl = template.Must(template.New("structureInfoDesc").Parse(`
<span>Number Of Items: <b>{{.count}}</b> {{.percent}}% </span>`))

func (i Info) desc() template.HTML {
	data := map[string]interface{}{
		"count":   i.percent.Count(),
		"percent": desc.Float(i.percent.Ratio()*100, 4),
	}
	return i.conds.Desc() + desc.ExecuteTemplateWeb(*tmpl, data)
}

func (b Backend) queryMaker(conds cond.Conditions, i int) (template.HTML, error) {
	var negConds cond.Conditions = slices.Map(conds, func(c cond.Condition) cond.Condition { return c.Negate() })
	return desc.QueryMaker{
		QueryId: i,
		Dataset: b.Source(),
		Children: []struct {
			ButtonName, DocId string
			ContentTmpl       *template.Template
		}{
			{"Standard", "standard", desc.ConstQuery(conds.Query(b.Source()))},
			{"Inversed", "inversed", desc.ConstQuery(negConds.Query(b.Source()))},
		},
	}.Make(), nil
}

func (b Backend) Desc(info Info, idx int) (template.HTML, error) {
	return info.desc(), nil
}

func (b Backend) QryMaker(info Info, i int) (template.HTML, error) {
	return b.queryMaker(info.conds, i)
}
