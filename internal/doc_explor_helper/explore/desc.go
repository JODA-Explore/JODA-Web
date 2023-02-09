package explore

import (
	"encoding/json"
	"html/template"
	"reflect"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
)

func DescRoot() template.HTML {
	return datasets[""].Tree.desc()
}

var datasetTreeTmpl = template.Must(
	template.New("datasetTree").
		Parse(`<br>
{{- define "datasetTreeRec"}}
<li>
    <span onmouseenter="addDatasetBtns(this, '{{.Dataset}}')" onmouseleave="deleteAllBtns(this)">{{.Dataset}} </span>
</li>
<ul>
    {{range .Children}}
    {{- template "datasetTreeRec" .}}
    {{end}}
</ul>
{{- end}}

<ul>
    {{range .Children}}
    {{- template "datasetTreeRec" .}}
    {{end}}
</ul>
`),
)

func (tree DatasetTree) desc() template.HTML {
	return desc.ExecuteTemplateWeb(*datasetTreeTmpl, tree)
}

var datasetInfoTmpl = template.Must(
	template.New("datasetInfo").
		Parse(`<br>
<details>
    <summary>
        <b>Query:</b>
    </summary>
    {{if eq .Query ""}}
    <p>unknown </p>
    {{else}}
        <div style="overflow-wrap: break-word; background-color:powderblue;margin-left: 1em;">
            {{.Query}}
        </div>
    {{end}}
</details>
<details>
    <summary>
        <b>History:</b>
    </summary>
    <div style="margin-left: 1em;">
        {{range $i, $opt := .Opts}}
        <details>
            <summary>
                <b>Option {{$i}}</b>
            </summary>
            <div class="grid_wrapper" onmouseenter="addOptBtns(this,{{$i}})" onmouseleave="deleteAllBtns(this)">
                <ul class="grid_opt">
                    <li> Backends
                        {{$opt.DescEnabledBackends}}
                    </li>
                    <li> Filter
                        {{$opt.Filter.Desc}}
                    </li>
                </ul>
            </div>
        </details>
        {{end}}
    </div>
</details>
<br>
`),
)

func (di DatasetInfo) Desc() template.HTML {
	return desc.ExecuteTemplateWeb(*datasetInfoTmpl, di)
}

func (opt SearchOpt) DescEnabledBackends() template.HTML {
	return desc.List(opt.EnabledBackendsName())
}

func DescDataset(dataset string) template.HTML {
	return template.HTML(
		`<br>Details For <b id="dataset-details-name">`+dataset+`</b>:`,
	) + datasets[dataset].Desc()
}

func (opt SearchOpt) ToJson() (res []byte, err error) {
	m := make(map[string]interface{})
	v := reflect.ValueOf(opt)
	for i := 0; i < v.NumField(); i++ {
		m[v.Type().Field(i).Name] = v.Field(i).Interface()
	}
	res, err = json.Marshal(m)
	return
}
