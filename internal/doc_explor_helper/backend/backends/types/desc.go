package types

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/cond"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/slices"
)

func (i Info) desc() template.HTML {
	return i.conds.Desc() + desc.Percent(i.percent, "Count", "Total")
}

func (b Backend) Desc(info Info, idx int) (template.HTML, error) {
	return info.desc(), nil
}

func negateConds(cs cond.Conditions) (res cond.Conditions) {
	res = slices.Map(cs, func(x cond.Condition) cond.Condition {
		if x.CondType() == cond.Exists {
			// exists && is null
			return x
		} else {
			return x.Negate()
		}
	})
	return
}

func (b Backend) condsToQryMaker(cs cond.Conditions, i int) template.HTML {
	neg := negateConds(cs)
	return desc.QueryMaker{
		QueryId: i,
		Dataset: b.Source(),
		Children: []struct {
			ButtonName, DocId string
			ContentTmpl       *template.Template
		}{
			{"Standard", "standard", desc.ConstQuery(cs.Query(b.Source()))},
			{"Inversed", "inversed", desc.ConstQuery(neg.Query(b.Source()))},
		},
	}.Make()
}

func (b Backend) QryMaker(info Info, idx int) (template.HTML, error) {
	cs := info.conds
	return b.condsToQryMaker(cs, idx), nil
}
