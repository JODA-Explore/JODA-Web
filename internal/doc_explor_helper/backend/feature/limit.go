package feature

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

type Limit struct {
	ratio float64
	min   int
}

func NewLimit(ratio float64, min int) Limit {
	return Limit{ratio, min}
}

func (l Limit) CandidatesNumber(f filter.Filter) int {
	return num.Max(l.min, num.Mul(f.MaxNumber, l.ratio))
}
