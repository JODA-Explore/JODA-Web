package cond

import (
	"html/template"
	"strings"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/slices"
)

type Condition interface {
	JsonPoint() string
	string() string
	IsNeg() bool
	Negate() Condition
	CondType() Type
	Label() string
	Desc() template.HTML
}

func New(
	jsonPoint string,
	neg bool,
	condType Type,
	extra interface{},
) Condition {
	switch condType {
	case Exists:
		return newExistsCond(jsonPoint, neg)
	case Equal:
		right := extra.(string)
		return newEqualCond(jsonPoint, right, neg)
	default:
		return newTypeCond(jsonPoint, neg, condType)
	}
}

func Str(c Condition) string {
	if c.IsNeg() {
		if c.CondType() == Equal {
			return strings.Replace(c.string(), "==", "!=", 1)
		}
		return "!" + c.string()
	} else {
		return c.string()
	}
}

type Conditions []Condition

func (cs Conditions) Query(dataSet string) string {
	var sb strings.Builder
	sb.WriteString("LOAD ")
	sb.WriteString(dataSet)
	sb.WriteString(" CHOOSE ")
	for i, x := range cs {
		sb.WriteString(Str(x))
		if i < len(cs)-1 {
			sb.WriteString(" && ")
		}
	}
	return sb.String()
}

func (cs Conditions) JsonPoints() (res []string) {
	return slices.Map(cs, func(c Condition) string { return c.JsonPoint() })
}

func (cs Conditions) Types() (res []Type) {
	return slices.Map(cs, func(c Condition) Type { return c.CondType() })
}

func (cs Conditions) filter(pred func(Condition) bool) (in, out Conditions) {
	for _, x := range cs {
		if pred(x) {
			in = append(in, x)
		} else {
			out = append(out, x)
		}
	}
	return
}

func (cs Conditions) GroupByNeg() (pos, neg Conditions) {
	pred := func(c Condition) bool {
		return !c.IsNeg()
	}
	return cs.filter(pred)
}

var tmpl = template.Must(template.New("Conditions").Parse(`
<ul>
    {{range .}}
    <li> {{.Desc}} </li>
    {{end}}
</ul>
`))

func (cs Conditions) GroupByType() (res []Conditions) {
	res = make([]Conditions, len(condTypeSlice))
	for _, x := range cs {
		t := x.CondType()
		res[t] = append(res[t], x)
	}
	return
}

func (cs Conditions) Desc() (res template.HTML) {
	add := func(css Conditions) {
		if len(css) == 0 {
			return
		}
		res += template.HTML(css[0].Label())
		res += desc.ExecuteTemplateWeb(*tmpl, css)
	}
	for _, x := range cs.GroupByType() {
		if len(x) == 0 {
			continue
		}
		pos, neg := x.GroupByNeg()
		add(pos)
		add(neg)
	}
	return
}
