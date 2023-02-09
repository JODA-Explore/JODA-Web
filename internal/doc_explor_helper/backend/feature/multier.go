package feature

import (
	"errors"
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/points"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
)

// Multier is a feature to process multiple infos into a single mixed info.
// it will not be used independtly but with a merger feature to group different infos.
type Multier[I result.Info] interface {
	Main[I]

	// MultiInfo generates extra info by the given grouped infos. Normally extra info contains some meta infos over the grouped infos.
	// since ExtraInfo is an interface, it can also be nil.
	MultiInfo([]I) result.ExtraInfo[I]
	MultiDesc([]I, result.ExtraInfo[I], int) (template.HTML, error)
	MultiQryMaker([]I, result.ExtraInfo[I], int) (template.HTML, error)

	// MultiNum return the limit number of info-candidates.
	MultiNum(filter.Filter) int
}

// multierRs generate results by the given Multier interface.
func multierRs[I result.Info](m Multier[I], topK topK[I], f filter.Filter) (rs result.Results) {
	switch be := m.(type) {
	case CountMerger[I]:
		return countRs(be, topK, f)
	}
	return
}

// multierTopK returns the top-K infos by its rating, k is calculated by MultiNum.
func multierTopK[I result.Info](m Multier[I], pts *points.Trie, f filter.Filter) (topK topK[I]) {
	topK = newTopK[I](m.MultiNum(f))
	addInfos := func(infos []I) {
		for _, info := range infos {
			topK.insert(info)
		}
	}
	walker := makeWalker(Main[I](m), addInfos)
	pts.Walk(walker)
	return
}

// mixedResult converts multi-infos to a single info, then generate the result by the single info,
// if the infos have only one member, then the returned result is exactly the same as the main interface without multier feature,
// if the infos have multiple members, a mixed info will be generated, then return the result of the mixed info.
func mixedResult[I result.Info](m Multier[I], infos []I) *result.Result {
	var info result.Info
	switch len(infos) {
	case 0:
		panic(errors.New("the number of infos can not be zero!"))
	case 1:
		info = infos[0]
	default:
		mixed := result.NewMixedInfo(infos, m.MultiInfo(infos))
		info = &mixed
	}
	r := result.New(info.Rating(), info)
	return &r
}
