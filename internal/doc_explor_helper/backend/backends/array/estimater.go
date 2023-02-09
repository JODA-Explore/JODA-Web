package array

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/feature"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/rating"
)

var _ feature.Estimater[Info] = Backend{} // Verify Merger is implemented.

func (b Backend) Estimate(info *Info) rating.Rater {
	r := rating.NewReversedCoverage(
		info.coverage.Ratio(),
		rating.PathFactor(0, false, info.JsonPoint()),
	)
	info.estimatedRating = &r
	return &r
}

func (b Backend) EstNum(f filter.Filter) int {
	return b.CandidatesNumber(f)
}

func (b Backend) Complete(info *Info) {
	b.complete(info, false)
}
