package result

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/rating"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

type Result struct {
	Info
	Rater    rating.Rater
	Desc     template.HTML
	QryMaker template.HTML
}

type Results []*Result

func New(r rating.Rater, i Info) Result {
	return Result{Rater: r, Info: i}
}

func (r Result) PrettyScore() template.HTML {
	return rating.PrettyScore(r.Rater)
}

// Extract put all queries from rs that pred returns true into in, and anothers into out.
func (rs *Results) Extract(pred func(*Result) bool) (in, out Results) {
	for _, x := range *rs {
		if pred(x) {
			// always sorted
			in = append(in, x)
		} else {
			out = append(out, x)
		}
	}
	return
}

func (rs *Results) ExtractByRating(rating float64) (higher, lower Results) {
	pred := func(r *Result) bool {
		return r.Rater.Score() >= rating
	}
	return rs.Extract(pred)
}

// ExtractByBackends classify queries by its main backend.
func (rs *Results) ExtractByBackends() []Results {
	res := make([]Results, backends.Num)
	for _, x := range *rs {
		id := x.Backend()
		res[id] = append(res[id], x)
	}
	return res
}

// Insert new to rs without number limit.
func (rs *Results) TryInsert(info Info, rater rating.Rater, f filter.Filter) bool {
	if rater.Score() < f.MinRating {
		return false
	}
	r := New(rater, info)
	rs.InsertWithLimit(f.MaxNumber, &r)
	return true
}

// Merge new to rs without number limit.
func (rs *Results) Merge(new Results) {
	rs.MergeWithLimit(0, new)
}

// Insert new to rs without number limit.
func (rs *Results) Insert(new *Result) {
	rs.InsertWithLimit(0, new)
}

// InsertWithLimit insert new to rs and select top maxN of it.
func (rs *Results) InsertWithLimit(maxN int, new *Result) {
	rs.MergeWithLimit(maxN, Results{new})
}

func (rs *Results) MergeAllWithLimit(maxN int, s ...Results) {
	for _, x := range s {
		rs.MergeWithLimit(maxN, x)
	}
}

// MergeWithLimit merge new to rs and select top maxN of it
func (rs *Results) MergeWithLimit(maxN int, new Results) {
	old := *rs
	ln := len(new)
	lo := len(old)
	if maxN == 0 {
		maxN = ln + lo
	}
	if ln == 0 {
		*rs = old[:num.Min(lo, maxN)]
		return
	}
	if lo == 0 {
		*rs = new[:num.Min(ln, maxN)]
		return
	}
	lr := num.Min(ln+lo, maxN)
	res := make(Results, lr)
	j, k := 0, 0
	for i := 0; i < lr; i++ {
		if j == len(old) {
			copy(res[i:], new[k:])
			break
		}
		if k == len(new) {
			copy(res[i:], old[j:])
			break
		}
		oldV, newV := old[j], new[k]
		if oldV.Rater.Score() < newV.Rater.Score() {
			res[i] = newV
			k++
		} else {
			res[i] = oldV
			j++
		}
	}
	*rs = res
}

func (r *Result) ShowRating() template.HTML {
	return desc.Rating(r.Rater.Score())
}
