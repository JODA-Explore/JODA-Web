package feature

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
)

// SimpleMerger means that there are only one group in total,
// all infos are either in this group or in a group with single element(i.e. itself).
type SimpleMerger[I result.Info] interface {
	Multier[I]
	InGroup(I) bool
}

func simpleMergeRs[I result.Info](sm SimpleMerger[I], topK topK[I], f filter.Filter) (rs result.Results) {
	var group []result.Info
	for _, x := range topK.Vals() {
		if sm.InGroup(x) {
			group = append(group, x)
		} else {
			rs.TryInsert(x, x.Rating(), f)
		}
	}
	mixed := result.NewMixedInfo(group, nil)
	rs.TryInsert(&mixed, mixed.Rating(), f)
	return
}
