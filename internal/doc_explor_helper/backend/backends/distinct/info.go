package distinct

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/freq"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/rating"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/jsontype"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

type Info struct {
	result.BasicInfo
	valueType            jsontype.Type
	min, max             uint64
	topPercent, coverage num.Percent
	sampleTop            freq.Freq
	estimatedRating      rating.Rater
}

func (i Info) EstimatedRating() rating.Rater {
	return i.estimatedRating
}

func (i *Info) updateRating() {
	var rat rating.Rater
	if i.topPercent.Empty() {
		// something wrong happened, skip this info
		rat = rating.ConstRating(0)
	} else {
		r := rating.NewReversedCoverage(
			i.topPercent.Ratio(),
			rating.PathFactor(0, false, i.JsonPoint()),
			rating.CoverageFactor{Val: i.coverage.Ratio()},
			rating.GlobalScale(-5),
		)
		rat = &r
	}
	i.SetRating(rat)
}
