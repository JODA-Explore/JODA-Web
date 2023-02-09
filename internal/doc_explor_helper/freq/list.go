package freq

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/topk"
)

type Freq struct {
	*topk.TopK[string, uint64]
}

func New(cap int) Freq {
	t := topk.New[string, uint64](cap)
	return Freq{&t}
}

var freqTmpl = template.Must(template.New("freq").Parse(`
<ul>
    {{range .}}
    <li> {{.Val}}   {{.N}} </li>
    {{end}}
</ul>`))

func (f Freq) Desc() template.HTML {
	return desc.ExecuteTemplateWeb(*freqTmpl, f.List())
}
