package rating

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
)

type Factor interface {
	ToScale() float64
	Desc() template.HTML
	Name() string
}

type Factors []Factor

func (fs Factors) sumScale() (res float64) {
	for _, f := range fs {
		if f != nil {
			res += f.ToScale()
		}
	}
	return
}

func (fs Factors) avgScale() (res float64) {
	return fs.sumScale() / float64(len(fs))
}

func (fs Factors) maxScale() (max float64, maxIndex int) {
	max = -101
	maxIndex = -1
	for i, f := range fs {
		if f == nil {
			continue
		}
		value := f.ToScale()
		if value > max {
			max = value
			maxIndex = i
		}
	}
	return
}

func (fs Factors) FindFirst(name string) (index int, found bool) {
	for i, f := range fs {
		if f != nil && f.Name() == name {
			return i, true
		}
	}
	return
}

var factorsTmpl = template.Must(template.New("factors").Parse(`
    {{range .}}
    <div>
        {{.Desc}}
    </div>
    {{end}}
`))

func (fs Factors) Desc() template.HTML {
	if fs == nil {
		return ""
	}
	return desc.Collapse(desc.Scale("Scale", fs.sumScale()), "", desc.ExecuteTemplateWeb(*factorsTmpl, fs), "")
}
