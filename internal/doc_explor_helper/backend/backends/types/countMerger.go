package types

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/feature"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

var _ feature.CountMerger[Info] = Backend{} // Verify Merger is implemented.

func (b Backend) TryMerge(x, y Info) bool {
	conds := append(x.conds, y.conds...)
	count := b.CountMust(conds.Query(b.Source()))
	if num.Diff(count, x.percent.Count()) > b.maxDiff {
		return false
	}
	return true
}

func (b Backend) CountMerge(infos []Info) (res [][]Info) {
	counts := feature.NewCounts(func(info Info) uint64 { return info.percent.Count() }, b.maxDiff)
	for _, info := range infos {
		counts.Insert(info)
	}
	return counts.Infos
}
