package feature

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/topk"
)

type topK[I result.Info] struct {
	*topk.TopK[I, float64]
}

func (tp topK[I]) insert(i I) {
	tp.Insert(i, i.Rating().Score())
}

// realTopK returns the top-k by its rating score,
// the given tp are top-k by its estimated rating score.
func (tp topK[I]) realTopK() (real topK[I]) {
	list := tp.List()
	real = newTopK[I](len(list))
	for _, x := range tp.List() {
		real.insert(x.Val)
	}
	return
}

func (tp topK[I]) results(f filter.Filter) (rs result.Results) {
	for _, x := range tp.List() {
		rs.TryInsert(x.Val, x.Val.Rating(), f)
	}
	return
}

func newTopK[I result.Info](n int) topK[I] {
	t := topk.New[I, float64](n)
	return topK[I]{&t}
}
