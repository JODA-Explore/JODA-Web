package structure

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/feature"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
)

var _ feature.Multier[Info] = Backend{} // Verify Reducer is implemented.

func (b Backend) MultiNum(f filter.Filter) int {
	return b.CandidatesNumber(f)
}

func (b Backend) MultiInfo(infos []Info) result.ExtraInfo[Info] {
	extra := Info{}
	for _, x := range infos {
		extra.conds = append(extra.conds, x.conds...)
	}
	extra.percent = b.newPercent(b.CountMust(extra.conds.Query(b.Source())))
	return extra
}

func (b Backend) MultiDesc(_ []Info, extra result.ExtraInfo[Info], _ int) (template.HTML, error) {
	info := extra.Info()
	return info.desc(), nil
}

func (b Backend) MultiQryMaker(_ []Info, extra result.ExtraInfo[Info], i int) (template.HTML, error) {
	return b.queryMaker(extra.Info().conds, i)
}
