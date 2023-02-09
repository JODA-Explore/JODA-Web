package array

import (
	"html/template"
	"strconv"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/cmd"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/jsontype"
)

var sampleHeadTmpl = template.Must(template.New("sampleHead").Parse(`
<b>Top {{.}} in Sample</b>
`))

func (b Backend) SampleDesc(info Info, idx int) (template.HTML, error) {
	sampleHeader := desc.ExecuteTemplateWeb(*sampleHeadTmpl, b.topNumber)
	return desc.Collapse(sampleHeader, "", info.sampleTop.Desc(), ""), nil
}

var arrayCollapseTmpl = template.Must(template.New("arrayCollapse").Parse(`<br>
{{$doc_id := printf "arrayCollapse_%v" .id}}
<details id="details_{{$doc_id}}">
    <summary onclick="loadMemberFreq('{{.dataset}}','{{.path}}', {{.number}}, {{.id}})">
       <b>Member Frequency</b>
    </summary>
    <br>
    <b><span id=totalNumber_{{.id}}>{{.count}}</span></b> distinct members in total.
    <div id="{{$doc_id}}"><br>{{.content}}</div>
</details>`))

func arrayCollapse(dataset, path string, number, id int, content template.HTML) template.HTML {
	data := map[string]interface{}{
		"dataset": dataset,
		"path":    path,
		"number":  number,
		"id":      id,
		"content": content,
	}
	return desc.ExecuteTemplateWeb(*arrayCollapseTmpl, data)
}

var arrayTmpl = template.Must(template.New("arrayAnalysis").Parse(`
{{.jsonPoint}}
{{.percent}}
{{.real}}
`))

// CompleteDesc use the info of Query to generate describe
func (b Backend) Desc(info Info, idx int) (template.HTML, error) {
	buttonbar := desc.ButtonBarMaker{
		QueryId:       idx,
		ButtonBarName: "memberFreq",
		Dataset:       b.Source(),
		Children: []struct {
			ButtonName, DocId string
			DivContent        template.HTML
		}{
			{"Top 10", "coverage-doughnets", desc.Canvas(idx, "coverage-doughnets")},
			{"Top 50", "coverage-bar", desc.Canvas(idx, "coverage-bar")},
		},
	}.Make()
	collapse := arrayCollapse(b.Source(), info.JsonPoint(), 50, idx, buttonbar)
	topNumber := strconv.Itoa(b.topNumber)

	data := map[string]interface{}{
		"jsonPoint": desc.JsonPoint(info.JsonPoint()),
		"real":      collapse,
		"percent":   desc.Percent(info.topPercent, "Top"+topNumber, "Total"),
	}
	return desc.ExecuteTemplateWeb(*arrayTmpl, data), nil
}

func (b Backend) QryMaker(info Info, idx int) (template.HTML, error) {
	vals, err := b.AllDistinctMembers(b.Source(), info.JsonPoint())
	if err != nil {
		return "", nil
	}
	overviewQ := cmd.Load(b.Source()) +
		cmd.As("", cmd.Fun("FLATTEN", cmd.Quote(info.JsonPoint()))) +
		cmd.Agg("", cmd.Group(`COUNT('')`, "count", cmd.Quote("")))
	var member string
	if info.arrayType == jsontype.String {
		member = `"%s"`
	} else {
		member = `%s`
	}
	selectQ := desc.SelectListQuery(
		vals,
		cmd.Load(b.Source()),
		"single-value"+strconv.Itoa(idx),
		cmd.Load(b.Source())+cmd.Choose(cmd.Fun("IN", member, cmd.Quote(info.JsonPoint()))),
	)

	qryMaker := desc.QueryMaker{
		QueryId: idx,
		Dataset: b.Source(),
		Children: []struct {
			ButtonName, DocId string
			ContentTmpl       *template.Template
		}{
			{"Overview", "overview", desc.ConstQuery(overviewQ)},
			{"Single Member", "single-member", selectQ},
		},
	}
	return qryMaker.Make(), nil
}
