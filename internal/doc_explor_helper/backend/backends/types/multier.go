package types

import (
	"html/template"
	"strconv"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/feature"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/cond"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

var _ feature.Multier[Info] = Backend{} // Verify Multier is implemented.

func (b Backend) MultiNum(f filter.Filter) int {
	return b.CandidatesNumber(f)
}

func (b Backend) MultiInfo(infos []Info) result.ExtraInfo[Info] {
	return nil
}

func (b Backend) MultiDesc(infos []Info, _ result.ExtraInfo[Info], idx int) (template.HTML, error) {
	var avgPercent, details template.HTML
	var avgTotalCount uint64
	var conds cond.Conditions
	for i, x := range infos {
		avgTotalCount += x.percent.Total()
		conds = append(conds, x.conds...)
		details += desc.Collapse(
			template.HTML("Condition "+strconv.Itoa(i)+": "),
			"",
			x.desc(),
			"",
		)
	}
	newCount := b.CountMust(conds.Query(b.Source()))
	if newCount != 0 && avgTotalCount != 0 {
		avgTotalCount = num.Div(avgTotalCount, len(infos))
		avgPercent = desc.Percent(
			num.NewPercent(newCount, avgTotalCount),
			"Count",
			"Total",
		)
	}
	return conds.Desc() + avgPercent + desc.Collapse(
		"<b>Details:</b>",
		template.HTML("It was reduced by "+strconv.Itoa(len(infos))+" conditions, since they have similiar results."),
		details,
		"",
	), nil
}

func (b Backend) MultiQryMaker(infos []Info, _ result.ExtraInfo[Info], idx int) (template.HTML, error) {
	var conds cond.Conditions
	for _, info := range infos {
		conds = append(conds, info.conds...)
	}
	return b.condsToQryMaker(conds, idx), nil
}
