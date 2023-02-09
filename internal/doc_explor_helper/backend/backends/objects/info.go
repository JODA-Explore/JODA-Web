package objects

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/rating"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

type Info struct {
	result.BasicInfo
	real, sample, coverage num.Percent
}

func (i *Info) updateRating() {
	rat := rating.NewCoverage(
		i.real.Ratio(),
		rating.PathFactor(400, false, i.JsonPoint()),
		rating.CoverageFactor{Val: i.coverage.Ratio()},
		rating.GlobalScale(40),
	)
	i.SetRating(&rat)
}
