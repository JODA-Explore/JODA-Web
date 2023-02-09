package distinct

import (
	"fmt"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/feature"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/rating"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/jsontype"
)

var _ feature.Estimater[Info] = Backend{} // Verify Estimater is implemented.

func (b Backend) Estimate(info *Info) rating.Rater {
	var r rating.Rating
	switch info.valueType {
	case jsontype.String:
		r = rating.NewReversedCoverage(
			info.coverage.Ratio(),
			sizeFactor(info.max),
			newMinMaxFactor(info.min, info.max),
			rating.PathFactor(0, false, info.JsonPoint()),
		)
	case jsontype.Int:
		r = rating.NewReversedCoverage(
			info.coverage.Ratio(),
			newMinMaxFactor(info.min, info.max),
			rating.PathFactor(0, false, info.JsonPoint()),
		)
	case jsontype.Float:
		r = rating.NewReversedCoverage(
			info.coverage.Ratio(),
			rating.PathFactor(0, false, info.JsonPoint()),
		)
	default:
		panic(fmt.Errorf("the value type of %v is none of the following types:string,int,float", info))
	}
	info.estimatedRating = &r
	return &r
}

func (b Backend) EstNum(f filter.Filter) int {
	return b.CandidatesNumber(f)
}

func (b Backend) Complete(info *Info) {
	b.complete(info, false)
}
