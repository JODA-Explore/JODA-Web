package feature

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/points"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/trie"
)

type Main[I result.Info] interface {
	ID() backends.ID

	// Filter will filter JSON points by the value of returned trie.Control.
	Filter() points.Walker

	// Infos generate a list of Infos according to the given path and points info,
	// the children of the path may be skipped depends on the value of returned trie.Control.
	Infos(path string, pi points.Info) ([]I, trie.Control)

	// Desc generate the description of the given info.
	Desc(I, int) (template.HTML, error)

	// QryMaker generate the query maker of the given info.
	QryMaker(I, int) (template.HTML, error)
}

// basicRs generate results by the given Main interface(without any other optional features).
func basicRs[I result.Info](m Main[I], pts *points.Trie, f filter.Filter) (rs result.Results) {
	addInfos := func(infos []I) {
		for _, info := range infos {
			rs.TryInsert(info, info.Rating(), f)
		}
	}
	walker := makeWalker(m, addInfos)
	pts.Walk(walker)
	return
}

// makeWalker helps to generate a Walker function,
// it will first filter json points, then generate infos, the argument addInfos decides how to process those infos.
func makeWalker[I result.Info](m Main[I], addInfos func([]I)) points.Walker {
	filter := m.Filter()
	return func(path string, v *points.Value, depth int) (ctl trie.Control) {
		if v.Info.CountNull == v.Info.CountTotal {
			return trie.Continue
		}
		ctl = filter(path, v, depth)
		if ctl != trie.None {
			return
		}
		infos, ctl := m.Infos(path, v.Info)
		addInfos(infos)
		return ctl
	}
}
