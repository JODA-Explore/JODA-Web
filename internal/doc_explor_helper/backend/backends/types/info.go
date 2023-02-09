package types

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/cond"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/rating"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

type Info struct {
	result.BasicInfo
	percent num.Percent
	conds   cond.Conditions
}

func (i *Info) updateRating() {
	rat := rating.NewCoverage(
		i.percent.Ratio(),
		rating.PathFactor(0, false, i.JsonPoint()),
	)
	i.SetRating(&rat)
}
