package objects

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/cmd"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
)

func (b Backend) SampleDesc(info Info, idx int) (template.HTML, error) {
	return desc.Collapse(
		"<b>Sample:</b>",
		"<br>Distinct Values in sample",
		desc.Percent(info.sample, "Distinct", "Total"),
		"",
	), nil
}

var estimatedObjTmpl = template.Must(template.New("estimatedObj").Parse(`<br>
{{.jsonPoint}}<br><br>
Count of Distinct Values (Estimated by sample):{{.percent}}
{{.sample}}
`))

var realObjTmpl = template.Must(template.New("realObj").Parse(`<br>
{{.jsonPoint}}<br><br>
Count of Distinct Values:{{.percent}}
`))

// CompleteDesc use the info of Query to generate describe
func (b Backend) Desc(info Info, idx int) (template.HTML, error) {
	data := map[string]interface{}{
		"jsonPoint": desc.JsonPoint(info.JsonPoint()),
		"percent":   desc.Percent(info.real, "Distinct", "Total"),
	}
	if b.UseSample(info) {
		return desc.ExecuteTemplateWeb(*estimatedObjTmpl, data), nil
	} else {
		return desc.ExecuteTemplateWeb(*realObjTmpl, data), nil
	}
}

func (b Backend) QryMaker(info Info, idx int) (template.HTML, error) {
	flattenedQ := cmd.Load(b.Source()) + cmd.As("", cmd.Fun("FLATTEN", cmd.Quote(info.JsonPoint())))
	qryMaker := desc.QueryMaker{
		QueryId: idx,
		Dataset: b.Source(),
		Children: []struct {
			ButtonName, DocId string
			ContentTmpl       *template.Template
		}{
			{"Flattened", "flattened", desc.ConstQuery(flattenedQ)},
		},
	}
	return qryMaker.Make(), nil
}
