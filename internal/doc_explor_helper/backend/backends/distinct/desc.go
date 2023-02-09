package distinct

import (
	"html/template"
	"strconv"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/cmd"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
)

var sampleHeadTmpl = template.Must(template.New("sampleHead").Parse(`
<b>Top {{.}} in Sample</b>
`))

func (b Backend) SampleDesc(info Info, idx int) (template.HTML, error) {
	sampleHeader := desc.ExecuteTemplateWeb(*sampleHeadTmpl, b.topNumber)
	return desc.Collapse(sampleHeader, "", info.sampleTop.Desc(), ""), nil
}

var distinctValuesCollapseTmpl = template.Must(template.New("distinctValuesCollapse").Parse(`<br>
{{$doc_id := printf "distinctValues_%v" .id}}
<details id="details_{{$doc_id}}">
    <summary onclick="loadDistinctValues('{{.dataset}}','{{.path}}', {{.number}}, {{.id}})">
       <b>Distinct Values</b>
    </summary>
    <div id="{{$doc_id}}"><br>{{.content}}</div>
</details>`))

func distinctValuesCollapse(dataset, path string, number, id int, content template.HTML) template.HTML {
	data := map[string]interface{}{
		"dataset": dataset,
		"path":    path,
		"number":  number,
		"id":      id,
		"content": content,
	}
	return desc.ExecuteTemplateWeb(*distinctValuesCollapseTmpl, data)
}

func (b Backend) Desc(info Info, idx int) (template.HTML, error) {
	if info.topPercent.Total() == 0 {
		return "", nil
	}

	percent := desc.Percent(info.topPercent, "Top"+strconv.Itoa(b.topNumber), "Total")

	buttonbar := desc.ButtonBarMaker{
		QueryId:       idx,
		ButtonBarName: "distinctValues",
		Dataset:       b.Source(),
		Children: []struct {
			ButtonName, DocId string
			DivContent        template.HTML
		}{
			{"Top 10", "coverage-doughnets", desc.Canvas(idx, "coverage-doughnets")},
			{"Top 50", "coverage-bar", desc.Canvas(idx, "coverage-bar")},
		},
	}.Make()
	collapse := distinctValuesCollapse(b.Source(), info.JsonPoint(), 50, idx, buttonbar)

	return desc.JsonPoint(info.JsonPoint()) + percent + collapse, nil
}

func (b Backend) QryMaker(info Info, idx int) (template.HTML, error) {
	vals, err := b.AllDistinctValues(b.Source(), info.JsonPoint())
	if err != nil {
		return "", err
	}
	aggQ := cmd.Load(b.Source()) + cmd.Agg("", cmd.Group(`COUNT('')`, "count", cmd.Quote(info.JsonPoint())))
	selectQ := desc.SelectListQuery(
		vals,
		cmd.Load(b.Source()),
		"single-value"+strconv.Itoa(idx),
		cmd.Load(b.Source())+cmd.Choose(cmd.Quote(info.JsonPoint())+" == "+`"%s"`),
	)

	qryMaker := desc.QueryMaker{
		QueryId: idx,
		Dataset: b.Source(),
		Children: []struct {
			ButtonName, DocId string
			ContentTmpl       *template.Template
		}{
			{"Agg", "agg", desc.ConstQuery(aggQ)},
			{"Single Value", "single-value", selectQ},
		},
	}
	return qryMaker.Make(), nil
}
