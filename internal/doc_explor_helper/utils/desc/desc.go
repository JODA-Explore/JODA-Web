package desc

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

var (
	singleQuoteReg = regexp.MustCompile(`(?P<singleQuote>'.*?')`)
	loadReg        = regexp.MustCompile(`(?P<load>LOAD [a-zA-Z]+ )`)
	exclamationReg = regexp.MustCompile(`\!`)
)

func replace(reg *regexp.Regexp, s, replace string) string {
	reg.MatchString(s)
	return reg.ReplaceAllString(s, replace)
}

func JsonPoint(jp string) template.HTML {
	return template.HTML(replace(singleQuoteReg, jp, `<font color="green">${singleQuote}</font>`))
}

func Rating(r float64) template.HTML {
	s := fmt.Sprintf("%.2f", r*100)
	res := `<span style="font-size:25px;">` + s[:2] + `</span><span style="font-size:15px;">` + s[2:] + `</span>`
	return template.HTML(res)
}

var inlineRatingTmpl = template.Must(template.New("InlineRating").Parse(`{{.name}}: <b>{{.score}}</b>`))

func InlineRating(name string, r float64) template.HTML {
	data := map[string]interface{}{
		"name":  name,
		"score": Float(100*r, 2),
	}
	return ExecuteTemplateWeb(*inlineRatingTmpl, data)
}

func Scale(name string, scale float64) template.HTML {
	var prefix string
	if scale >= 0 {
		prefix = `<span style="font-size:130%;">&#43;</span>` // plus sign
	}
	return template.HTML(name+" :  "+prefix) + Float(scale, 2) + template.HTML("%")
}

func RemoveLoad(jp string) string {
	return replace(loadReg, jp, ``)
}

func ExecuteTemplate(tmpl template.Template, data interface{}) (string, error) {
	var tmpRes bytes.Buffer
	err := tmpl.Execute(&tmpRes, data)
	if err != nil {
		return "", err
	}
	return tmpRes.String(), nil
}

func ExecuteTemplateWeb(tmpl template.Template, data interface{}) template.HTML {
	str, err := ExecuteTemplate(tmpl, data)
	if err != nil {
		panic(err)
	}
	return template.HTML(str)
}

var percentTmpl = template.Must(template.New("percent").Parse(`<div>
    <br>
    <div class="frac">
        <span>{{.countName}}</span>
        <span class="symbol">/</span>
        <span class="bottom">{{.totalName}}</span>
    </div>
    <span style="font-size:150%;"><b>:</b></span>
    <span style="font-size:150%;">
        <sup>{{.count}}</sup> &frasl; <sub> {{.total}}</sub>
    </span>
    <i>({{.percent}})</i>
</div>`))

func Percent(
	p num.Percent,
	countName, totalName string,
) template.HTML {
	data := make(map[string]interface{})
	data["count"] = p.Count()
	data["countName"] = countName
	data["total"] = p.Total()
	data["totalName"] = totalName
	data["percent"] = fmt.Sprintf("%.2f%%", p.Ratio()*100)

	return ExecuteTemplateWeb(*percentTmpl, data)
}

var collapseTmpl = template.Must(template.New("collapse").Parse(`<br>
<details>
    <summary>
        {{.head}}
    </summary>
    {{.desc}} 
    <div style="background-color:{{.blockcolor}};margin-left: 1em;">{{.block}} </div>
</details>`))

func Collapse(head, desc, block template.HTML, blockcolor string) template.HTML {
	data := make(map[string]interface{})
	data["head"] = head
	data["desc"] = desc
	data["blockcolor"] = blockcolor
	data["block"] = block
	return ExecuteTemplateWeb(*collapseTmpl, data)
}

var lazyCollapseTmpl = template.Must(template.New("lazyCollapse").Parse(`<br>
{{$doc_id := printf "%v_%v" .docName .id}}
<details id="details_{{$doc_id}}">
    <summary  onclick="fetchApi('{{.endpoint}}',{{.id}},'{{.dataset}}','{{$doc_id}}')">
        {{.head}}
    </summary>
    <div id="{{$doc_id}}" style="margin-left: 1em;"></div>
</details>`))

func LazyCollapse(
	head template.HTML,
	docName, endpoint, dataset string,
	id int,
) template.HTML {
	data := make(map[string]interface{})
	data["head"] = head
	data["docName"] = docName
	data["endpoint"] = endpoint
	data["dataset"] = dataset
	data["id"] = id
	return ExecuteTemplateWeb(*lazyCollapseTmpl, data)
}

func Float(flt float64, prec int) template.HTML {
	res := strconv.FormatFloat(flt, 'f', prec, 64)
	return template.HTML(strings.TrimRight(strings.TrimRight(res, "0"), "."))
}

var listTmpl = template.Must(template.New("list").Parse(`
<ul>
    {{range .}}
    <li> {{.}} </li>
    {{end}}
</ul>`))

func List[T any](l []T) template.HTML {
	return ExecuteTemplateWeb(*listTmpl, l)
}

var buttonTmpl = template.Must(template.New("Button").Parse(`
<div id="{{.elemName}}_{{.idx}}" style="margin-left: 1em;">
<button class="w3-hover-green w3-ripple" id="{{.buttonName}}_{{.idx}}"
        onclick="fetchApi('{{.endpoint}}',{{.idx}},'{{.dataset}}','{{.elemName}}_{{.idx}}')">
    <b>{{.buttonShow}}</b>
</button>
</div>
`))

func Button(
	buttonName, buttonShow, endpoint, dataset, elemName string,
	idx int,
) template.HTML {
	data := make(map[string]interface{})
	data["buttonName"] = buttonName
	data["buttonShow"] = buttonShow
	data["endpoint"] = endpoint
	data["dataset"] = dataset
	data["elemName"] = elemName
	data["idx"] = idx
	return ExecuteTemplateWeb(*buttonTmpl, data)
}

// (event,'{{$buttonBarId}}','{{$buttonBarName}}_{{$queryId}}_{{.DocId}}', {{.Endpoint}}, {{$dataset}},{{$queryId}})
var buttonBarTmpl = template.Must(
	template.New("ButtonBar").Parse(`
{{$queryId := .QueryId}}
{{$dataset := .Dataset}}
{{$buttonBarName := .ButtonBarName}}
{{$buttonBarId := printf "%v_%v" $buttonBarName $queryId}}
<div class="w3-container" id="{{$buttonBarId}}">
    <div class="w3-bar w3-black main-buttonbar">
        {{- range .Children}}
        {{$id := printf "%v_%v_%v" $buttonBarName $queryId .DocId }}
        <button id='{{$id}}_btn' class="w3-bar-item w3-btn" onclick="handleButtonBar('{{$id}}_btn','{{$buttonBarId}}',{{$id}})">
            <b>&nbsp;{{.ButtonName}}&nbsp;</b>
        </button>
        {{end}}
    </div>
    <div class="w3-container w3-border">
        {{- range .Children}}
        <div id="{{$buttonBarName}}_{{$queryId}}_{{.DocId}}" class="w3-container w3-animate-opacity content" style="display:none">{{.DivContent}}</div>
        {{- end}}
    </div>
</div>
`),
)

type ButtonBarMaker struct {
	QueryId                int
	ButtonBarName, Dataset string
	Children               []struct {
		ButtonName, DocId string
		DivContent        template.HTML
	}
}

func (bm ButtonBarMaker) Make() template.HTML {
	return ExecuteTemplateWeb(*buttonBarTmpl, bm)
}

var canvasTmpl = template.Must(
	template.New("canvas").
		Parse(`
<div class="chart-container">
    <canvas class="{{.canvasType}}" id="{{.canvasType}}_{{.queryId}}" aria-label="not support" role="img">
        <p> Your browser does not support the canvas element. </p>
    </canvas>
</div>
`),
)

func Canvas(queryId int, canvasType string) template.HTML {
	data := map[string]interface{}{
		"queryId":    queryId,
		"canvasType": canvasType,
	}
	return ExecuteTemplateWeb(*canvasTmpl, data)
}

var commonButtonsTmpl = template.Must(
	template.New("commonButtons").
		Parse(`<br>
 <div class="w3-bar">
     <button class="w3-btn w3-green w3-ripple" onclick="handleCopyQuery('{{.queryContentId}}')">
         <b>Copy</b>
     </button>
     <button class="w3-btn w3-green w3-ripple" onclick="handleExecQuery('{{.queryContentId}}');">
         <b>Execute</b>
     </button>
     <button class="w3-btn w3-green w3-ripple" onclick="handleExploreQuery('{{.queryContentId}}','{{.dataset}}',{{.queryId}})">
         <b>Explore</b>
     </button>
 </div><br>
`),
)

type QueryMaker struct {
	QueryId  int
	Dataset  string
	Children []struct {
		ButtonName, DocId string
		ContentTmpl       *template.Template
	}
}

func (qm QueryMaker) Make() template.HTML {
	buttonbarName := "QueryGenerator"
	newQueryContentId := func(queryId int, docId string) string {
		return "queryContent_" + buttonbarName + "_" + strconv.Itoa(queryId) + "_" + docId
	}
	children := make([]struct {
		ButtonName, DocId string
		DivContent        template.HTML
	}, len(qm.Children))
	for i, child := range qm.Children {
		data := map[string]interface{}{
			"queryContentId": newQueryContentId(qm.QueryId, child.DocId),
			"dataset":        qm.Dataset,
			"queryId":        qm.QueryId,
		}
		contentRes := ExecuteTemplateWeb(*child.ContentTmpl, data)
		commonButtonsRes := ExecuteTemplateWeb(*commonButtonsTmpl, data)
		children[i].ButtonName = child.ButtonName
		children[i].DocId = child.DocId
		children[i].DivContent = contentRes + commonButtonsRes
	}
	return ButtonBarMaker{qm.QueryId, buttonbarName, qm.Dataset, children}.Make()
}

var constQueryGeneratorTmpl = template.Must(
	template.New("constQueryGenerator").Delims("<<", ">>").
		Parse(`<br><p>Query:</p><div id="{{.queryContentId}}" style="background-color:powderblue;"><<.query>></div>`),
)

func constQueryStr(query string) string {
	data := make(map[string]interface{})
	data["query"] = query
	str, err := ExecuteTemplate(*constQueryGeneratorTmpl, data)
	if err != nil {
		panic(err)
	}
	return str
}

func ConstQuery(query string) *template.Template {
	constQuery, err := template.New("constQuery").Parse(constQueryStr(query))
	if err != nil {
		log.Println("ConstQuery error: ", err)
	}
	return constQuery
}

var selectListQueryGeneratorTmpl = template.Must(
	template.New("selectListQueryGenerator").Delims("<<", ">>").
		Parse(`<br>
<<$selectListId := printf "selectList_%v" .selectListId>>
Select one of the distinct values: <select id="<<$selectListId>>" onchange="selectQuery(<<.queryTmpl>>,{{.queryContentId}},<<$selectListId>>,<<.initValue>>)" style="width:100%;max-width:90%;">
    <option>--choose value--</option>
    <<range .values>>
    <option> <<.>> </option>
    <<end>>
</select>
`),
)

func SelectListQuery(
	values []string,
	initValue, selectListId string,
	queryTmpl string,
) *template.Template {
	data := make(map[string]interface{})
	data["values"] = values
	data["selectListId"] = selectListId
	data["queryTmpl"] = queryTmpl
	data["initValue"] = initValue
	var tmpRes bytes.Buffer
	err1 := selectListQueryGeneratorTmpl.Execute(&tmpRes, data)
	if err1 != nil {
		panic(err1)
	}
	selectListQuery, err2 := template.New("selectListQuery").Parse(tmpRes.String() + constQueryStr(initValue))
	if err2 != nil {
		log.Println("selectList error: ", err2)
	}
	return selectListQuery
}

var jsonPointTmpl = template.Must(
	template.New("jsonPoint").
		Parse(`<p>Json Point: <b>{{.jsonPoint}}</b></p>`),
)

func JsonPointDesc(jsonPoint string) template.HTML {
	data := make(map[string]interface{})
	data["jsonPoint"] = jsonPoint
	return ExecuteTemplateWeb(*jsonPointTmpl, data)
}
