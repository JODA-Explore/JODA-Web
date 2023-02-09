package structure

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/cond"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/rating"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

type Info struct {
	result.BasicInfo
	conds   cond.Conditions
	percent num.Percent
}

func (i *Info) updateRating() {
	rat := rating.NewCoverage(i.percent.Ratio(), rating.PathFactor(3, false, i.JsonPoint()))
	i.SetRating(&rat)
}

func (i Info) Info() Info {
	return i
}

func (i Info) Factors() rating.Factors {
	return rating.Factors{newHybridCondsBonus(i.conds)}
}
