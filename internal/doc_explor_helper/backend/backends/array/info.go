package array

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/freq"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/rating"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/jsontype"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

type Info struct {
	result.BasicInfo
	arrayType            jsontype.Type
	coverage, topPercent num.Percent
	sampleTop            freq.Freq
	estimatedRating      rating.Rater
}

func (i *Info) updateRating() {
	rat := rating.NewReversedCoverage(
		i.topPercent.Ratio(),
		rating.PathFactor(0, false, i.JsonPoint()),
		rating.CoverageFactor{Val: i.coverage.Ratio()},
		rating.GlobalScale(50),
	)
	i.SetRating(&rat)
}

func (i Info) EstimatedRating() rating.Rater {
	return i.estimatedRating
}
